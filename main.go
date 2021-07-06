package main

import (
	"github.com/CorentinCrz/abstracts/api/server"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

func init() {
	err := gotenv.Load()
	if err != nil {
		panic("failed to load .env file")
	}
}


func main()  {
	router := mux.NewRouter()
	s := server.New(router)
	s.Run()
}