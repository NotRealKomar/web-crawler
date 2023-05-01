package main

import (
	"encoding/json"
	"web-crawler/modules/elastic/client"
	"web-crawler/modules/elastic/repositories"
	"web-crawler/modules/logger"

	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal(err)
	}

	repository := &repositories.ContentRepository{}

	// parserService := parser.ParserService{}
	// httpClientService := http.HttpClientService{}

	// url, _ := url.Parse("https://go.dev/doc/tutorial/handle-errors")

	// reader, _ := httpClientService.Get(url)
	// data, _ := parserService.Parse(*reader)

	// for _, value := range data.Data {
	// 	fmt.Println(value)
	// }

	mux := http.NewServeMux()

	mux.HandleFunc("/", notFound)
	mux.HandleFunc("/healthcheck", healthcheck)
	mux.HandleFunc("/status", status)
	mux.HandleFunc("/content", getItemsRoute(repository))

	http.ListenAndServe(":3000", logger.NewLogMiddleware(mux))
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

func getItemsRoute(repository *repositories.ContentRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		documents, findManyErr := repository.FindMany()
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
