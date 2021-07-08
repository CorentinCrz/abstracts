package controller

import (
	"encoding/json"
	"fmt"
	"github.com/CorentinCrz/abstracts/model"
	"github.com/CorentinCrz/abstracts/service"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gorilla/mux"
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

func (c *Controller) PostBook(w http.ResponseWriter, r *http.Request) {
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

	err = c.Db.CreateBook(book, nil)
	if err != nil {
		c.ErrorHandler(w, err)
		return
	}
	c.respond(w, r, book, 200)
}

func (c *Controller) PutBook(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)
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

	err = c.Db.UpdateBook(params["id"], book)
	if err != nil {
		c.ErrorHandler(w, err)
		return
	}
	c.respond(w, r, book, 200)
}

func (c *Controller) GetBook(w http.ResponseWriter, r *http.Request) {
	author := r.FormValue("author")
	title := r.FormValue("title")
	abstract := r.FormValue("abstract")
	fmt.Println(title)
	b, err := c.Db.GetBook(&author, &title, &abstract)
	if err != nil {
		c.respond(w, r, err, http.StatusInternalServerError)
	}
	if b == nil {
		c.respond(w, r, "Books not found", http.StatusNotFound)
	}
	c.respond(w, r, b, http.StatusOK)
}

func (c *Controller) DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reponse, err := c.Db.DeleteBook(params["id"])
	if err != nil {
		c.respond(w, r, err, 400)
		return
	}
	c.respond(w, r, reponse, 200)
}
