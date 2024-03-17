package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

//go:embed index.html
var index string

func make_route(inner func() string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(405)
		}

		w.Header().Add("Content-Type", "text/html; charset=utf-8")

		txt := inner()
		w.Write([]byte(txt))
	}
}

func as_html(inner func() (*string, *string)) func() string {
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

func commandFmt(cmd string, args ...string) func() (*string, *string) {
	return func() (*string, *string) {
		cmd := exec.Command(cmd, args...)
		out, err := cmd.CombinedOutput()
		val := string(out)

		if err != nil {
			msg := fmt.Sprintf("%s: %v\n\n%s", cmd, err, val)
			return nil, &msg
		}

		return &val, nil
	}
}

func main() {
	uptime := commandFmt("uptime")
	zpoolStatus := commandFmt("zpool", "status")
	caddyStatus := commandFmt("systemctl", "status", "caddy")

	index_route := make_route(func() string {
		return index
	})

	zpool_status_route := make_route(as_html(zpoolStatus))
	caddy_status_route := make_route(as_html(caddyStatus))
	uptime_route := make_route(as_html(uptime))

	http.HandleFunc("/zpool-status", zpool_status_route)
	http.HandleFunc("/caddy-status", caddy_status_route)
	http.HandleFunc("/uptime", uptime_route)
	http.HandleFunc("/", index_route)

	log.Println("Will listen on 127.0.0.1:4525")
	err := http.ListenAndServe("127.0.0.1:4525", nil)
	log.Fatalf("http listener exited: %v", err)
}
