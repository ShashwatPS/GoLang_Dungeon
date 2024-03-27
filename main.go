package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func BookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	page := vars["page"]

	responseData := map[string]string{
		"title": title,
		"page":  page,
	}

	jsonResponse, err := json.Marshal(responseData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/books/{title}/page/{page}", BookHandler).Methods("GET")

	http.ListenAndServe(":3000", r)
}
