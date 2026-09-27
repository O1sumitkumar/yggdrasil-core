package models

import (
	"fmt"
	"sync"
)

// Catalog loads and indexes curated model manifests.
type Catalog struct {
	mu      sync.RWMutex
	path    string
	entries map[string]CatalogEntry
}

// NewCatalog loads catalog.json from the manifest search path or embedded fallback.
func NewCatalog(dataDir string) (*Catalog, error) {
	path, err := FindManifestPath("catalog.json", dataDir)
	if err != nil {
		return NewCatalogEmbedded()
	}
	var manifest CatalogManifest
	if err := loadJSON(path, &manifest); err != nil {
		return nil, err
	}
	c := &Catalog{
		path:    path,
		entries: make(map[string]CatalogEntry, len(manifest.Models)),
	}
	for _, m := range manifest.Models {
		c.entries[m.ID] = m
	}
	return c, nil
}

// Upsert adds or replaces a catalog entry (used for dynamic HF installs).
func (c *Catalog) Upsert(e CatalogEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.entries == nil {
		c.entries = make(map[string]CatalogEntry)
	}
	c.entries[e.ID] = e
}

// Get returns a catalog entry by ID.
func (c *Catalog) Get(id string) (CatalogEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[id]
	return e, ok
}

// List returns all catalog entries.
func (c *Catalog) List() []CatalogEntry {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]CatalogEntry, 0, len(c.entries))
	for _, e := range c.entries {
		out = append(out, e)
	}
	return out
}

// LoadPresets loads purpose presets.
func LoadPresets(dataDir string) ([]PurposePreset, error) {
	path, err := FindManifestPath("presets.json", dataDir)
	if err != nil {
		return LoadPresetsEmbedded()
	}
	var manifest PresetsManifest
	if err := loadJSON(path, &manifest); err != nil {
		return nil, err
	}
	if len(manifest.Purposes) == 0 {
		return nil, fmt.Errorf("presets manifest is empty")
	}
	return manifest.Purposes, nil
}
