package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"web-crawler/modules/documents"
	"web-crawler/modules/elastic/client"
	"web-crawler/modules/logger"
)

type ContentRepository struct{}
type queryObject map[string]any

func (*ContentRepository) FindMany() ([]documents.ContentDocument, error) {
	client, _ := client.GetClient()
	search := client.Search

	query := queryObject{
		"query": queryObject{
			"match_all": queryObject{},
		},
	}

	data, marshalErr := json.Marshal(query)

	if marshalErr != nil {
		return nil, marshalErr
	}

	response, searchErr := client.Search(
		search.WithContext(context.Background()),
		search.WithIndex("content"),
		search.WithBody(bytes.NewReader(data)),
		search.WithTrackTotalHits(true),
		search.WithPretty(),
	)

	if searchErr != nil {
		return nil, searchErr
	}
	defer response.Body.Close()

	if response.IsError() {
		logger.Log(response)

		return nil, errors.New("something went wrong with search request")
	}

	buffer := new(bytes.Buffer)
	_, copyErr := io.Copy(buffer, response.Body)

	if copyErr != nil {
		return nil, copyErr
	}

	logger.Log(buffer.String())

	return []documents.ContentDocument{}, nil
}
