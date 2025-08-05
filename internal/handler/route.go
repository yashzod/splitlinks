package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/yashzod/splitlinks/internal/app"
	"github.com/yashzod/splitlinks/internal/repository"
	"github.com/yashzod/splitlinks/internal/service"
	"gorm.io/gorm"
)

var db *gorm.DB = repository.InitDB()

func Route(ctx *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/r/")
		if slug == "" {
			http.Error(w, "missing slug", http.StatusBadRequest)
			return
		}

		link, err := service.GetRedirectLink(ctx.DB, slug)
		if err != nil {
			http.Error(w, "url not found or expired", http.StatusNotFound)
			return
		}
		if link == "" {
			http.Error(w, "url not found", http.StatusNotFound)
			return
		}

		fmt.Println("Redirecting to:", link)
		http.Redirect(w, r, link, http.StatusFound)
	}
}

func CreateExperiment(ctx *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			fmt.Println("JSON decode error:", err)
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		slug, err := service.CreateExperiment(ctx.DB, data)
		if err != nil {
			fmt.Println("create experiment error:", err)
			http.Error(w, "failed to create experiment", http.StatusInternalServerError)
			return
		}

		response := map[string]string{"slug": slug}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			// Encoding failed after headers; log but client already got status
			fmt.Println("encode response error:", err)
		}
	}
}

func GetExperiment(ctx *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.Query().Get("slug")
		if slug == "" {
			http.Error(w, "missing slug", http.StatusBadRequest)
			return
		}

		query := map[string]string{"slug": slug}
		exp, err := service.GetExperiment(query, ctx.DB)
		if err != nil {
			fmt.Println("get experiment error:", err)
			http.Error(w, "experiment not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(exp); err != nil {
			fmt.Println("encode response error:", err)
		}
	}
}
