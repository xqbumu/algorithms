package main

import (
	"algorithms/pkg/netfake"
	"io"
	"log"
	"net/http"
	"os"
)

var client = &http.Client{}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("please input a https url follow with exe to handshake.")
	}
	req, err := http.NewRequest(http.MethodGet, os.Args[1], nil)
	if err != nil {
		panic(err)
	}
	client.Transport = netfake.NewTransport()

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	log.Printf("Response body with len: %d", len(body))
}
