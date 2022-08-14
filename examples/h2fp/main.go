package main

import (
	"algorithms/assets"
	"algorithms/pkg/netfake"
	"crypto/tls"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
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

	// server := http.Server{
	// 	Addr:    ":8443",
	// 	Handler: r,
	// }
	// if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
	// 	panic(err)
	// }

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

	ln = tls.NewListener(netfake.NewListener(ln), cfg)

	if err := http.Serve(ln, r); err != nil {
		panic(err)
	}
}
