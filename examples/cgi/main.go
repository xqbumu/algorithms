/*
ExpiresByType text/html "access plus 0 seconds"
ExpiresByType text/plain "access plus 0 seconds"

Options +ExecCGI
AddHandler cgi-script .cgi

RewriteEngine On
RewriteBase /
RewriteCond %{REQUEST_URI} !^/static.*
RewriteCond %{REQUEST_URI} ^([a-zA-z0-9]+)/(.*)$
RewriteRule (.*) /cgi-bin/%1.cgi/%2 [L]
RewriteRule (.*) /cgi-bin/restdb.cgi/$1 [L]
*/
package main

import (
	"algorithms/pkg/restdb"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cgi"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

func main() {
	mux := chi.NewMux()
	mux.Use(
		// Timeout(time.Second*0),
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("x-cgi-prefix", prefix())
				next.ServeHTTP(w, r)
			})
		},
	)

	mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "url: %+v\n", r.URL)
		fmt.Fprintf(w, "path: %+v\n", r.URL.Path)
		fmt.Fprintf(w, "not found")
	})

	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, time.Now().Format(time.RFC3339))
	})

	mux.Get("/env", func(w http.ResponseWriter, r *http.Request) {
		printEnv(w)
	})

	mux.Route(prefix(), func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "prefix: %s\n", prefix())
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(w, "pwd err: %s\n", err.Error())
			} else {
				fmt.Fprintf(w, "pwd: %s\n", cwd)
			}
			fmt.Fprintf(w, "time: %s\n", time.Now().Format(time.RFC3339))
		})

		r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
			if err := os.Chdir(path.Join("..", "etc")); err != nil {
				fmt.Fprint(w, err.Error())
				return
			}
			confPath := r.Header.Get("x-cgi-config")
			if len(confPath) == 0 {
				confPath = "restdb.json"
			}
			confBytes, err := os.ReadFile(confPath)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, err.Error())
				return
			}

			app, err := restdb.NewApp(confBytes)
			if err != nil {
				w.WriteHeader(http.StatusNotImplemented)
				fmt.Fprint(w, err.Error())
				return
			}

			app.Handler(w, r.WithContext(restdb.WithPrefix(r.Context(), prefix())))
		})
	})

	cgi.Serve(mux)
}

func Timeout(timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer func() {
				cancel()
				if ctx.Err() == context.DeadlineExceeded {
					os.Exit(http.StatusGatewayTimeout)
				}
			}()

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func prefix() string {
	filename := filepath.Base(os.Args[0])
	return fmt.Sprintf("/%s", strings.TrimSuffix(filename, filepath.Ext(filename)))
}

func printEnv(w io.Writer) {
	for _, line := range os.Environ() {
		fmt.Fprintln(w, line)
	}
}
