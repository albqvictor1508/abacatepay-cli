package utils

import (
	"net"
	"time"
)

const GoogleDNS = "8.8.8:53"

func IsOnline() bool {
	connection, err := net.DialTimeout("tcp", GoogleDNS, 2*time.Second)

	if err != nil {
		return false
	}

	defer connection.Close()

	return true
}
