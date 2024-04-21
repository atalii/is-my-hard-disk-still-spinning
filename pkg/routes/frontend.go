package routes

import (
	"github.com/atalii/is-my-hard-disk-still-spinning/v2/pkg/conf"
	"html/template"
	"log"
	"net/http"

	_ "embed"
)

//go:generate tailwindcss -m -i static/input.css -o static/styles.css -c tailwind.config.js

//go:embed static/index.html
var index string

//go:embed static/styles.css
var styles []byte

type IndexData struct {
	Links    []conf.Link
	Services []conf.Service
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.New("index").Parse(index)
	if err != nil {
		log.Printf("couldn't build index template: %v", err)
		w.WriteHeader(405)
		return
	}

	if err := tmpl.Execute(w, IndexData{
		Links:    conf.Links(),
		Services: conf.Services(),
	}); err != nil {
		log.Printf("index template execution failed: %v", err)
	}
}

func Styles(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/css")
	w.Write(styles)
}
