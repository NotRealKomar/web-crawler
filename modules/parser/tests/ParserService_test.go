package parser_test

import (
	"bytes"
	"io"
	"net/url"
	"reflect"
	"testing"
	"web-crawler/modules/logger"
	"web-crawler/modules/parser"
)

var parserService *parser.ParserService

func beforeTest(t *testing.T) {
	parserService = parser.NewParserService(
		&logger.NewMockLoggerService(t).LoggerService,
	)
}

func TestParse(t *testing.T) {
	beforeTest(t)

	testCases := []struct {
		message  string
		input    string
		expected parser.ParseData
	}{
		{
			message: "Should parse valid HTML data",
			input: `
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
			`,
			expected: parser.ParseData{
				Data: []string{"Hello", "test"},
				Links: []url.URL{
					{
						Scheme: "https",
						Host:   "www.baconmockup.com",
					},
				},
			},
		},
		{
			message: "Should parse HTML data with only links",
			input: `
				<!DOCTYPE html>
				<html>
					<head>
						<meta charset="utf-8">
					</head>
					<body>
						<a href="https://www.baconmockup.com">bacon</a>
						<a href="https://www.google.com">google</a>
						<a href="https://wordpress.org/plugins/any-ipsum/">ipsums</a>
					</body>
				</html>
			`,
			expected: parser.ParseData{
				Data: []string{"bacon", "google", "ipsums"},
				Links: []url.URL{
					{
						Scheme: "https",
						Host:   "www.baconmockup.com",
					},
					{
						Scheme: "https",
						Host:   "www.google.com",
					},
					{
						Scheme: "https",
						Host:   "wordpress.org",
						Path:   "/plugins/any-ipsum/",
					},
				},
			},
		},
		{
			message: "Should parse HTML data with only data",
			input: `
        <!DOCTYPE html>
        <html>
          <head>
            <meta charset="utf-8">
          </head>
          <body>
            <h1>Hello</h1>
						<p>World</p>
          </body>
        </html>
      `,
			expected: parser.ParseData{
				Data: []string{"Hello", "World"},
			},
		},
	}

	for _, testCase := range testCases {
		data, err := parserService.Parse(
			io.NopCloser(bytes.NewBufferString(testCase.input)),
		)
		if err != nil {
			t.Errorf("%v - test failed: %v", testCase.message, err)
		}

		if !reflect.DeepEqual(data.Data, testCase.expected.Data) {
			t.Errorf("%v - test failed: expected %#v, got %#v", testCase.message, testCase.expected.Data, data.Data)
		}

		if !reflect.DeepEqual(data.Links, testCase.expected.Links) {
			t.Errorf("%v - test failed: expected %#v, got %#v", testCase.message, testCase.expected.Links, data.Links)
		}
	}
}
