package main

import (
	"encoding/json"
	"go-video-converter/internal/transcoder"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/Json")
	res := HealthResponse{Status: "ok"}
	json.NewEncoder(w).Encode(res)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var validate *validator.Validate = validator.New()

func post(w http.ResponseWriter, r *http.Request) {
	var response transcoder.ConvertParam

	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "error decoding body"})
		return
	}

	err := validate.Struct(response)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)

	mux.HandleFunc("POST /conversions", post)

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal("failed to start the server", err)
	}
}
