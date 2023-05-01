package indexers

import (
	"bytes"
	"encoding/json"
	"web-crawler/modules/documents"
	"web-crawler/modules/elastic/client"
	"web-crawler/modules/logger"
)

type ContentIndexer struct{}

func (*ContentIndexer) IndexFromDocument(document documents.ContentDocument) (*documents.ContentDocument, error) {
	client, getClientErr := client.GetClient()

	if getClientErr != nil {
		return nil, getClientErr
	}

	data, marshalErr := json.Marshal(document)

	if marshalErr != nil {
		return nil, marshalErr
	}

	response, indexErr := client.Index("content", bytes.NewReader(data))

	if indexErr != nil {
		return nil, indexErr
	}
	defer response.Body.Close()

	logger.Log(response)

	return &document, nil
}
