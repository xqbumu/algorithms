package main

import (
	"algorithms/assets"
	"algorithms/pkg/netfake"
	"crypto/tls"
	"html/template"
	"io"
	"io/fs"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile | log.Lmicroseconds)

	hander := newHander()

	go func() {
		server := runHttps(hander)
		defer server.Close()
	}()

	server := runHttp(hander)
	defer server.Close()
}

func runHttp(r http.Handler) http.Server {
	server := http.Server{
		ReadTimeout:  time.Second * 300,
		WriteTimeout: time.Second * 300,
		Handler:      r,
	}

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	if err := server.Serve(netfake.NewListener("http", ln)); err != nil {
		panic(err)
	}

	return server
}

func runHttps(r http.Handler) *http.Server {
	certFile, err := assets.FS.ReadFile("certs/_wildcard.example.arpa.pem")
	if err != nil {
		panic(err)
	}
	keyFile, err := assets.FS.ReadFile("certs/_wildcard.example.arpa-key.pem")
	if err != nil {
		panic(err)
	}

	cfg := &tls.Config{}
	cfg.Certificates = make([]tls.Certificate, 1)
	cfg.Certificates[0], err = tls.X509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		ReadTimeout:  time.Second * 300,
		WriteTimeout: time.Second * 300,
		Handler:      r,
		TLSConfig:    cfg,
	}

	ln, err := net.Listen("tcp", ":8443")
	if err != nil {
		panic(err)
	}

	if err := server.ServeTLS(netfake.NewListener("https", ln), "", ""); err != nil {
		panic(err)
	}

	return server
}

func newHander() http.Handler {
	tpls := template.New("")
	tpls.Funcs(template.FuncMap{
		"br": func(i, j int) bool {
			if i == 0 {
				return false
			} else {
				return i%j == (j - 1)
			}
		},
	})
	tpls, err := tpls.ParseFS(assets.FS, "views/*.tpl")
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Get("/flags", func(w http.ResponseWriter, r *http.Request) {
		flags, err := fs.Glob(assets.FS, "static/flags/*.png")
		if err != nil {
			panic(err)
		}
		if err := tpls.ExecuteTemplate(w, "flags.tpl", map[string]any{
			"flags": flags,
		}); err != nil {
			panic(err)
		}
	})
	r.Handle("/static/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
		http.FileServer(http.FS(assets.FS)).ServeHTTP(w, r)
	}))

	return r
}

func doRequest() {
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
