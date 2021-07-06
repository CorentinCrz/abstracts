package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/CorentinCrz/abstracts/api/server"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	"log"
	"strconv"
	"strings"
	"sync"
)

func init() {
	err := gotenv.Load()
	if err != nil {
		panic("failed to load .env file")
	}
}


func main()  {
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
	//
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
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

	// 2. Index documents concurrently
	type Book struct {
		title string
		author string
		abstract string
	}
	//
	for i, book := range []Book{
			Book{title: "Texting",author: "Bob Wiki", abstract: "Text messaging, or texting, is the act of composing and sending electronic messages, typically consisting of alphabetic and numeric characters, between two or more users of mobile devices, desktops/laptops, or other type of compatible computer. Text messages may be sent over a cellular network, or may also be sent via an Internet connection."},
			Book{title: "Philippines",author: "Boby Wiki", abstract: "SMS was introduced to selected markets in the Philippines in 1995. In 1998, Philippine mobile-service providers launched SMS more widely across the country, with initial television marketing campaigns targeting hearing-impaired users. The service was initially free with subscriptions, but Filipinos quickly exploited the feature to communicate for free instead of using voice calls, which they would be charged for. After telephone companies realized this trend, they began charging for SMS. The rate across networks is 1 peso per SMS (about US$0.023). Even after users were charged for SMS, it remained cheap, about one-tenth of the price of a voice call. This low price led to about five million Filipinos owning a cell phone by 2001.[58] Because of the highly social nature of Philippine culture and the affordability of SMS compared to voice calls, SMS usage shot up. Filipinos used texting not only for social messages but also for political purposes, as it allowed the Filipinos to express their opinions on current events and political issues.[59] It became a powerful tool for Filipinos in promoting or denouncing issues and was a key factor during the 2001 EDSA II revolution, which overthrew then-President Joseph Estrada, who was eventually found guilty of corruption. According to 2009 statistics, there are about 72 million mobile-service subscriptions (roughly 80% of the Filipino population), with around 1.39 billion SMS messages being sent daily.[60][61] Because of the large number of text messages being sent, the Philippines became known as the \"text capital of the world\" during the late 1990s until the early 2000s."},
			Book{title: "New Zealand",author: "Toto", abstract: "There are three mobile network companies in New Zealand."},
			Book{title: "Africa",author: "Bob Wiki", abstract: "Text messaging will become a key revenue driver for mobile network operators in Africa over the next couple of years."},
		} {
		wg.Add(1)

		go func(i int, title string, author string, abstract string) {
			defer wg.Done()

			// Build the request body.
			var b strings.Builder
			b.WriteString(`{"title" : "`)
			b.WriteString(title)
			b.WriteString(`","author" : "`)
			b.WriteString(author)
			b.WriteString(`","abstract" : "`)
			b.WriteString(abstract)
			b.WriteString(`"}`)

			// Set up the request object.
			req := esapi.IndexRequest{
				Index:      "books",
				DocumentID: strconv.Itoa(i + 1),
				Body:       strings.NewReader(b.String()),
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
		}(i, book.title, book.author, book.abstract)
	}
	wg.Wait()

	log.Println(strings.Repeat("-", 37))

	// 3. Search for the indexed documents
	//
	// Build the request body.
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "books",
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("books"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	log.Println(strings.Repeat("=", 37))

	router := mux.NewRouter()
	s := server.New(router)
	s.Run()
}