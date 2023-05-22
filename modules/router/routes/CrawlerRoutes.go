package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"web-crawler/modules/DI"
	"web-crawler/modules/crawler"
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

		processIdChannel := make(chan string)
		go crawler.InitializeCrawl(crawlUrl, processIdChannel)

		processId := <-processIdChannel

		w.Write([]byte("Crawl process started.\nJob Id:\n" + processId + "\n\n"))
	}
}

func GetCheckCrawlRoute() types.RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		crawler := crawler.CrawlerService{}
		DI.Inject(&crawler)

		query := r.URL.Query()

		jobId := query.Get("id")
		if jobId == "" {
			w.Write([]byte("Id cannot be empty\n"))
			return
		}

		status := crawler.GetCrawlStatus(jobId)

		w.Write([]byte("Job #" + jobId + " - \"" + string(status) + "\"\n"))
	}
}
