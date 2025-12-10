package utils

import (
	"net"
	"time"
)

var googleDNS = "8.8.8:53"

func IsOnline() bool {
	timeout := 2 * time.Second

	conn, err := net.DialTimeout("tcp", googleDNS, timeout)
	if err != nil {
		return false
	}

	defer conn.Close()
	return true
}
