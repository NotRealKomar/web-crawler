package crawler

import (
	"net/url"
	"strings"
	"web-crawler/modules/documents"
	"web-crawler/modules/elastic/repositories"
	"web-crawler/modules/http"
	"web-crawler/modules/logger"
	"web-crawler/modules/parser"

	"golang.org/x/exp/slices"
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

func (service *CrawlerService) InitializeCrawl(link *url.URL) error {
	processedLinks := []url.URL{}

	crawlErr := service.crawl(link, &processedLinks, 0, 0)

	logger.Log("Crawl process finished, processed", len(processedLinks), "links")
	return crawlErr
}

func (service *CrawlerService) crawl(link *url.URL, processedLinks *[]url.URL, linkCount int, depth int) error {
	logger.Log("Start crawling on depth", depth, "in", link.String())

	if depth > MAX_LINK_COUNT {
		logger.Log("Link count", linkCount, "in depth", depth, "exceeded max link count in", link.String())
		return nil
	}

	if depth > MAX_DEPTH {
		logger.Log("Depth", depth, "exceeded max depth in", link.String())
		return nil
	}

	response, getErr := service.httpClient.Get(link)
	if getErr != nil {
		return getErr
	}

	parseData, parseErr := service.parserService.Parse(*response)
	if parseErr != nil {
		return parseErr
	}

	document := documents.NewContentDocument()
	document.Source = link.String()
	document.Data = strings.Join(parseData.Data, " ")

	service.contentRepository.Save(*document)

	*processedLinks = append(*processedLinks, *link)
	depth += 1

	innerLinkCount := 0

	for _, link := range parseData.Links {
		if innerLinkCount <= MAX_LINK_COUNT &&
			depth <= MAX_DEPTH &&
			isLinkUnique(link, processedLinks) {
			crawlErr := service.crawl(&link, processedLinks, innerLinkCount, depth)
			if crawlErr != nil {
				return crawlErr
			}

			innerLinkCount += 1
		}
	}

	return nil
}

func isLinkUnique(linkToCheck url.URL, links *[]url.URL) bool {
	return slices.IndexFunc(*links, func(link url.URL) bool {
		return linkToCheck.String() == link.String()
	}) == -1
}
