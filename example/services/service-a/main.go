package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func index(rw http.ResponseWriter, req *http.Request) {
	_, _ = rw.Write([]byte("OK"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods(http.MethodGet)
	log.Println("Starting service, on port: 8080")
	log.Panic(http.ListenAndServe(":8080", r))
}
