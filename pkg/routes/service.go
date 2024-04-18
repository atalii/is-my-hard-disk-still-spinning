package routes

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

//go:embed static/service.up.html
var svUp string

//go:embed static/service.err.html
var svErr string

//go:embed static/service.inactive.html
var svIn string

//go:embed static/service.reloading.html
var svRl string

//go:embed static/service.activating.html
var svAe string

//go:embed static/service.deactivating.html
var svDe string

//go:embed static/service.failed.html
var svFl string

func ServiceStatusRoute(w http.ResponseWriter, r *http.Request) {
	serviceName := r.PathValue("service")

	var tmpl *template.Template

	status, err := status(fmt.Sprintf("%s.service", serviceName))
	if err != nil {
		log.Printf("service route: %v", err)
		tmpl, err = template.New("serviceErr").Parse(svErr)
	} else {
		switch status {
		case "reloading":
			tmpl, err = template.New("serviceRl").Parse(svRl)
		case "active":
			tmpl, err = template.New("serviceUp").Parse(svUp)
		case "inactive":
			tmpl, err = template.New("serviceIn").Parse(svIn)
		case "activating":
			tmpl, err = template.New("serviceAe").Parse(svDe)
		case "deactivating":
			tmpl, err = template.New("serviceDe").Parse(svDe)
		case "failed":
			tmpl, err = template.New("serviceFl").Parse(svFl)
		default:
			log.Printf("unknown ActiveState value: %s", status)
			tmpl, err = template.New("serviceErr").Parse(svErr)
		}
	}

	if err != nil {
		log.Printf("status route template parsing: %v", err)
		w.WriteHeader(500)
		return
	}

	if err := tmpl.Execute(w, Service{
		Name: serviceName,
	}); err != nil {
		log.Printf("service.up.html: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
}

func status(serviceUnit string) (string, error) {
	systemdConn.lock.Lock()
	defer systemdConn.lock.Unlock()

	prop, err := systemdConn.conn.GetUnitPropertyContext(context.TODO(), serviceUnit, "ActiveState")
	if err != nil {
		return "", fmt.Errorf("can't read status: %w", err)
	}

	return prop.Value.Value().(string), nil
}
