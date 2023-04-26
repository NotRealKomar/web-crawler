package main

import (
	"web-crawler/modules/elastic"
	"web-crawler/modules/logger"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal(err)
	}

	// parserService := parser.ParserService{}
	// httpClientService := http.HttpClientService{}
	elasticService := elastic.NewService()

	// url, _ := url.Parse("https://go.dev/doc/tutorial/handle-errors")

	// reader, _ := httpClientService.Get(url)
	// data, _ := parserService.Parse(*reader)

	// for _, value := range data.Data {
	// 	fmt.Println(value)
	// }

	elasticService.Status()
}
