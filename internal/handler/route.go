package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/yashzod/splitlinks/internal/service"

	"github.com/yashzod/splitlinks/internal/repository"
	"gorm.io/gorm"
)

var db *gorm.DB = repository.InitDB()

func Route(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/r/")
	if slug == "" {
		http.Error(w, "Missing slug", http.StatusBadRequest)
		return
	}

	link, err := service.GetRedirectLink(slug)
	if err != nil {
		http.Error(w, "Url not found or expired", http.StatusNotFound)
	}
	fmt.Println("Redirecting slug:", link)

	http.Redirect(w, r, link, http.StatusFound)
}

func CreateExperiment(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	// Parse JSON body
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Println("JSON decode error:", err)
		http.Error(w, "decode error", http.StatusBadRequest)
		return
	}

	// Call service layer to create experiment
	err = service.CreateExperiment(data, db)
	if err != nil {
		http.Error(w, "Failed to create experiment", http.StatusInternalServerError)
		return
	}

	// Respond success
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Experiment created successfully")
}

func GetExperiment(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		http.Error(w, "Missing slug", http.StatusBadRequest)
		return
	}

	query := map[string]string{"slug": slug}

	exp, err := service.GetExperiment(query, db)
	if err != nil {
		http.Error(w, "Experiment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exp)
}