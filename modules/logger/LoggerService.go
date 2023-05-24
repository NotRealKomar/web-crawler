package logger

import (
	"log"
	"os"
)

const LOG_PREFIX = "[CRAWLER]: "

type LoggerService struct {
	logChannel chan string
	logger     *log.Logger
}

func NewLoggerService() *LoggerService {
	logChannel := make(chan string)
	logger := log.New(os.Stdout, LOG_PREFIX, log.Flags())

	return &LoggerService{
		logChannel,
		logger,
	}
}

func (service *LoggerService) EnableChannelLogging() {
	defer close(service.logChannel)

	for value := range service.logChannel {
		service.Log(value)
	}
}

func (service *LoggerService) GetChannel() chan string {
	return service.logChannel
}

func (service *LoggerService) Log(values ...any) {
	service.logger.Println(values...)
}

func (service *LoggerService) Fatal(values ...any) {
	service.logger.Fatal(values...)
}
