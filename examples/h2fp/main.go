package main

import (
	"algorithms/assets"
	"algorithms/pkg/netfake"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile | log.Lmicroseconds)
	var (
		err error
	)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	certFile, err := assets.FS.ReadFile("certs/_wildcard.example.arpa.pem")
	if err != nil {
		panic(err)
	}
	keyFile, err := assets.FS.ReadFile("certs/_wildcard.example.arpa-key.pem")
	if err != nil {
		panic(err)
	}

	server := http.Server{
		ReadTimeout:  time.Second * 300,
		WriteTimeout: time.Second * 300,
		Handler:      r,
	}

	cfg := &tls.Config{}
	cfg.Certificates = make([]tls.Certificate, 1)
	cfg.Certificates[0], err = tls.X509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	ln, err := net.Listen("tcp", ":8443")
	if err != nil {
		panic(err)
	}

	ln = tls.NewListener(netfake.NewListener(ln, "server"), cfg)

	go func() {
		if err := server.Serve(ln); err != nil {
			panic(err)
		}
	}()
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, "https://127.0.0.1:8443", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Host", "www.example.arpa")

	client := newClient()
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer func() {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		log.Printf("body length: %s", body)
	}()
}

func newClient() *http.Client {
	return &http.Client{
		Transport: netfake.NewTransport(),
		Timeout:   time.Second * 300,
	}
}
