package main

import (
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
		logger.Fatal(err)
	}

	registerDependencies()

	router := router.GetRouter()

	loggerService := logger.LoggerService{}
	DI.Inject(&loggerService)

	go loggerService.EnableLogging()

	logger.Log("Run the server on http://localhost:3000/")
	http.ListenAndServe(":3000", logger.NewLogMiddleware(router))
}

func registerDependencies() {
	logger := logger.NewLoggerService()
	repository := repositories.NewContentRepository(logger)

	DI.Register(repository)
	DI.Register(
		services.NewContentSearchService(
			*repository,
		),
	)
	DI.Register(
		crawler.NewCrawlerService(
			repository,
			&httpModule.HttpClientService{},
			&parser.ParserService{},
			logger,
		),
	)
	DI.Register(
		logger,
	)
}
