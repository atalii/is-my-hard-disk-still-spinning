package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/atalii/is-my-hard-disk-still-spinning/v2/pkg/cmd"
	"github.com/atalii/is-my-hard-disk-still-spinning/v2/pkg/routes"
	"github.com/atalii/is-my-hard-disk-still-spinning/v2/pkg/sys"
)

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
				"<div class=\"val\">%v</div>",
				*val,
			)
		}

		if err != nil {
			return fmt.Sprintf(
				"<div class=\"err\">%v</div>",
				*err,
			)
		}

		panic("unreachable")
	}
}

func main() {
	if err := routes.InitState(); err != nil {
		log.Fatalf("cannot start: %v", err)
	}

	uptime := sys.Uptime()

	zpool_status_route := makeRoute(asHtml(cmd.Runner("zpool", "status")))
	uptime_route := makeRoute(asHtml(uptime))

	http.HandleFunc("/stats/systemd/{service}", routes.ServiceStatusRoute)
	http.HandleFunc("/stats/zpool-status", zpool_status_route)
	http.HandleFunc("/stats/uptime", uptime_route)
	http.HandleFunc("/styles.css", routes.Styles)
	http.HandleFunc("/", routes.Index)

	log.Println("Will listen on 127.0.0.1:4525")
	err := http.ListenAndServe("127.0.0.1:4525", nil)
	log.Fatalf("http listener exited: %v", err)
}
