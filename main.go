package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"network/Ume/src/middleware"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", middleware.AuthHandler).Methods("POST")

	log.Println("Server started at :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
