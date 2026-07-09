package main

import (
	"encoding/json"
	"fmt"
	"go-video-converter/internal/transcoder"
	"log"
	"net/http"
	"github.com/go-playground/validator/v10"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func health(w http.ResponseWriter, r *http.Request) {
	if err := writeJson(w, http.StatusOK, HealthResponse{Status: "ok"}); err != nil {
		log.Println("error while writing response", err)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var validate *validator.Validate = validator.New()

func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("error encoding the data, %w", err)
	}

	return nil
}

func post(w http.ResponseWriter, r *http.Request) {
	var response transcoder.ConvertParam

	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		if err := writeJson(w, http.StatusBadRequest, ErrorResponse{Error: "error decoding body"}); err != nil {
			log.Println("error while writing response", err)
		}
		return
	}

	err := validate.Struct(response)

	if err != nil {
		if err := writeJson(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"}); err != nil {
			log.Println("error while writing response", err)
		}
		return
	}

	if err := writeJson(w, http.StatusOK, response); err != nil {
		log.Println("error while writing response", err)
	}

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
