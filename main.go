package main

import (
	"bytes"
	"encoding/json"
	"io"

	"web-crawler/modules/crawler"
	"web-crawler/modules/documents"
	"web-crawler/modules/elastic/client"
	"web-crawler/modules/elastic/repositories"
	httpModule "web-crawler/modules/http"
	"web-crawler/modules/logger"
	"web-crawler/modules/parser"

	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal(err)
	}

	repository := &repositories.ContentRepository{}
	crawler := crawler.NewCrawlerService(
		*repository,
		httpModule.HttpClientService{},
		parser.ParserService{},
	)

	router := mux.NewRouter()

	router.HandleFunc("/", notFound).Methods(http.MethodGet)
	router.HandleFunc("/healthcheck", healthcheck).Methods(http.MethodGet)
	router.HandleFunc("/status", status).Methods(http.MethodGet)

	router.HandleFunc("/content", getItemsRoute(repository)).Methods(http.MethodGet)
	router.HandleFunc("/content/new", getSaveRoute(repository)).Methods(http.MethodPost)

	router.HandleFunc("/crawler/test", getCrawlRoute(crawler)).Methods(http.MethodPost)

	logger.Log("Run the server on http://localhost:3000/")

	http.ListenAndServe(":3000", logger.NewLogMiddleware(router))
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("404 NOT FOUND"))
}

func status(w http.ResponseWriter, r *http.Request) {
	response, statusErr := client.Status()
	if statusErr != nil {
		w.Write([]byte(statusErr.Error()))
	}

	w.Write([]byte(response))
}

func getCrawlRoute(service *crawler.CrawlerService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
		}

		crawlUrl, parseRequestUriErr := url.ParseRequestURI(requestUrl)
		if parseRequestUriErr != nil {
			w.Write([]byte(parseRequestUriErr.Error()))
		}

		crawlErr := service.InitializeCrawl(crawlUrl)
		if crawlErr != nil {
			w.Write([]byte(crawlErr.Error()))
		}

		logger.Log("Crawl process finished")
		w.Write([]byte("Crawl process finished\n"))
	}
}

func getItemsRoute(repository *repositories.ContentRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		documents, findManyErr := repository.GetMany()
		if findManyErr != nil {
			w.Write([]byte(findManyErr.Error()))
			return
		}

		output, marshalErr := json.Marshal(documents)
		if marshalErr != nil {
			w.Write([]byte(marshalErr.Error()))
			return
		}

		w.Write(output)
	}
}

func getSaveRoute(repository *repositories.ContentRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		buffer := new(bytes.Buffer)
		io.Copy(buffer, r.Body)

		document := documents.NewContentDocument()

		unmarshalErr := json.Unmarshal(buffer.Bytes(), document)
		if unmarshalErr != nil {
			w.Write([]byte(unmarshalErr.Error()))
			return
		}

		newDocument, saveErr := repository.Save(*document)
		if saveErr != nil {
			w.Write([]byte(saveErr.Error()))
			return
		}

		response, marshalErr := json.Marshal(newDocument)
		if marshalErr != nil {
			w.Write([]byte(marshalErr.Error()))
			return
		}

		w.Write(response)
	}
}
