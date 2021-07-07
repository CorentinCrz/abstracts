package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CorentinCrz/abstracts/model"
	"log"
)

func formatResearch(author *string, title *string, abstract *string) string  {
	str := ""
	if &author != nil && *author != "" {
		str += "author: " + *author + ", "
	}
	if &title != nil && *title != "" {
		str += "title: " + *title + ", "
	}
	if &abstract != nil && *abstract != "" {
		str += "abstract: " + *abstract + ", "
	}
	return str
}

func (e *Elastic) GetBook(author *string, title *string, abstract *string) ([]model.Book, error)  {
	var r  map[string]interface{}
	res, err := e.es.Search(
		e.es.Search.WithContext(context.Background()),
		e.es.Search.WithIndex("books"),
		// e.es.Search.WithBody(&buf),
		e.es.Search.WithQuery(formatResearch(author, title, abstract)),
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