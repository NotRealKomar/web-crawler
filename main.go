package main

import (
	"log"
	"web-crawler/modules/DI"
	"web-crawler/modules/crawler"
	"web-crawler/modules/elastic/repositories"
	"web-crawler/modules/elastic/services"
	httpModule "web-crawler/modules/http"
	"web-crawler/modules/logger"
	"web-crawler/modules/parser"
	"web-crawler/modules/router"

	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	registerDependencies()

	router := router.GetRouter()

	loggerService := logger.LoggerService{}
	DI.Inject(&loggerService)

	go loggerService.EnableChannelLogging()

	loggerService.Log("Run the server on http://localhost:3000/")
	http.ListenAndServe(":3000", logger.NewLogMiddleware(router, &loggerService))
}

func registerDependencies() {
	logger := logger.NewLoggerService()
	repository := repositories.NewContentRepository(logger)

	DI.Register(
		logger,
		nil,
	)
	DI.Register(repository, nil)
	DI.Register(
		services.NewContentSearchService(
			*repository,
		),
		nil,
	)
	DI.Register(
		crawler.NewCrawlerService(
			repository,
			&httpModule.HttpClientService{},
			&parser.ParserService{},
			logger,
		),
		nil,
	)
}
