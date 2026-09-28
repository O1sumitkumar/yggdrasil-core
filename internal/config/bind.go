package config

import (
	"net"
	"strings"
)

// ListensBeyondLoopback reports whether host is not a loopback bind.
// Unspecified addresses such as 0.0.0.0 and :: are beyond loopback.
func ListensBeyondLoopback(host string) bool {
	h := strings.TrimSpace(strings.ToLower(host))
	h = strings.Trim(h, "[]")
	if h == "" || h == "localhost" {
		return false
	}
	ip := net.ParseIP(h)
	if ip == nil {
		return true
	}
	return !ip.IsLoopback()
}
