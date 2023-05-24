package parser_test

import (
	"bytes"
	"io"
	"testing"
	"web-crawler/modules/logger"
	"web-crawler/modules/parser"
)

const mockHtmlData = `
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="utf-8">
		</head>
		<body>
			<h1>Hello</h1>
			<a href="https://www.baconmockup.com">test</a>
		</body>
	</html>
`

var parserService *parser.ParserService

func beforeTest() {
	mockLogger := logger.NewMockLoggerService()
	logger := logger.LoggerService(*mockLogger)

	parserService = parser.NewParserService(&logger)
}

func TestParse(t *testing.T) {
	beforeTest()

	reader := io.NopCloser(bytes.NewReader([]byte(mockHtmlData)))

	response, err := parserService.Parse(reader)
	if err != nil {
		t.Error(err)
	}

	if len(response.Data) == 0 {
		t.Error("Parse failed: response data is empty")
	}

	if response.Data[0] != "Hello" {
		t.Error("Parse failed: response data is not 'Hello'")
	}

	if response.Data[1] != "test" {
		t.Error("Parse failed: response data is not 'test'")
	}

	if len(response.Links) == 0 {
		t.Error("Parse failed: response links is empty")
	}

	if response.Links[0].String() != "https://www.baconmockup.com" {
		t.Error("Parse failed: response links is not 'https://www.baconmockup.com'")
	}
}
