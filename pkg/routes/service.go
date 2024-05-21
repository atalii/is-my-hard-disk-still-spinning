package routes

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/coreos/go-systemd/v22/dbus"

	"github.com/atalii/is-my-hard-disk-still-spinning/v2/pkg/conf"
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

func ServiceAllRoute(w http.ResponseWriter, r *http.Request) {
	systemdConn.lock.Lock()
	units, err := systemdConn.conn.ListUnitsContext(context.TODO())
	systemdConn.lock.Unlock()

	if err != nil {
		log.Printf("all services route: %v", err)
		w.WriteHeader(500)
		return
	}

	slices.SortFunc(units, func(a dbus.UnitStatus, b dbus.UnitStatus) int {
		return strings.Compare(a.Name, b.Name)
	})

	for _, unit := range units {
		if !strings.HasSuffix(unit.Name, ".service") {
			continue
		}

		tmpl, err := templateService(unit.Name)

		if err != nil {
			log.Printf("all services route templating: %s: %v", unit.Name, err)
			w.WriteHeader(500)
			continue
		}

		if err := tmpl.Execute(w, conf.Service{
			Name: unit.Name,
		}); err != nil {
			log.Printf("all service template execution: %s: %v", unit.Name, err)
			w.WriteHeader(500)
			continue
		}
	}

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
}

func ServiceStatusRoute(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("service")

	tmpl, err := templateService(fmt.Sprintf("%s.service", name))

	if err != nil {
		log.Printf("status route template parsing: %v", err)
		w.WriteHeader(500)
		return
	}

	if err := tmpl.Execute(w, conf.Service{
		Name: name,
	}); err != nil {
		log.Printf("service status route template execution: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
}

func templateService(serviceName string) (*template.Template, error) {
	var tmpl *template.Template

	status, err := status(serviceName)

	if err != nil {
		log.Printf("service route: on %s: %v", serviceName, err)
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
		return nil, err
	} else {
		return tmpl, nil
	}
}

func status(serviceUnit string) (string, error) {
	log.Printf("checking status of unit: %s", serviceUnit)

	systemdConn.lock.Lock()
	defer systemdConn.lock.Unlock()

	prop, err := systemdConn.conn.GetUnitPropertyContext(context.TODO(), serviceUnit, "ActiveState")
	if err != nil {
		return "", fmt.Errorf("can't read status: %w", err)
	}

	return prop.Value.Value().(string), nil
}
