//go:build windows

package llamacpp

import (
	"os"
	"os/exec"
	"strconv"
	"time"
)

func configureLlamaProcess(cmd *exec.Cmd) {}

func stopLlamaProcess(cmd *exec.Cmd, graceful time.Duration) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	pid := cmd.Process.Pid
	_ = cmd.Process.Signal(os.Interrupt)
	done := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(done)
	}()
	select {
	case <-done:
		return
	case <-time.After(graceful):
	}
	_ = exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/T", "/F").Run()
	select {
	case <-done:
	case <-time.After(800 * time.Millisecond):
	}
}
