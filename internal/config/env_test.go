package config

import "testing"

func TestSplitCSVAndApplyEnv(t *testing.T) {
	t.Setenv("YGGDRASIL_NODE_ID", "nid-1")
	t.Setenv("YGGDRASIL_NODE_NAME", "Node One")
	t.Setenv("YGGDRASIL_API_HOST", "0.0.0.0")
	t.Setenv("YGGDRASIL_API_PORT", "7331")
	t.Setenv("YGGDRASIL_INTERNAL_HOST", "0.0.0.0")
	t.Setenv("YGGDRASIL_INTERNAL_PORT", "7332")
	t.Setenv("YGGDRASIL_ADVERTISE_HOST", "node-a")
	t.Setenv("YGGDRASIL_DISCOVERY_ENABLED", "true")
	t.Setenv("YGGDRASIL_STATIC_PEERS", "node-b:7332, node-c:7332")
	t.Setenv("YGGDRASIL_WEB_UI_DIR", "/opt/yggdrasil/web")
	t.Setenv("YGGDRASIL_WEB_UI_ENABLED", "true")

	cfg := DefaultConfig()
	ApplyEnvOverrides(&cfg)
	if cfg.NodeID != "nid-1" || cfg.NodeName != "Node One" {
		t.Fatalf("identity: %+v", cfg)
	}
	if cfg.APIHost != "0.0.0.0" || cfg.InternalHost != "0.0.0.0" {
		t.Fatalf("hosts: %+v", cfg)
	}
	if cfg.AdvertiseHost != "node-a" {
		t.Fatalf("advertise=%q", cfg.AdvertiseHost)
	}
	if len(cfg.StaticPeers) != 2 || cfg.StaticPeers[0] != "node-b:7332" {
		t.Fatalf("peers=%v", cfg.StaticPeers)
	}
	if cfg.WebUIDir != "/opt/yggdrasil/web" || !cfg.WebUIEnabled {
		t.Fatalf("web ui: dir=%q enabled=%v", cfg.WebUIDir, cfg.WebUIEnabled)
	}
}
