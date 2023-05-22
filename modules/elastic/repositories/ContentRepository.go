package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"web-crawler/modules/documents"
	"web-crawler/modules/elastic/client"
	repositoryErrors "web-crawler/modules/elastic/repositories/errors"
	"web-crawler/modules/elastic/repositories/types"
	"web-crawler/modules/helpers"
	"web-crawler/modules/logger"
	generalTypes "web-crawler/modules/types"
)

type ContentRepository struct {
	loggerService *logger.LoggerService
}

const INDEX_NAME = "content"

func NewContentRepository(loggerService *logger.LoggerService) *ContentRepository {
	return &ContentRepository{loggerService}
}

func (repository *ContentRepository) GetMany() ([]documents.ContentDocument, error) {
	body, marshalErr := json.Marshal(
		types.QueryObject{
			"query": types.QueryObject{
				"match_all": types.QueryObject{},
			},
		},
	)
	if marshalErr != nil {
		return nil, marshalErr
	}

	return repository.getManyByQuery(body)
}

func (repository *ContentRepository) GetManyByKeyword(
	search string,
	pagination *generalTypes.PaginationOptions,
) ([]documents.ContentDocument, error) {
	body, marshalErr := json.Marshal(
		types.QueryObject{
			"size": pagination.Items,
			"from": pagination.Items * pagination.Page,
			"query": types.QueryObject{
				"wildcard": types.QueryObject{
					"data": types.QueryObject{
						"value": "*" + search + "*",
					},
				},
			},
		},
	)
	if marshalErr != nil {
		return nil, marshalErr
	}

	return repository.getManyByQuery(body)
}

func (repository *ContentRepository) Save(document documents.ContentDocument) {
	client, _ := client.GetClient()
	create := client.Create
	responseData := &types.CreateResponse{}

	body, marshalErr := json.Marshal(document)
	if marshalErr != nil {
		panic(marshalErr)
	}

	response, clientErr := create(
		INDEX_NAME,
		document.Id,
		bytes.NewReader(body),
		create.WithContext(context.Background()),
		create.WithPretty(),
	)
	if clientErr != nil {
		panic(clientErr)
	}

	if response.IsError() {
		logger.Log(*response)

		panic(repositoryErrors.CreateFailedException)
	}

	decodeErr := helpers.DecodeResponseBody(responseData, response.Body)
	if decodeErr != nil {
		panic(decodeErr)
	}

	if responseData.Result == "updated" {
		repository.loggerService.GetChannel() <- "unexpected result for 'create' request with id" + responseData.Id
	}
}

func (*ContentRepository) getManyByQuery(query []byte) ([]documents.ContentDocument, error) {
	client, _ := client.GetClient()
	search := client.Search
	responseData := &types.SearchResponse[documents.ContentDocument]{}

	response, searchErr := search(
		search.WithContext(context.Background()),
		search.WithIndex(INDEX_NAME),
		search.WithBody(bytes.NewReader(query)),
		search.WithTrackTotalHits(true),
		search.WithPretty(),
	)

	if searchErr != nil {
		return nil, searchErr
	}

	if response.IsError() {
		logger.Log(response)

		return nil, errors.New(repositoryErrors.SearchFailedException)
	}

	decodeErr := helpers.DecodeResponseBody(responseData, response.Body)
	if decodeErr != nil {
		return nil, decodeErr
	}

	if responseData.Hits.Hits == nil {
		return nil, errors.New(repositoryErrors.NoDocumentsException)
	}

	output := []documents.ContentDocument{}

	for _, document := range responseData.Hits.Hits {
		output = append(output, document.Source)
	}

	return output, nil
}
