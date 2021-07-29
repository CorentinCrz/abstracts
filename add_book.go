package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/CorentinCrz/abstracts/model"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	"github.com/subosito/gotenv"
)

func init() {
	err := gotenv.Load()
	if err != nil {
		panic("failed to load .env file")
	}
}

func add() {
	log.SetFlags(0)

	var (
		r  map[string]interface{}
		wg sync.WaitGroup
	)

	// Initialize a client with the default settings.
	//
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	//
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 1. Get cluster info

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	books := []model.Book{
		{Title: "E88E", Author: "Bob Wiki", Abstract: "Text messaging, or texting, is the act of composing and sending electronic messages, typically consisting of alphabetic and numeric characters, between two or more users of mobile devices, desktops/laptops, or other type of compatible computer. Text messages may be sent over a cellular network, or may also be sent via an Internet connection."},
		{Title: "EE981", Author: "Boby Wiki", Abstract: "SMS was introduced to selected markets in the Philippines in 1995. In 1998, Philippine mobile-service providers launched SMS more widely across the country, with initial television marketing campaigns targeting hearing-impaired users. The service was initially free with subscriptions, but Filipinos quickly exploited the feature to communicate for free instead of using voice calls, which they would be charged for. After telephone companies realized this trend, they began charging for SMS. The rate across networks is 1 peso per SMS (about US$0.023). Even after users were charged for SMS, it remained cheap, about one-tenth of the price of a voice call. This low price led to about five million Filipinos owning a cell phone by 2001.[58] Because of the highly social nature of Philippine culture and the affordability of SMS compared to voice calls, SMS usage shot up. Filipinos used texting not only for social messages but also for political purposes, as it allowed the Filipinos to express their opinions on current events and political issues.[59] It became a powerful tool for Filipinos in promoting or denouncing issues and was a key factor during the 2001 EDSA II revolution, which overthrew then-President Joseph Estrada, who was eventually found guilty of corruption. According to 2009 statistics, there are about 72 million mobile-service subscriptions (roughly 80% of the Filipino population), with around 1.39 billion SMS messages being sent daily.[60][61] Because of the large number of text messages being sent, the Philippines became known as the \"text capital of the world\" during the late 1990s until the early 2000s."},
		{Title: "EE772", Author: "Toto", Abstract: "There are three mobile network companies in New Zealand."},
		{Title: "E77E3", Author: "Bob Wiki", Abstract: "Text messaging will become a key revenue driver for mobile network operators in Africa over the next couple of years."},
	}

	for i, book := range books {
		wg.Add(1)

		go func(i int, title string, author string, abstract string) {
			defer wg.Done()
			book.Title = title
			book.Author = author
			book.Abstract = abstract
			book.Id = uuid.New().String()

			// Set up the request object.
			bookJson, err := json.Marshal(book)

			// Set up the request object.
			req := esapi.IndexRequest{
				Index:      "books",
				DocumentID: strconv.Itoa(i + 1),
				Body:       strings.NewReader(string(bookJson)),
				Refresh:    "true",
			}

			// Perform the request with the client.
			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and indexed document version.
					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}
		}(i, book.Title, book.Author, book.Abstract)
	}
	wg.Wait()
}
