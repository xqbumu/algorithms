package main

import (
	"net/http"
	"os"
	"text/template"

	"github.com/olivere/vite"
)

func main() {
	viteFragment, err := vite.HTMLFragment(vite.Config{
		FS:        os.DirFS("frontend/dist"),
		IsDev:     *isDev,
		ViteURL:   "http://localhost:5173", // optional: defaults to this
		ViteEntry: "src/main.js",           // reccomended as highly dependent on your app
	})
	if err != nil {
		panic(err)
	}

	indexTemplate := `
  <head>
      <meta charset="UTF-8" />
      <title>My Go Application</title>
      {{ .Vite.Tags }}
  </head>
  <body></body>
  `

	tmpl := template.Must(template.New("name").Parse(indexTemplate))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pageData := map[string]interface{}{
			"Vite": viteFragment,
		}

		if err = tmpl.Execute(w, pageData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
