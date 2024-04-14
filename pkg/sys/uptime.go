package sys

import (
	"fmt"
	"syscall"
)

func Uptime() func() (*string, *string) {
	return func() (*string, *string) {
		var sysinfo syscall.Sysinfo_t
		if err := syscall.Sysinfo(&sysinfo); err != nil {
			err := err.Error()
			return nil, &err
		}

		uptime_minutes := sysinfo.Uptime / 60
		uptime_hours   := uptime_minutes / 60
		uptime_days    := uptime_hours / 24

		val := fmt.Sprintf(
			"Up %d days, %d hours, and %d minutes.",
			uptime_days, uptime_hours % 24, uptime_minutes % 60,
		)

		return &val, nil
	}
}
