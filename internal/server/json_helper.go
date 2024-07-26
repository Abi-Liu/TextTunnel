package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJson(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error Marshalling payload: %v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(500)
	}
}

func RespondWithError(w http.ResponseWriter, statusCode int, msg string) {
	type Res struct {
		Error string `json:"error"`
	}

	if statusCode > 499 {
		log.Printf("Internal server error: %s", msg)
	}

	RespondWithJson(w, statusCode, Res{Error: msg})
}
