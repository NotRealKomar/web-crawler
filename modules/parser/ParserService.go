package parser

import (
	"errors"
	"io"
	"net/url"
	"regexp"
	"strings"

	"web-crawler/modules/logger"
	parserErrors "web-crawler/modules/parser/errors"

	"golang.org/x/exp/slices"
	"golang.org/x/net/html"
)

const (
	HTML_ANCHOR_TAG            = "a"
	HTML_ANCHOR_HREF_ATTRIBUTE = "href"
	HTML_PARAGRAPH_TAG         = "p"
	HTML_SPAN_TAG              = "span"
	HTML_HEADER_TAG            = "h1"

	REGEX_NEWLINE_EXPRESSION = "(?m)\n {1,}"
)

type ParseData struct {
	Data  []string
	Links []url.URL
}

type ParseInputReader io.ReadCloser
type ParserService struct{}

func (*ParserService) Parse(reader ParseInputReader) (*ParseData, error) {
	output := &ParseData{}

	tokenizer := html.NewTokenizer(reader)
	defer reader.Close()

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				break
			}

			logger.Log(parserErrors.ErrorTokenException, tokenType)

			return nil, errors.New(parserErrors.ErrorTokenException)
		}

		if tokenType == html.StartTagToken {
			token := tokenizer.Token()

			if token.Data == HTML_ANCHOR_TAG {
				idx := slices.IndexFunc(token.Attr, func(attr html.Attribute) bool {
					return attr.Key == HTML_ANCHOR_HREF_ATTRIBUTE
				})

				if idx != -1 {
					link := token.Attr[idx].Val

					url, err := url.ParseRequestURI(link)

					if err == nil && url.IsAbs() {
						output.Links = append(output.Links, *url)
					}
				}
			}

			if token.Data == HTML_PARAGRAPH_TAG ||
				token.Data == HTML_SPAN_TAG ||
				token.Data == HTML_ANCHOR_TAG ||
				token.Data == HTML_HEADER_TAG {
				innerTokenType := tokenizer.Next()

				if innerTokenType == html.TextToken {
					innerToken := tokenizer.Token()
					newLineRegexp := regexp.MustCompile(REGEX_NEWLINE_EXPRESSION)

					data := newLineRegexp.ReplaceAllString(strings.TrimSpace(innerToken.Data), " ")

					if len(data) > 0 {
						output.Data = append(output.Data, data)
					}
				}
			}
		}
	}

	return output, nil
}
