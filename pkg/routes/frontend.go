package routes

import (
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

type Site struct {
	Url  string
	Name string
}

type Service struct {
	Name string
}

type IndexData struct {
	Sites    []Site
	Services []Service
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

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
				Url:  "https://wiki-js.home.tali.network",
				Name: "wiki.js",
			},
			Site{
				Url:  "https://jf-home.tali.network",
				Name: "jellyfin",
			},
			Site{
				Url:  "https://in-home.tali.network",
				Name: "invoke-ai",
			},
		},
		Services: []Service{
			Service{
				Name: "tailscaled",
			},
			Service{
				Name: "postgresql",
			},
			Service{
				Name: "caddy",
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

	w.Header().Add("Content-Type", "text/css")
	w.Write(styles)
}
