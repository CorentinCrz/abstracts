package service

import (
	"context"
	"encoding/json"
	"github.com/CorentinCrz/abstracts/model"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strings"
)

func (e *Elastic) CreateBook(book model.CreateBook) (error)  {
	var b strings.Builder
	b.WriteString(`{"title" : "`)
	b.WriteString(book.Title)
	b.WriteString(`","author" : "`)
	b.WriteString(book.Author)
	b.WriteString(`","abstract" : "`)
	b.WriteString(book.Abstract)
	b.WriteString(`"}`)

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      "books",
		Body:       strings.NewReader(b.String()),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), e.es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Printf("[%s] Error indexing document", res.Status())
		return err
	}
	return nil
}

func (e *Elastic) GetBook() ([]model.Book, error)  {
	var r  map[string]interface{}

	// Perform the search request.
	res, err := e.es.Search(
		e.es.Search.WithContext(context.Background()),
		e.es.Search.WithIndex("books"),
		// e.es.Search.WithBody(&buf),
		e.es.Search.WithTrackTotalHits(true),
		e.es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error getting response")
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the ID and document source for each hit.
	var b []model.Book
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"]
		b = append(b, model.Book{
			Title: source.(map[string]interface{})["title"],
			Author: source.(map[string]interface{})["author"],
			Abstract: source.(map[string]interface{})["abstract"],
		})
	}
	return b, nil
}