package hardware

import (
	"context"
	"runtime"
	"testing"
)

func TestDetectReturnsBasics(t *testing.T) {
	d := &Detector{DataPath: t.TempDir()}
	inv, err := d.Detect(context.Background())
	if err != nil {
		t.Fatalf("detect: %v", err)
	}
	if inv.OS != runtime.GOOS {
		t.Fatalf("os: got %s want %s", inv.OS, runtime.GOOS)
	}
	if inv.Arch != runtime.GOARCH {
		t.Fatalf("arch: got %s want %s", inv.Arch, runtime.GOARCH)
	}
	if inv.CPU.Cores <= 0 {
		t.Fatalf("expected CPU cores > 0")
	}
	if inv.Memory.TotalBytes == 0 {
		t.Fatalf("expected memory total > 0")
	}
	if inv.Disk.AvailableBytes == 0 && inv.Disk.TotalBytes == 0 {
		t.Log("disk totals unavailable; tolerated")
	}
	if len(inv.Accelerators) == 0 {
		t.Fatalf("expected at least one accelerator entry")
	}
}
