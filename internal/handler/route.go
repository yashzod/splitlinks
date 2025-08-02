package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yashzod/splitlinks/internal/service"
)

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
