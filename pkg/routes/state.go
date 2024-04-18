package routes

import (
	"fmt"
	"sync"

	"github.com/coreos/go-systemd/v22/dbus"
)

var systemdConn struct{
	lock sync.RWMutex
	conn *dbus.Conn
}

func InitState() error {
	return initSystemdConn()
}


func initSystemdConn() error {
	systemdConn.lock.Lock()
	defer systemdConn.lock.Unlock()

	conn, err := dbus.New()
	if err != nil {
		return fmt.Errorf("failed to initialize state: %w", err)
	}

	systemdConn.conn = conn
	return nil
}
