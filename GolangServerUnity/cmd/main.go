package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zleodai/FSAESoftware/GolangServerUnity/handlers"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/test", handlers.Test).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
