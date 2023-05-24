package logger

import (
	"log"
	"os"
)

type MockLoggerService LoggerService

const MOCK_LOG_PREFIX = "[MOCK]: "

func NewMockLoggerService() *MockLoggerService {
	return &MockLoggerService{
		logger: log.New(os.Stdout, LOG_PREFIX, log.Flags()),
	}
}

func (service *MockLoggerService) Log(values ...any) {
	// Do nothing
}

func (service *MockLoggerService) Fatal(values ...any) {
	// Do nothing
}

func (service *MockLoggerService) EnableChannelLogging() {
	// Do nothing
}
