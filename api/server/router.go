package server

import "net/http"

func Test(w http.ResponseWriter, req *http.Request)  {
	w.Write([]byte("Hello World"))
}

func (s *Server) initializeRoute()  {
	s.Router.HandleFunc("/", Test).Methods("GET")
}
