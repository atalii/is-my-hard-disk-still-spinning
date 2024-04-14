package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"

	"github.com/atalii/is-my-hard-disk-still-spinning/v2/pkg/cmd"
	"github.com/atalii/is-my-hard-disk-still-spinning/v2/pkg/sys"
)

//go:embed index.html
var index string

func makeRoute(inner func() string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(405)
		}

		w.Header().Add("Content-Type", "text/html; charset=utf-8")

		txt := inner()
		w.Write([]byte(txt))
	}
}

func asHtml(inner func() (*string, *string)) func() string {
	return func() string {
		val, err := inner()

		if val != nil {
			return fmt.Sprintf(
				"<div class=\"val\"><pre>%v</pre></div>",
				*val,
			)
		}

		if err != nil {
			return fmt.Sprintf(
				"<div class=\"err\"><pre>%v</pre></div>",
				*err,
			)
		}

		panic("unreachable")
	}
}

func main() {
	uptime := sys.Uptime()
	zpoolStatus := cmd.Runner("zpool", "status")
	caddyStatus := cmd.Runner("systemctl", "status", "caddy")

	index_route := makeRoute(func() string {
		return index
	})

	zpool_status_route := makeRoute(asHtml(zpoolStatus))
	caddy_status_route := makeRoute(asHtml(caddyStatus))
	uptime_route := makeRoute(asHtml(uptime))

	http.HandleFunc("/zpool-status", zpool_status_route)
	http.HandleFunc("/caddy-status", caddy_status_route)
	http.HandleFunc("/uptime", uptime_route)
	http.HandleFunc("/", index_route)

	log.Println("Will listen on 127.0.0.1:4525")
	err := http.ListenAndServe("127.0.0.1:4525", nil)
	log.Fatalf("http listener exited: %v", err)
}
