package controller

import (
	"encoding/json"
	"fmt"
	"github.com/CorentinCrz/abstracts/model"
	"github.com/CorentinCrz/abstracts/service"
	"github.com/elastic/go-elasticsearch/v8"
	"io/ioutil"
	"log"
	"net/http"
)

type Controller struct {
	Db *service.Elastic
}

func (c *Controller) respond(w http.ResponseWriter, _ *http.Request, data interface{}, status int) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Cannot format json. err=%v\n", err)
	}
}

func (c *Controller) ErrorHandler(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	http.Error(w, http.StatusText(500), 500)
}


func New(es *elasticsearch.Client) *Controller {
	return &Controller{
		Db: service.New(es),
	}
}

func (c *Controller) PostBook(w http.ResponseWriter, r *http.Request)  {
	jsonBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		c.ErrorHandler(w, err)
		return
	}

	var book model.CreateBook
	err = json.Unmarshal(jsonBody, &book)
	if err != nil {
		c.ErrorHandler(w, err)
		return
	}

	err = c.Db.CreateBook(book)
	if err != nil {
		c.ErrorHandler(w, err)
		return
	}
	c.respond(w, r, book, 200)
}

func (c *Controller) GetBook(w http.ResponseWriter, r *http.Request)  {
	//params := mux.Vars(r)
	title := r.FormValue("title")
	fmt.Println(title)
	w.Write([]byte(title))
	b, _ := c.Db.GetBook()
	c.respond(w, r, b, 200)
}