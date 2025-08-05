package main

import (
	"log"
	"net/http"

	"github.com/yashzod/splitlinks/internal/app"
	"github.com/yashzod/splitlinks/internal/handler"
	"github.com/yashzod/splitlinks/internal/middleware"
	"github.com/yashzod/splitlinks/internal/repository"
)

func main() {
	ctx := &app.App{}
	db := repository.InitDB()
	ctx.DB = db
	mux := http.NewServeMux()
	mux.HandleFunc("/r/", handler.Route(ctx))
	mux.HandleFunc("/create_experiment/", handler.CreateExperiment(ctx))
	mux.HandleFunc("/get_experiment", handler.GetExperiment(ctx))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", middleware.CorsMiddleware(mux)))
}
