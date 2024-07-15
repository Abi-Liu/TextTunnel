package client

import (
	"net/http"
	"time"
)

func CreateHttpClient() *http.Client {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	return client
}
