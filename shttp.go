package main

import (
	"log"
	"net/http"
)

type sHTTPHandler struct {
	handler http.Handler
}

func sHTTPServer(handler http.Handler) http.Handler {
	return &sHTTPHandler{handler}
}

func (s *sHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("oi")
	s.handler.ServeHTTP(w, r)
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", sHTTPServer(http.FileServer(http.Dir(".")))))
}
