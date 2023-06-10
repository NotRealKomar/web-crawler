package crawler_test

import (
	"net/url"
	"testing"
	"time"
	"web-crawler/modules/crawler"
	elasticMocks "web-crawler/modules/elastic/repositories/mocks"
	httpMocks "web-crawler/modules/http/mocks"
	"web-crawler/modules/logger"
	"web-crawler/modules/parser"
)

var crawlerService *crawler.CrawlerService

const TEST_CRAWL_TIMEOUT = 1

func beforeTest(t *testing.T) {
	t.Parallel()

	mockLogger := logger.NewMockLoggerService(t)

	go mockLogger.EnableChannelLogging()

	crawlerService = crawler.NewCrawlerService(
		elasticMocks.NewMockContentRepository(t),
		httpMocks.NewMockHttpClientService(t),
		parser.NewMockParserService(t),
		&mockLogger.LoggerService,
	)

	t.Logf("%#v", crawlerService)
}

func TestInitializeCrawl(t *testing.T) {
	beforeTest(t)

	mockURL := &url.URL{
		Scheme: "https",
		Host:   "localhost",
		Path:   "/",
	}
	processIdChannel := make(chan string)

	go crawlerService.InitializeCrawl(mockURL, processIdChannel)

	jobId := <-processIdChannel

	time.Sleep(time.Second * (TEST_CRAWL_TIMEOUT + crawler.DONE_SIGNAL_TIMEOUT))

	jobStatus := crawlerService.GetJobStatus(jobId)
	t.Logf("Job ID: %s in status %v", jobId, jobStatus)

	if jobId == "" {
		t.Error("InitializeCrawl failed: job Id is empty")
	}

	if jobStatus != crawler.DONE {
		t.Error("InitializeCrawl failed: job status is not DONE")
	}
}
