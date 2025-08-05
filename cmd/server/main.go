package main

import (
	"log"
	"net/http"

	"github.com/yashzod/splitlinks/internal/handler"
	"github.com/yashzod/splitlinks/internal/middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/r/", handler.Route)
	mux.HandleFunc("/create_experiment/", handler.CreateExperiment)
	mux.HandleFunc("/get_experiment", handler.GetExperiment)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", middleware.CorsMiddleware(mux)))
}
