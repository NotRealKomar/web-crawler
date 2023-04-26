package elastic

import (
	"bytes"
	"io"
	"os"
	"web-crawler/modules/logger"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

type ElasticService struct {
	client *elasticsearch.Client
}

func NewService() *ElasticService {
	service := ElasticService{}

	cacertFile, readFileErr := os.ReadFile("CHANGE_ME")

	if readFileErr != nil {
		logger.Log(readFileErr)

		return nil
	}

	configuration := elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200",
		},
		Username: "elastic",
		Password: "CHANGE_ME",
		CACert:   cacertFile,
	}

	client, _ := elasticsearch.NewClient(configuration)
	service.client = client

	return &service
}

func (service *ElasticService) Status() {
	res, healthErr := service.client.Cluster.Health()

	if healthErr != nil {
		logger.Log(healthErr)
		return
	}
	defer res.Body.Close()

	buffer := new(bytes.Buffer)
	_, copyErr := io.Copy(buffer, res.Body)

	if copyErr != nil {
		logger.Log(copyErr)
		return
	}

	logger.Log(buffer.String())
}
