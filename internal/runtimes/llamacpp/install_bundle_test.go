package llamacpp

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestBundledLlamaServerBeside(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "yggdrasil-daemon")
	if runtime.GOOS != "darwin" {
		if got := bundledLlamaServerBeside(exe); got != "" {
			t.Fatalf("got %s", got)
		}
		return
	}
	bin := filepath.Join(dir, "llamacpp", "llama-server")
	if err := os.MkdirAll(filepath.Dir(bin), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bin, []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := bundledLlamaServerBeside(exe); got != bin {
		t.Fatalf("got %s, want %s", got, bin)
	}
}
