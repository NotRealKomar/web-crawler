package main

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"

	"web-crawler/modules/crawler"
	"web-crawler/modules/documents"
	"web-crawler/modules/elastic/client"
	"web-crawler/modules/elastic/repositories"
	"web-crawler/modules/elastic/services"
	httpModule "web-crawler/modules/http"
	"web-crawler/modules/logger"
	"web-crawler/modules/parser"
	"web-crawler/modules/types"

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

	service := services.NewContentSearchService(
		*repository,
	)
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
	router.HandleFunc("/content/search", getSearchRoute(service)).Methods(http.MethodGet)

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
		return
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
			return
		}

		crawlUrl, parseRequestUriErr := url.ParseRequestURI(requestUrl)
		if parseRequestUriErr != nil {
			w.Write([]byte(parseRequestUriErr.Error()))
			return
		}

		crawlErr := service.InitializeCrawl(crawlUrl)
		if crawlErr != nil {
			w.Write([]byte(crawlErr.Error()))
			return
		}

		logger.Log("Crawl process finished")
		w.Write([]byte("Crawl process finished\n"))
	}
}

func getSearchRoute(service *services.ContentSearchService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		searchQuery := query.Get("search")
		if searchQuery == "" {
			w.Write([]byte("Search cannot be empty\n"))
			return
		}

		pageQuery := query.Get("page")
		if pageQuery == "" {
			pageQuery = "0"
		}

		page, atoiErr := strconv.Atoi(pageQuery)
		if atoiErr != nil {
			w.Write([]byte(atoiErr.Error()))
			return
		}

		pagination := types.NewPaginationOptions()
		pagination.Page = page

		searchResponse, searchByKeywordErr := service.SearchByKeyword(searchQuery, pagination)
		if searchByKeywordErr != nil {
			w.Write([]byte(searchByKeywordErr.Error()))
			return
		}

		response, marshalErr := json.Marshal(searchResponse)
		if marshalErr != nil {
			w.Write([]byte(marshalErr.Error()))
			return
		}

		w.Write(response)
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
