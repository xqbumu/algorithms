package main

import (
	"algorithms/examples/h2fp/assets"
	"algorithms/examples/h2fp/pkg/httpfake"
	"algorithms/examples/h2fp/pkg/netfake"
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
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile | log.Lmicroseconds)

	hander := newHander()

	go func() {
		server := runHttps(hander)
		defer server.Close()
	}()

	go func() {
		server := runHttp2c(hander)
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

	ln, err := net.Listen("tcp", ":9080")
	if err != nil {
		panic(err)
	}

	if err := server.Serve(netfake.NewListener("http", ln)); err != nil {
		panic(err)
	}

	return server
}

func runHttp2c(r http.Handler) http.Server {
	h2s := &http2.Server{}
	h1s := http.Server{
		ReadTimeout:  time.Second * 300,
		WriteTimeout: time.Second * 300,
		Handler:      h2c.NewHandler(r, h2s),
	}

	ln, err := net.Listen("tcp", ":9081")
	if err != nil {
		panic(err)
	}

	if err := h1s.Serve(netfake.NewListener("http", ln)); err != nil {
		panic(err)
	}

	return h1s
}

func runHttps(r http.Handler) *httpfake.Server {
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

	server := &httpfake.Server{
		ReadTimeout:  time.Second * 300,
		WriteTimeout: time.Second * 300,
		Handler: httpfake.HandlerFunc(func(rw httpfake.ResponseWriter, req *httpfake.Request) {
			r.ServeHTTP(httpfake.StdResponseWriter(rw), httpfake.StdRequest(req))
		}),
		TLSConfig: cfg,
	}

	ln, err := net.Listen("tcp", ":9443")
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
	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
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
	req, err := http.NewRequest(http.MethodGet, "https://127.0.0.1:9443", nil)
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
