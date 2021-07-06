package main

import (
	"github.com/CorentinCrz/abstracts/api/server"
	"github.com/gorilla/mux"
)

func main()  {
	router := mux.NewRouter()
	s := server.New(router)
	s.Run()
}