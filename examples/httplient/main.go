package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	client := http.Client{}

	// url := "https://logs-prod3.grafana.net/loki/api/v1/rules"
	url := "https://arkjit.grafana.net/api/orgs/arkjit/instances"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer glsa_OTPWQaushoD5cCBNh4CpYNrMmpyrvco0_6b13085e")
	// req.SetBasicAuth("316475", "eyJrIjoiZWU5ZTRjM2EwMmYzZjkwMzlmZWMwYzU4ZmVjNGVmMDIzYzY5NTU2OSIsIm4iOiJhaXJkYi1zY291dCIsImlkIjo3MzcxMDN9")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	log.Printf("%s", bytes)
}
