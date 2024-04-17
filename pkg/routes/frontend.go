package routes

import (
	"log"
	"net/http"
	"html/template"

	_ "embed"
)

//go:generate tailwindcss -m -i static/input.css -o static/styles.css -c tailwind.config.js

//go:embed static/index.html
var index string

//go:embed static/styles.css
var styles []byte

type Site struct {
	Url  string
	Name string
}

type IndexData struct{
	Sites []Site
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}

	tmpl, err := template.New("index").Parse(index)
	if err != nil {
		log.Printf("couldn't build index template: %v", err)
		w.WriteHeader(405)
		return
	}

	if err := tmpl.Execute(w, IndexData{
		Sites: []Site{
			Site{
				Url: "https://wiki-js.home.tali.network",
				Name: "wiki.js",
			},
			Site{
				Url: "https://jf-home.tali.network",
				Name: "jellyfin",
			},
			Site{
				Url: "https://in-home.tali.network",
				Name: "invoke-ai",
			},
		},
	}); err != nil {
		log.Printf("index template execution failed: %v", err)
	}
}

func Styles(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
	}

	w.Header().Add("Content-Type", "text/css");
	w.Write(styles)
}
