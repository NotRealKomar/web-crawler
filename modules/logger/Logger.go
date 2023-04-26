package logger

import (
	"log"
	"os"
)

const LOG_PREFIX = "[CRAWLER]:"

var logger *log.Logger

func getLogger() *log.Logger {
	if logger == nil {
		logger = log.New(os.Stdout, LOG_PREFIX, log.Flags())
	}

	return logger
}

func Log(values ...any) {
	logger := getLogger()

	logger.Println(values...)
}

func Fatal(values ...any) {
	logger := getLogger()

	logger.Fatal(values...)
}
