package parser

import (
	"errors"
	"io"
	"strings"
	"web-crawler/modules/logger"

	"golang.org/x/net/html"
)

type ParseData struct {
	Data []string
	// links []string
}

type ParseInputReader io.ReadCloser
type ParserService struct{}

func (*ParserService) Parse(reader ParseInputReader) (output *ParseData, parseError error) {
	output = &ParseData{}

	tokenizer := html.NewTokenizer(reader)
	defer reader.Close()

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				break
			}

			logger.Log(ErrorTokenException, tokenType)

			return nil, errors.New(ErrorTokenException)
		}

		if tokenType == html.TextToken {
			data := strings.ReplaceAll(tokenizer.Token().Data, "\n", "")
			data = strings.ReplaceAll(data, "\t", "")

			if len(strings.ReplaceAll(data, " ", "")) > 0 {
				output.Data = append(output.Data, strings.Trim(data, " "))
			}
		}
	}

	return output, nil
}
