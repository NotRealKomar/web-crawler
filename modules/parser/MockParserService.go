package parser

import (
	"net/url"
	"testing"
	"web-crawler/modules/logger"
)

type MockParserService struct {
	ParserService
}

func NewMockParserService(t *testing.T) *MockParserService {
	return &MockParserService{
		ParserService: *NewParserService(
			&logger.NewMockLoggerService(t).LoggerService,
		),
	}
}

func (mock *MockParserService) Parse(reader ParseInputReader) (*ParseData, error) {
	return &ParseData{
		Data: []string{"Hello", "World"},
		Links: []url.URL{
			{
				Scheme: "https",
				Host:   "www.example.com",
				Path:   "/index.html",
			},
			{
				Scheme: "https",
				Host:   "www.google.com",
				Path:   "/index.html",
			},
			{
				Scheme: "https",
				Host:   "wordpress.org",
				Path:   "/plugins/any-ipsum/",
			},
		},
	}, nil
}
