package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error Marshalling payload: %v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	type Res struct {
		Error string `json:"error"`
	}

	if statusCode > 499 {
		log.Printf("Internal server error: %s", msg)
	}

	respondWithJson(w, statusCode, Res{Error: msg})
}
