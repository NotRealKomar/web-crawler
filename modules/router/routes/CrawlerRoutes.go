package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"web-crawler/modules/DI"
	"web-crawler/modules/crawler"
)

func GetCrawlRoute() func(w http.ResponseWriter, r *http.Request) {
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

		crawlErr := crawler.InitializeCrawl(crawlUrl)
		if crawlErr != nil {
			w.Write([]byte(crawlErr.Error()))
			return
		}

		w.Write([]byte("Crawl process finished\n"))
	}
}
