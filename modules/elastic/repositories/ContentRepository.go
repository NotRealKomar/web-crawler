package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"

	"web-crawler/modules/documents"
	"web-crawler/modules/elastic/client"
	repositoryErrors "web-crawler/modules/elastic/repositories/errors"
	"web-crawler/modules/elastic/repositories/types"
	"web-crawler/modules/logger"
)

type ContentRepository struct{}

const INDEX_NAME = "content"

func (*ContentRepository) GetMany() ([]documents.ContentDocument, error) {
	client, _ := client.GetClient()
	search := client.Search

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

	response, searchErr := client.Search(
		search.WithContext(context.Background()),
		search.WithIndex(INDEX_NAME),
		search.WithBody(bytes.NewReader(body)),
		search.WithTrackTotalHits(true),
		search.WithPretty(),
	)

	if searchErr != nil {
		return nil, searchErr
	}
	defer response.Body.Close()

	if response.IsError() {
		logger.Log(response)

		return nil, errors.New(repositoryErrors.SearchFailedException)
	}

	buffer := new(bytes.Buffer)
	_, copyErr := io.Copy(buffer, response.Body)

	if copyErr != nil {
		return nil, copyErr
	}

	responseData := &types.SearchResponse[documents.ContentDocument]{}
	json.Unmarshal(buffer.Bytes(), responseData)

	if responseData.Hits.Hits == nil {
		return nil, errors.New(repositoryErrors.NoDocumentsException)
	}

	output := []documents.ContentDocument{}

	for _, document := range responseData.Hits.Hits {
		output = append(output, document.Source)
	}

	return output, nil
}
