package config

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	DefaultAPIPort      = 7331
	DefaultInternalPort = 7332
	DefaultBindLoopback = "127.0.0.1"
	ServiceType         = "_localai._tcp"
	ServiceDomain       = "local."
)

// DefaultDataDir returns the OS-appropriate application data directory.
func DefaultDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Yggdrasil")
	case "windows":
		if local := os.Getenv("LOCALAPPDATA"); local != "" {
			return filepath.Join(local, "Yggdrasil")
		}
		return filepath.Join(home, "AppData", "Local", "Yggdrasil")
	default:
		if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
			return filepath.Join(xdg, "yggdrasil")
		}
		return filepath.Join(home, ".local", "share", "yggdrasil")
	}
}

// DefaultConfig returns a Config with safe local-first defaults.
func DefaultConfig() Config {
	dataDir := DefaultDataDir()
	return Config{
		DataDir:          dataDir,
		DBPath:           filepath.Join(dataDir, "yggdrasil.db"),
		ModelsDir:        filepath.Join(dataDir, "models"),
		RuntimesDir:      filepath.Join(dataDir, "runtimes"),
		LogsDir:          filepath.Join(dataDir, "logs"),
		WebUIDir:         "",
		APIHost:          DefaultBindLoopback,
		APIPort:          DefaultAPIPort,
		InternalHost:     DefaultBindLoopback,
		InternalPort:     DefaultInternalPort,
		LANAPIEnabled:    false,
		WebUIEnabled:     true,
		DiscoveryEnabled: true,
		NodeName:         hostnameOr("Yggdrasil-Node"),
	}
}

func hostnameOr(fallback string) string {
	h, err := os.Hostname()
	if err != nil || h == "" {
		return fallback
	}
	return h
}
