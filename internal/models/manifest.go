package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// CatalogManifest is the top-level catalog file.
type CatalogManifest struct {
	Version int            `json:"version"`
	Models  []CatalogEntry `json:"models"`
}

// PresetsManifest holds purpose presets.
type PresetsManifest struct {
	Version  int             `json:"version"`
	Purposes []PurposePreset `json:"purposes"`
}

// FindManifestPath searches for a manifest file relative to cwd and data dir.
func FindManifestPath(name string, dataDir string) (string, error) {
	candidates := []string{
		filepath.Join("manifests", "models", name),
	}
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			filepath.Join(wd, "manifests", "models", name),
			filepath.Join(wd, "..", "manifests", "models", name),
		)
	}
	if dataDir != "" {
		candidates = append(candidates, filepath.Join(dataDir, "manifests", "models", name))
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return c, nil
		}
	}
	return "", fmt.Errorf("manifest %s not found", name)
}

func loadJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}
