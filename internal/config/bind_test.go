package config

import "testing"

func TestListensBeyondLoopback(t *testing.T) {
	loopback := []string{"", "localhost", "127.0.0.1", "::1", "[::1]"}
	for _, host := range loopback {
		if ListensBeyondLoopback(host) {
			t.Fatalf("%q should be loopback", host)
		}
	}
	remote := []string{"0.0.0.0", "::", "192.168.1.20", "node-a"}
	for _, host := range remote {
		if !ListensBeyondLoopback(host) {
			t.Fatalf("%q should require authentication", host)
		}
	}
}
