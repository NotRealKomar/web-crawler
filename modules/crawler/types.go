package crawler

import (
	"net/url"
	"web-crawler/modules/elastic/repositories"
	"web-crawler/modules/http"
	"web-crawler/modules/logger"
	"web-crawler/modules/parser"
)

type CrawlLink struct {
	Url   *url.URL
	Depth int
}

type CrawlStatus string

const (
	IN_PROGRESS CrawlStatus = "in_progress"
	DONE        CrawlStatus = "done"
)

type CrawlJob map[string]CrawlStatus

type CrawlerService struct {
	contentRepository *repositories.ContentRepository
	httpClient        *http.HttpClientService
	parserService     *parser.ParserService
	loggerService     *logger.LoggerService
	jobList           CrawlJob
}
