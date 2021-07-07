package controller

import (
	"encoding/json"
	"fmt"
	"github.com/CorentinCrz/abstracts/service"
	"github.com/elastic/go-elasticsearch/v8"
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


func New(es *elasticsearch.Client) *Controller {
	return &Controller{
		Db: service.New(es),
	}
}

func (c *Controller) GetBook(w http.ResponseWriter, r *http.Request)  {
	author := r.FormValue("author")
	title := r.FormValue("title")
	abstract := r.FormValue("abstract")
	fmt.Println(title)
	b, _ := c.Db.GetBook(&author, &title, &abstract)
	c.respond(w, r, b, 200)
}