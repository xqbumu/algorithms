package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
	"os"
	"time"
)

func main() {
	server := http.NewServeMux()

	server.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, time.Now().Format(time.RFC3339))
		if r.URL.Query().Get("env") == "true" {
			for _, line := range os.Environ() {
				fmt.Fprintln(w, line)
			}
		}
	})

	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, line := range os.Environ() {
			fmt.Fprintln(w, line)
		}
	})

	cgi.Serve(server)
}
