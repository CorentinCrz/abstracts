package server

import (
	"github.com/CorentinCrz/abstracts/controller"
	"net/http"
)

func Test(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello World"))
}

func (s *Server) initializeRoute() {
	c := controller.New(s.es)
	s.Router.HandleFunc("/", Test).Methods("GET")
	s.Router.HandleFunc("/books", c.PostBook).Methods("POST")
	s.Router.HandleFunc("/books", c.GetBook).Methods("GET")
	s.Router.HandleFunc("/books/{id}", c.PutBook).Methods("PUT")
	s.Router.HandleFunc("/books/{id}", c.DeleteBook).Methods("DELETE")
}
