package DI

import (
	"reflect"
	"web-crawler/modules/logger"
)

var dependencies map[reflect.Type]any

func Register[T any](value T) {
	if dependencies == nil {
		dependencies = make(map[reflect.Type]any)
	}

	key := reflect.TypeOf(value)

	if key.Kind() != reflect.Pointer {
		logger.Fatal("Cannot register dependency \"", key, "\": must be a pointer type")
	}

	dependencies[key] = value
}

func Inject[T any](output *T) {
	key := reflect.TypeOf(output)
	dependency := dependencies[key]

	if dependency != nil && reflect.ValueOf(dependency).Type().AssignableTo(key) {
		*output = *dependency.(*T)
		return
	}

	logger.Fatal("Cannot get dependency with key: ", key)
}
