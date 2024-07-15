package main

import (
	"io"
	"log"

	"github.com/Abi-Liu/TextTunnel/internal/client"
)

func main() {
	type Res struct {
		Ok string
	}

	client := client.CreateHttpClient()
	res, err := client.Get("http://localhost:8080/health")
	if err != nil {
		log.Print(err)
	}

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
	}

	log.Print(string(dat))
}
