package models

import "github.com/yeixio/yggdrasil-core/pkg/contracts"

// CatalogEntry is a model in the curated manifest.
type CatalogEntry struct {
	ID                string                      `json:"id"`
	DisplayName       string                      `json:"display_name"`
	Summary           string                      `json:"summary,omitempty"`
	Family            string                      `json:"family,omitempty"`
	Variant           string                      `json:"variant,omitempty"`
	Parameters        string                      `json:"parameters,omitempty"`
	SizeBytes         uint64                      `json:"size_bytes,omitempty"`
	MemoryNeededBytes uint64                      `json:"memory_needed_bytes,omitempty"`
	Context           int                         `json:"context,omitempty"`
	Capabilities      contracts.ModelCapabilities `json:"capabilities"`
	Source            contracts.ModelSource       `json:"source"`
	Purpose           []string                    `json:"purpose,omitempty"`
	Tags              []string                    `json:"tags,omitempty"`
	Runtime           []string                    `json:"runtime,omitempty"`
	RecommendedRoles  []string                    `json:"recommended_roles,omitempty"`
	Dynamic           bool                        `json:"dynamic,omitempty"`
}

// PurposePreset describes a use-case preset.
type PurposePreset struct {
	ID              string   `json:"id"`
	Label           string   `json:"label"`
	Description     string   `json:"description"`
	PreferredModels []string `json:"preferred_models"`
}

// DownloadState tracks an in-progress download.
type DownloadState struct {
	ID              string
	ModelID         string
	Status          string
	BytesDownloaded uint64
	BytesTotal      uint64
	TempPath        string
	Error           string
}
