package service

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

type Elastic struct {
	es *elasticsearch.Client
}

func New(es *elasticsearch.Client) *Elastic  {
	return &Elastic{
		es: es,
	}
}

func InitEs()  *elasticsearch.Client {
	// Initialize a client with the default settings.
	//
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	//
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	return es
}