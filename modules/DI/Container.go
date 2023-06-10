package DI

import (
	"reflect"
	"web-crawler/modules/logger"
)

var dependencies map[reflect.Type]any
var loggerService *logger.LoggerService = logger.NewLoggerService()

func Register[T any](value T, key reflect.Type) {
	if dependencies == nil {
		dependencies = make(map[reflect.Type]any)
	}

	if key == nil {
		key = reflect.TypeOf(value)
	}

	if key.Kind() != reflect.Pointer {
		loggerService.Fatal("Cannot register dependency \"", key, "\": must be a pointer type")
	}

	dependencies[key] = value
}

func Inject[T any](output *T) {
	if dependencies == nil {
		loggerService.Fatal("Container has no dependencies")
	}

	key := reflect.TypeOf(output)
	dependency := dependencies[key]

	if dependency != nil && reflect.ValueOf(dependency).Type().AssignableTo(key) {
		*output = *dependency.(*T)
		return
	}

	loggerService.Fatal("Cannot get dependency with key: ", key)
}

func ClearDependencies() {
	dependencies = nil
}
