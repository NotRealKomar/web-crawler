package http

import (
	"io"
	"net/http"
	"net/url"
	"web-crawler/modules/logger"
)

type HttpClientServiceBase interface {
	Get(url *url.URL) (*io.ReadCloser, error)
}

type HttpClientService struct {
	HttpClientServiceBase
	logger logger.LoggerServiceBase
	client *http.Client
}

func NewHttpClientService(
	logger logger.LoggerServiceBase,
	client *http.Client,
) *HttpClientService {
	return &HttpClientService{
		logger: logger,
		client: client,
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
