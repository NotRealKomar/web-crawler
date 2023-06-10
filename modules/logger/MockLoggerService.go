package logger

import (
	"log"
	"os"
	"testing"
)

type MockLoggerService struct {
	LoggerService
	*testing.T
}

const MOCK_LOG_PREFIX = "[MOCK]: "

func NewMockLoggerService(t *testing.T) *MockLoggerService {
	return &MockLoggerService{
		LoggerService{
			logChannel: make(chan string),
			logger:     log.New(os.Stdout, MOCK_LOG_PREFIX, log.Flags()),
		},
		t,
	}
}

func (service *MockLoggerService) Log(values ...any) {
	// Do nothing
}

func (service *MockLoggerService) Fatal(values ...any) {
	// Do nothing
}

func (service *MockLoggerService) EnableChannelLogging() {
	defer close(service.logChannel)

	for value := range service.logChannel {
		service.T.Log(value)
	}
}
