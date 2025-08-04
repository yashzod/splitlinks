package main

import (
	"log"
	"net/http"

	"github.com/yashzod/splitlinks/internal/handler"
)

func main() {
	http.HandleFunc("/r/", handler.Route)
	http.HandleFunc("/create_experiment/", handler.CreateExperiment)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
