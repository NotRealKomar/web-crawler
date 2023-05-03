package client

import (
	"bytes"
	"io"
	"os"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

type elasticClientConfig struct {
	username   string
	password   string
	address    string
	cacertPath string
}

var client *elasticsearch.Client

func GetClient() (*elasticsearch.Client, error) {
	if client == nil {
		err := generateNewClient()

		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func Status() (string, error) {
	client, getClientErr := GetClient()

	if getClientErr != nil {
		return "", getClientErr
	}

	res, healthErr := client.Cluster.Health()

	if healthErr != nil {
		return "", healthErr
	}
	defer res.Body.Close()

	buffer := new(bytes.Buffer)
	_, copyErr := io.Copy(buffer, res.Body)

	if copyErr != nil {
		return "", copyErr
	}

	return buffer.String(), nil
}

func generateNewClient() error {
	config := loadConfig()

	cacertFile, readFileErr := os.ReadFile(config.cacertPath)

	if readFileErr != nil {
		return readFileErr
	}

	configuration := elasticsearch.Config{
		Addresses: []string{
			config.address,
		},
		Username: config.username,
		Password: config.password,
		CACert:   cacertFile,
	}

	newClient, _ := elasticsearch.NewClient(configuration)
	client = newClient

	return nil
}

func loadConfig() elasticClientConfig {
	return elasticClientConfig{
		username:   os.Getenv("ELASTICSEARCH_USERNAME"),
		password:   os.Getenv("ELASTICSEARCH_PASSWORD"),
		address:    os.Getenv("ELASTICSEARCH_ADDRESS"),
		cacertPath: os.Getenv("ELASTICSEARCH_CACERT_PATH"),
	}
}
