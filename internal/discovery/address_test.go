package discovery

import (
	"net"
	"strconv"
	"testing"
)

func TestJoinHostPortFormat(t *testing.T) {
	got := net.JoinHostPort("192.168.1.50", strconv.Itoa(7332))
	if got != "192.168.1.50:7332" {
		t.Fatalf("got %q", got)
	}
	// IPv6 needs brackets
	got6 := net.JoinHostPort("::1", "7332")
	if got6 != "[::1]:7332" {
		t.Fatalf("got %q", got6)
	}
}
