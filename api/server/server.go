package server

import (
	"fmt"
	"github.com/CorentinCrz/abstracts/service"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"os"
)

type Server struct {
	Router *mux.Router
	es *elasticsearch.Client
}


func New(router *mux.Router) *Server {
	return &Server{
		Router: router,
		es: service.InitEs(),
	}
}

func (s *Server) Run() *Server  {
	s.initializeRoute()
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{os.Getenv("FRONT_URL"), os.Getenv("DOCUMENTATION_URL")},
		AllowedHeaders: []string{"Authorization", "Content-Type", "accept"},
		AllowedMethods: []string{"POST", "GET", "PUT", "DELETE", "PATCH"},
		AllowCredentials: true,
		Debug: false,
	})
	handler := crs.Handler(s.Router)

	err := http.ListenAndServe(":" + os.Getenv("PORT"), handler)
	if err != nil {
		fmt.Printf("%v", err)
		panic(err)
	}
	return s

}
