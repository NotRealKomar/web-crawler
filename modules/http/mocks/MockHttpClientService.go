package http

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"
	httpModule "web-crawler/modules/http"
	"web-crawler/modules/logger"
)

type MockHttpClientService struct {
	httpModule.HttpClientService
	*testing.T
}

func NewMockHttpClientService(t *testing.T) *MockHttpClientService {
	var MockClient = &http.Client{}

	clientService := httpModule.NewHttpClientService(
		&logger.NewMockLoggerService(t).LoggerService,
		MockClient,
	)

	return &MockHttpClientService{
		*clientService,
		t,
	}
}

func (m *MockHttpClientService) Get(url *url.URL) (*io.ReadCloser, error) {
	output := io.NopCloser(bytes.NewReader([]byte("")))

	return &output, nil
}
