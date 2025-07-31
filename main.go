package main

import (
	"log"
	"net/http"
	"strings"
)

type Experiment struct {
	VariantA string
	VariantB string
}

var experiments = map[string]Experiment{
	"abc123": {"https://site.com/a", "https://site.com/b"},
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	linkID := strings.TrimPrefix(r.URL.Path, "/r/")
	exp, ok := experiments[linkID]
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Basic A/B logic: 50-50 split
	if r.URL.Query().Get("v") == "b" {
		http.Redirect(w, r, exp.VariantB, http.StatusFound)
	} else {
		http.Redirect(w, r, exp.VariantA, http.StatusFound)
	}
}

func main() {
	http.HandleFunc("/r/", handleRedirect)
	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
