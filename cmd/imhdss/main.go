package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/atalii/is-my-hard-disk-still-spinning/v2/pkg/cmd"
	"github.com/atalii/is-my-hard-disk-still-spinning/v2/pkg/conf"
	"github.com/atalii/is-my-hard-disk-still-spinning/v2/pkg/routes"
)

func makeRoute(inner func() string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
		} else {
			return fmt.Sprintf(
				"<div class=\"err\">%v</div>",
				*err,
			)
		}
	}
}

func main() {
	confPath := flag.String("config", "/etc/imhdss/conf.kdl", "Location of the configuration file.")

	flag.Parse()
	if err := conf.ReadConf(*confPath); err != nil {
		log.Fatalf("while reading configuration: %s: %v", *confPath, err)
	} else {
		run()
	}
}

func run() {
	if err := routes.InitState(); err != nil {
		log.Fatalf("cannot start: %v", err)
	}

	zpool_status_route := makeRoute(asHtml(cmd.Runner("zpool", "status")))

	http.HandleFunc("GET /stats/systemd/{service}", routes.ServiceStatusRoute)
	http.HandleFunc("GET /stats/zpool-status", zpool_status_route)
	http.HandleFunc("GET /stats/uptime", routes.UptimeRoute)
	http.HandleFunc("GET /styles.css", routes.Styles)
	http.HandleFunc("GET /", routes.Index)

	log.Println("Will listen on 127.0.0.1:4525")
	err := http.ListenAndServe("127.0.0.1:4525", nil)
	log.Fatalf("http listener exited: %v", err)
}
