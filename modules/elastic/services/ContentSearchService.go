package services

import (
	"strings"
	"web-crawler/modules/elastic/models"
	"web-crawler/modules/elastic/repositories"
	"web-crawler/modules/types"
)

const RESPONSE_DATA_OFFSET = 10

type ContentSearchService struct {
	repository repositories.ContentRepositoryBase
}

func NewContentSearchService(
	repository repositories.ContentRepositoryBase,
) *ContentSearchService {
	return &ContentSearchService{
		repository,
	}
}

func (service *ContentSearchService) SearchByKeyword(search string, pagination *types.PaginationOptions) (
	*types.PaginatedResponse[models.SearchResponseModel],
	error,
) {
	documents, getManyByKeywordErr := service.repository.GetManyByKeyword(search, pagination)
	if getManyByKeywordErr != nil {
		return nil, getManyByKeywordErr
	}

	outputModels := []models.SearchResponseModel{}
	for _, document := range documents {
		idx := strings.Index(document.Data, search)

		startIdx, endIdx := idx-len(search)-RESPONSE_DATA_OFFSET, idx+len(search)*2+RESPONSE_DATA_OFFSET

		if startIdx < 0 {
			startIdx = 0
		}
		if endIdx > len(document.Data) {
			endIdx = len(document.Data)
		}

		data := document.Data[startIdx:endIdx]

		outputModels = append(outputModels, models.SearchResponseModel{
			Data:   "..." + data + "...",
			Source: document.Source,
		})
	}

	return &types.PaginatedResponse[models.SearchResponseModel]{
		Models: outputModels,
		Page:   pagination.Page + 1,
	}, nil
}
