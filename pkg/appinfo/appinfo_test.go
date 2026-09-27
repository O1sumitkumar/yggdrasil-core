package appinfo

import "testing"

func TestDefaultsMatchDaemon(t *testing.T) {
	if DefaultAPIAddr() != "127.0.0.1:7331" {
		t.Fatalf("API address = %q", DefaultAPIAddr())
	}
	if DefaultDataDir() == "" {
		t.Fatal("data directory is empty")
	}
}
