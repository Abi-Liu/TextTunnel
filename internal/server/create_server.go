package server

import (
	"net/http"
	"time"
)

func CreateServer() (*http.Server, error) {
	config, err := CreateConfig()
	if err != nil {
		return nil, err
	}

	r := config.NewRouter()

	server := &http.Server{
		Addr:              ":" + config.Env.PORT,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return server, nil
}
