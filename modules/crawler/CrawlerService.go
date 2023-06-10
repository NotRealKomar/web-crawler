package crawler

import (
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"
	"web-crawler/modules/documents"
	"web-crawler/modules/elastic/repositories"
	"web-crawler/modules/http"
	"web-crawler/modules/logger"
	"web-crawler/modules/parser"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

func NewCrawlerService(
	contentRepository repositories.ContentRepositoryBase,
	httpClient http.HttpClientServiceBase,
	parserService parser.ParserServiceBase,
	loggerService logger.LoggerServiceBase,
) *CrawlerService {
	return &CrawlerService{
		contentRepository,
		httpClient,
		parserService,
		loggerService,
		make(CrawlJob),
	}
}

func (service *CrawlerService) GetJobStatus(jobId string) CrawlStatus {
	return service.jobList[jobId]
}

func (service *CrawlerService) InitializeCrawl(link *url.URL, processIdChannel chan string) {
	done := make(chan struct{})
	maxLinkCount := calculateMaxLinkCount(MAX_LINKS_PER_CRAWL, MAX_DEPTH_LEVEL)
	processedLinks := &[]url.URL{}

	processId := uuid.NewString()
	service.jobList[processId] = IN_PROGRESS

	processIdChannel <- processId

	linkCounter := 0

	go service.crawl(
		CrawlLink{
			Url:   link,
			Depth: 0,
		},
		&linkCounter,
		&maxLinkCount,
		processedLinks,
		done,
	)

	<-done

	time.Sleep(time.Second * DONE_SIGNAL_TIMEOUT)

	service.loggerService.GetChannel() <- "Crawl process is finished, processed " + strconv.Itoa(len(*processedLinks)) + " links out of " + strconv.Itoa(maxLinkCount)
	service.jobList[processId] = DONE
}

func (service *CrawlerService) crawl(
	link CrawlLink,
	linkCount *int,
	maxLinkCount *int,
	processedLinks *[]url.URL,
	statusChannel chan struct{},
) {
	if link.Depth > MAX_DEPTH_LEVEL {
		return
	}

	service.loggerService.GetChannel() <- "Start crawling in " + link.Url.String()

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

	go service.contentRepository.Save(*document)

	linksToParse := []url.URL{}

	*processedLinks = append(*processedLinks, *link.Url)

	for _, link := range parseData.Links {
		if getIsLinkIsUnique(processedLinks, &link) {
			linksToParse = append(linksToParse, link)
		}
	}

	service.loggerService.GetChannel() <- "Found " + strconv.Itoa(len(linksToParse)) + " unique links in " + link.Url.String()

	if len(linksToParse) == 0 {
		service.loggerService.GetChannel() <- "Terminating crawl process in " + link.Url.String()
		statusChannel <- struct{}{} // stop the crawl process entirely if didn't get any links

		return
	}

	if len(linksToParse) > MAX_LINKS_PER_CRAWL {
		linksToParse = linksToParse[0:MAX_LINKS_PER_CRAWL]
	}

	*linkCount += 1
	if *linkCount >= *maxLinkCount {
		statusChannel <- struct{}{}

		return
	}

	for _, crawlLink := range linksToParse {
		go func(crawlLink url.URL) {
			service.crawl(
				CrawlLink{
					Url:   &crawlLink,
					Depth: link.Depth + 1,
				},
				linkCount,
				maxLinkCount,
				processedLinks,
				statusChannel,
			)
		}(crawlLink)
	}

	service.loggerService.GetChannel() <- "Finish crawling in " + link.Url.String()
}

func getIsLinkIsUnique(linkSlice *[]url.URL, link *url.URL) bool {
	return slices.IndexFunc(*linkSlice, func(processedLink url.URL) bool {
		return processedLink.String() == link.String()
	}) == -1
}

func calculateMaxLinkCount(linksPerLevel int, maxDepth int) int {
	maxCount := 0.0

	for i := 0; i <= maxDepth; i++ {
		maxCount += math.Pow(float64(linksPerLevel), float64(i))
	}

	return int(maxCount)
}
