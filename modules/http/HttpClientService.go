package http

import (
	"io"
	"net/http"
	"net/url"
	"web-crawler/modules/logger"
)

type HttpClientService struct{}

func (*HttpClientService) Get(url *url.URL) (*io.ReadCloser, error) {
	response, getErr := http.Get(url.String())

	if getErr != nil {
		logger.Log(getErr)

		return nil, getErr
	}

	return &response.Body, nil
}
