package config

import (
	"os"
	"strconv"
	"strings"
)

// ApplyEnvOverrides mutates cfg from YGGDRASIL_* environment variables.
// Used for Docker / CI cluster nodes with fixed identities and static peers.
func ApplyEnvOverrides(cfg *Config) {
	if cfg == nil {
		return
	}
	if v := os.Getenv("YGGDRASIL_NODE_ID"); v != "" {
		cfg.NodeID = v
	}
	if v := os.Getenv("YGGDRASIL_NODE_NAME"); v != "" {
		cfg.NodeName = v
	}
	if v := os.Getenv("YGGDRASIL_API_HOST"); v != "" {
		cfg.APIHost = v
	}
	if v := os.Getenv("YGGDRASIL_API_PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.APIPort = n
		}
	}
	if v := os.Getenv("YGGDRASIL_INTERNAL_HOST"); v != "" {
		cfg.InternalHost = v
	}
	if v := os.Getenv("YGGDRASIL_INTERNAL_PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.InternalPort = n
		}
	}
	if v := os.Getenv("YGGDRASIL_ADVERTISE_HOST"); v != "" {
		cfg.AdvertiseHost = v
	}
	if v := os.Getenv("YGGDRASIL_DISCOVERY_ENABLED"); v != "" {
		cfg.DiscoveryEnabled = parseBool(v)
	}
	if v := os.Getenv("YGGDRASIL_STATIC_PEERS"); v != "" {
		cfg.StaticPeers = splitCSV(v)
	}
	if v := os.Getenv("YGGDRASIL_WEB_UI_DIR"); v != "" {
		cfg.WebUIDir = v
	}
	if v := os.Getenv("YGGDRASIL_WEB_UI_ENABLED"); v != "" {
		cfg.WebUIEnabled = parseBool(v)
	}
}

// EnvTruthy reports whether a YGGDRASIL_* (or any) env var is a truthy flag.
func EnvTruthy(key string) bool {
	return parseBool(os.Getenv(key))
}

func parseBool(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func splitCSV(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
