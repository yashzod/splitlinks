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
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusNotFound)
	}
	service.CreateExperiment(data, db)
	fmt.Println(data)
}
