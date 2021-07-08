package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CorentinCrz/abstracts/model"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	"log"
	"strings"
)

func formatResearch(author *string, title *string, abstract *string) string {
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

func (e *Elastic) CreateBook(book model.CreateBook, id *string) error {
	bk := &model.Book{
		Author: book.Author,
		Abstract: book.Abstract,
		Title: book.Title,
	}
	if id != nil {
		bk.Id = *id
	} else {
		bk.Id = uuid.New().String()
	}

	// Set up the request object.
	bookJson, err := json.Marshal(bk)
	if err != nil {
		return err
	}
	req := esapi.IndexRequest{
		Index:   "books",
		Body:    strings.NewReader(string(bookJson)),
		Refresh: "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), e.es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Printf("[%s] Error indexing document", res)
		return err
	}
	return nil
}

func (e *Elastic) UpdateBook(id string, book model.CreateBook) error  {
	_, err := e.DeleteBook(id)
	if err != nil {
		return err
	}

	err = e.CreateBook(book, &id)
	if err != nil {
		return err
	}

	return nil
}

func (e *Elastic) GetBook(author *string, title *string, abstract *string) ([]model.Book, error) {
	var r map[string]interface{}
	res, err := e.es.Search(
		e.es.Search.WithContext(context.Background()),
		e.es.Search.WithIndex("books"),
		// e.esbooks?.Search.WithBody(&buf),
		e.es.Search.WithQuery(formatResearch(author, title, abstract)),
		e.es.Search.WithTrackTotalHits(true),
		e.es.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New("Error getting response")
	}
	
	
	

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing the response body: %s", err))
	}
	// Print the ID and document source for each hit.
	var b []model.Book
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"]
		b = append(b, model.Book{
			Id: fmt.Sprintf("%v", source.(map[string]interface{})["id"]),
			Title: fmt.Sprintf("%v", source.(map[string]interface{})["title"]),
			Author: fmt.Sprintf("%v", source.(map[string]interface{})["author"]),
			Abstract: fmt.Sprintf("%v", source.(map[string]interface{})["abstract"]),
		})
	}
	return b, nil
}

func (e *Elastic) DeleteBook(id string) (*esapi.Response, error) {

	var index = []string{"books"}
	req := esapi.DeleteByQueryRequest{
		Index: index,
		Query: "id: " + id,
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), e.es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	return res, err
}
