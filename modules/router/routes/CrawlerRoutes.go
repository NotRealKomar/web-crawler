package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"web-crawler/modules/DI"
	"web-crawler/modules/crawler"
	"web-crawler/modules/logger"
	"web-crawler/modules/types"
)

func GetCrawlRoute() types.RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		crawler := crawler.CrawlerService{}
		DI.Inject(&crawler)

		defer r.Body.Close()

		buffer := new(bytes.Buffer)
		io.Copy(buffer, r.Body)

		var body map[string]string
		unmarshalErr := json.Unmarshal(buffer.Bytes(), &body)
		if unmarshalErr != nil {
			w.Write([]byte(unmarshalErr.Error()))
			return
		}

		requestUrl := body["url"]
		if requestUrl == "" {
			w.Write([]byte("No URL provided\n"))
			return
		}

		crawlUrl, parseRequestUriErr := url.ParseRequestURI(requestUrl)
		if parseRequestUriErr != nil {
			w.Write([]byte(parseRequestUriErr.Error()))
			return
		}

		messageChannel := make(chan string)

		go crawler.InitializeCrawl(crawlUrl, messageChannel)
		go logger.LogChannel(messageChannel)

		w.Write([]byte("Crawl process started\n"))
	}
}

func GetCheckCrawlRoute() types.RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
