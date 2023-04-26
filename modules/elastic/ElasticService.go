package elastic

import (
	"bytes"
	"io"
	"os"
	"web-crawler/modules/logger"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

type elasticServiceConfig struct {
	username   string
	password   string
	address    string
	cacertPath string
}

type ElasticService struct {
	client *elasticsearch.Client
}

func NewService() *ElasticService {
	service := ElasticService{}
	config := loadConfig()

	cacertFile, readFileErr := os.ReadFile(config.cacertPath)

	if readFileErr != nil {
		logger.Log(readFileErr)

		return nil
	}

	configuration := elasticsearch.Config{
		Addresses: []string{
			config.address,
		},
		Username: config.username,
		Password: config.password,
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

func loadConfig() elasticServiceConfig {
	return elasticServiceConfig{
		username:   os.Getenv("ELASTICSEARCH_USERNAME"),
		password:   os.Getenv("ELASTICSEARCH_PASSWORD"),
		address:    os.Getenv("ELASTICSEARCH_ADDRESS"),
		cacertPath: os.Getenv("ELASTICSEARCH_CACERT_PATH"),
	}
}
