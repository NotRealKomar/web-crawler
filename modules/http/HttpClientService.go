package http

import (
	"io"
	"net/http"
	"net/url"
	"web-crawler/modules/logger"
)

type HttpClientService struct {
	logger *logger.LoggerService
}

func NewHttpClientService(logger *logger.LoggerService) *HttpClientService {
	return &HttpClientService{
		logger,
	}
}

func (client *HttpClientService) Get(url *url.URL) (*io.ReadCloser, error) {
	response, getErr := http.Get(url.String())

	if getErr != nil {
		client.logger.Log(getErr)

		return nil, getErr
	}

	return &response.Body, nil
}
