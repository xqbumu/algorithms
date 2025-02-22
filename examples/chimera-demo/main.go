package main

import (
	"algorithms/examples/chimera-demo/handlers"
	"log"
	"net/http"

	"github.com/matt1484/chimera"
)

func main() {
	api := chimera.NewAPI()
	api.Use(func(req *http.Request, ctx chimera.RouteContext, next chimera.NextFunc) (chimera.ResponseWriter, error) {
		resp, err := next(req)
		return resp, err
	})
	chimera.Get(api, "/test/{path}", handlers.Test)
	chimera.Get(api, "/counter", handlers.NewCounter().DoHandler)

	addr := "0.0.0.0:8000"
	log.Printf("starting server at %s", addr)
	api.Start(addr)
}
