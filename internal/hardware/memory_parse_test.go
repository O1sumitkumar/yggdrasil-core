package hardware

import "testing"

func TestParseVMStatAvailable(t *testing.T) {
	text := `Mach Virtual Memory Statistics: (page size of 16384 bytes)
Pages free:                               1000.
Pages active:                            20000.
Pages inactive:                           2000.
Pages speculative:                         500.
Pages wired down:                         8000.
`
	got, ok := parseVMStatAvailable(text)
	if !ok {
		t.Fatal("expected parse")
	}
	want := uint64(3500 * 16384)
	if got != want {
		t.Fatalf("got %d want %d", got, want)
	}
}

func TestParseSwapUsage(t *testing.T) {
	total, used, ok := parseSwapUsage("total = 2048.00M  used = 512.00M  free = 1536.00M")
	if !ok {
		t.Fatal("expected parse")
	}
	if total != 2048*1024*1024 || used != 512*1024*1024 {
		t.Fatalf("total %d used %d", total, used)
	}
}
