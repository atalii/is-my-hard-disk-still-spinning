package routes

import (
	"fmt"
	"log"
	"net/http"
	"syscall"
)

func UptimeRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")

	uptime, err := uptime()
	if err != nil {
		log.Printf("uptime: %v", err)
		w.Write([]byte(fmt.Sprintf("<div class='err'>%v</div>", err)))
	} else {
		w.Write([]byte(uptime))
	}
}

func uptime() (string, error) {
	var sysinfo syscall.Sysinfo_t
	if err := syscall.Sysinfo(&sysinfo); err != nil {
		return "", err
	}

	uptime_minutes := sysinfo.Uptime / 60
	uptime_hours := uptime_minutes / 60
	uptime_days := uptime_hours / 24

	val := fmt.Sprintf(
		"Online and uninterrupted for %d days, %d hours, and %d minutes.",
		uptime_days, uptime_hours%24, uptime_minutes%60,
	)

	return val, nil
}
