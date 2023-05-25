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
	contentRepository repositories.ContentRepositoryBase
	httpClient        http.HttpClientServiceBase
	parserService     parser.ParserServiceBase
	loggerService     logger.LoggerServiceBase
	jobList           CrawlJob
}
