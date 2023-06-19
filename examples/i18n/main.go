//go:generate gotext -srclang=en-US update -out=catalog_gen.go -lang=en-US,zh algorithms/examples/i18n

package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	log.Println(language.Vietnamese)
	// go:generate gotext update -out catalog_gen.go -lang=en,zh
	// Initialize a router and add the path and handler for the homepage.
	mux := chi.NewMux()
	mux.Get("/{locale}", http.HandlerFunc(handleHome))
	// mux.Handle("/generize", http.HandlerFunc(pkg.Generize))

	// Start the HTTP server using the router.
	log.Print("starting server on :4018...")
	err := http.ListenAndServe(":4018", mux)
	log.Fatal(err)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	// Extract the locale from the URL path. This line of code is likely to
	// be different for you if you are using an alternative router.
	locale := chi.URLParam(r, "locale")

	// Declare variable to hold the target language tag.
	var lang language.Tag

	// Use language.MustParse() to assign the appropriate language tag
	// for the locale.
	switch locale {
	case "en-us":
		lang = language.MustParse("en-US")
	case "zh":
		lang = language.MustParse("zh")
	default:
		http.NotFound(w, r)
		return
	}

	// Initialize a message.Printer which uses the target language.
	p := message.NewPrinter(lang)
	log.Println(p)
	// Print the welcome message translated into the target language.
	p.Fprintf(w, "Welcome %s!", r.Host)
}
