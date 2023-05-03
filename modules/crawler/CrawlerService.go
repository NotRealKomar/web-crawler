package crawler

import (
	"net/url"
	"strings"
	"web-crawler/modules/documents"
	"web-crawler/modules/elastic/repositories"
	"web-crawler/modules/http"
	"web-crawler/modules/logger"
	"web-crawler/modules/parser"
)

const MAX_DEPTH = 3
const MAX_LINK_COUNT = 5

type CrawlerService struct {
	contentRepository repositories.ContentRepository
	httpClient        http.HttpClientService
	parserService     parser.ParserService
}

func NewCrawlerService(
	contentRepository repositories.ContentRepository,
	httpClient http.HttpClientService,
	parserService parser.ParserService,
) *CrawlerService {
	return &CrawlerService{
		contentRepository,
		httpClient,
		parserService,
	}
}

func (service *CrawlerService) InitializeCrawl(url *url.URL) error {
	return service.crawl(url, 0, 0)
}

func (service *CrawlerService) crawl(url *url.URL, linkCount int, depth int) error {
	logger.Log("Start crawling on depth", depth, "in", url.String())

	if depth >= MAX_LINK_COUNT {
		logger.Log("Link count", linkCount, "in depth", depth, "exceeded max link count in", url.String())
		return nil
	}

	if depth >= MAX_DEPTH {
		logger.Log("Depth", depth, "exceeded max depth in", url.String())
		return nil
	}

	response, getErr := service.httpClient.Get(url)
	if getErr != nil {
		return getErr
	}

	parseData, parseErr := service.parserService.Parse(*response)
	if parseErr != nil {
		return parseErr
	}

	document := documents.NewContentDocument()
	document.Source = url.String()
	document.Data = strings.Join(parseData.Data, " ")

	service.contentRepository.Save(*document)

	depth += 1

	for idx, link := range parseData.Links {
		if idx < MAX_LINK_COUNT && depth < MAX_DEPTH {
			crawlErr := service.crawl(&link, idx, depth)
			if crawlErr != nil {
				return crawlErr
			}
		}
	}

	return nil
}
