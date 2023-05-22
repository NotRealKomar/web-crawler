package routes

import (
	"net/http"
	"web-crawler/modules/elastic/client"
	"web-crawler/modules/types"
)

func GetHealthcheckRoute() types.RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}
}

func GetNotFoundRoute() types.RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("404 NOT FOUND"))
	}
}

func GetStatusRoute() types.RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		response, statusErr := client.Status()
		if statusErr != nil {
			w.Write([]byte(statusErr.Error()))
			return
		}

		w.Write([]byte(response))
	}
}
