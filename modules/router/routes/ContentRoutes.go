package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"web-crawler/modules/DI"
	"web-crawler/modules/elastic/repositories"
	"web-crawler/modules/elastic/services"
	"web-crawler/modules/types"
)

func GetSearchRoute() types.RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		service := services.ContentSearchService{}
		DI.Inject(&service)

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

func GetItemsRoute() types.RouteHandler {
	repository := repositories.ContentRepository{}
	DI.Inject(&repository)

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
