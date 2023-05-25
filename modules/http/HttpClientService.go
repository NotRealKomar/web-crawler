package http

import (
	"io"
	"net/http"
	"net/url"
	"web-crawler/modules/logger"
)

type HttpClientService struct {
	logger *logger.LoggerService
	client *http.Client
}

func NewHttpClientService(
	logger *logger.LoggerService,
	client *http.Client,
) *HttpClientService {
	return &HttpClientService{
		logger,
		client,
	}
}

func (client *HttpClientService) Get(url *url.URL) (*io.ReadCloser, error) {
	response, getErr := client.client.Get(url.String())

	if getErr != nil {
		client.logger.Log(getErr)

		return nil, getErr
	}

	return &response.Body, nil
}
