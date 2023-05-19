package crawler

import (
	"net/url"
	"strconv"
	"strings"
	"web-crawler/modules/documents"
	"web-crawler/modules/elastic/repositories"
	"web-crawler/modules/http"
	"web-crawler/modules/parser"
)

const MAX_LINK_COUNT = 5
const MAX_DEPTH = 4
const LINK_BUFFER_SIZE = 1

type CrawlLink struct {
	Url   *url.URL
	Depth int
}

type CrawlerService struct {
	contentRepository *repositories.ContentRepository
	httpClient        *http.HttpClientService
	parserService     *parser.ParserService
}

func NewCrawlerService(
	contentRepository *repositories.ContentRepository,
	httpClient *http.HttpClientService,
	parserService *parser.ParserService,
) *CrawlerService {
	return &CrawlerService{
		contentRepository,
		httpClient,
		parserService,
	}
}

func (service *CrawlerService) InitializeCrawl(link *url.URL, messageChannel chan string) {
	done := make(chan struct{})

	go service.crawl(
		CrawlLink{
			Url:   link,
			Depth: 0,
		},
		messageChannel,
		done,
	)

	<-done
}

func (service *CrawlerService) crawl(
	link CrawlLink,
	messageChannel chan string,
	statusChannel chan struct{},
) {
	if link.Depth > MAX_DEPTH {
		return
	}

	messageChannel <- "Start crawling in " + link.Url.String()

	response, getErr := service.httpClient.Get(link.Url)
	if getErr != nil {
		panic(getErr)
	}

	parseData, parseErr := service.parserService.Parse(*response)
	if parseErr != nil {
		panic(parseErr)
	}

	document := documents.NewContentDocument()
	document.Source = link.Url.String()
	document.Data = strings.Join(parseData.Data, " ")

	service.contentRepository.Save(*document)

	linksToParse := parseData.Links

	messageChannel <- "Found " + strconv.Itoa(len(linksToParse)) + " links in " + link.Url.String()

	if len(linksToParse) == 0 {
		messageChannel <- "Terminating crawl process in " + link.Url.String()

		return
	}

	if len(linksToParse) > MAX_LINK_COUNT {
		linksToParse = linksToParse[0:MAX_LINK_COUNT]
	}

	for _, crawlLink := range linksToParse {
		go service.crawl(
			CrawlLink{
				Url:   &crawlLink,
				Depth: link.Depth + 1,
			},
			messageChannel,
			statusChannel,
		)
	}

	statusChannel <- struct{}{}
}
