package startup

import (
	"crypto/tls"
	"fold/internal/config"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
)

var EsClient *elasticsearch.Client

// This func will be invoked by main, do all the necessary checks and startup ops here
func InitServer() {

	//connect elasticsearch
	cfg := elasticsearch.Config{

		Addresses: []string{config.GlobalConfig.ES_URL},
		Username:  config.GlobalConfig.ES_UserName,
		Password:  config.GlobalConfig.ES_Password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	var err error
	// Create Elasticsearch client
	EsClient, err = elasticsearch.NewClient(cfg)

	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %v", err)
	}

}
