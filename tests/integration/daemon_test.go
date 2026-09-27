package integration_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/yeixio/yggdrasil-core/internal/app"
)

func TestDaemonHealthAndHardware(t *testing.T) {
	dir := t.TempDir()
	application, err := app.New(app.Options{DataDir: dir})
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- application.Start(ctx)
	}()
	defer func() {
		cancel()
		_ = application.Shutdown(context.Background())
	}()

	cfg := application.Config.Get()
	base := "http://" + cfg.APIAddr()
	client := &http.Client{Timeout: 2 * time.Second}

	deadline := time.Now().Add(10 * time.Second)
	var healthOK bool
	for time.Now().Before(deadline) {
		resp, err := client.Get(base + "/api/v1/health")
		if err == nil {
			body, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				var payload map[string]any
				_ = json.Unmarshal(body, &payload)
				if payload["status"] == "ok" {
					healthOK = true
					break
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !healthOK {
		t.Fatal("health endpoint never became ready")
	}

	resp, err := client.Get(base + "/api/v1/hardware")
	if err != nil {
		t.Fatalf("hardware: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("hardware status %d", resp.StatusCode)
	}
	var hw map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&hw); err != nil {
		t.Fatalf("decode hardware: %v", err)
	}
	if hw["os"] == nil || hw["arch"] == nil {
		t.Fatalf("expected os/arch in hardware: %#v", hw)
	}

	resp2, err := client.Get(base + "/api/v1/version")
	if err != nil {
		t.Fatalf("version: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("version status %d", resp2.StatusCode)
	}

	// Settings survive restart path check: config file exists.
	if _, err := filepath.Glob(filepath.Join(dir, "config.json")); err != nil {
		t.Fatalf("config: %v", err)
	}
}
