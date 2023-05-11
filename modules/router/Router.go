package router

import (
	"net/http"
	"web-crawler/modules/router/routes"

	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", routes.GetNotFoundRoute()).Methods(http.MethodGet)
	router.HandleFunc("/healthcheck", routes.GetHealthcheckRoute()).Methods(http.MethodGet)
	router.HandleFunc("/status", routes.GetStatusRoute()).Methods(http.MethodGet)

	router.HandleFunc("/content", routes.GetItemsRoute()).Methods(http.MethodGet)
	router.HandleFunc("/content/search", routes.GetSearchRoute()).Methods(http.MethodGet)

	router.HandleFunc("/crawler/test", routes.GetCrawlRoute()).Methods(http.MethodPost)

	return router
}
