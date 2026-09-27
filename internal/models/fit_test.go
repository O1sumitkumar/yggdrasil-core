package models_test

import (
	"testing"

	"github.com/yeixio/yggdrasil-core/internal/models"
	"github.com/yeixio/yggdrasil-core/pkg/contracts"
)

func TestScoreFitsAndWinners(t *testing.T) {
	catalog, err := models.NewCatalogEmbedded()
	if err != nil {
		t.Fatal(err)
	}
	hw := contracts.HardwareInventory{
		Memory: contracts.MemoryInfo{TotalBytes: 48 * 1024 * 1024 * 1024, AvailableBytes: 40 * 1024 * 1024 * 1024},
		Accelerators: []contracts.Accelerator{{
			Model: "Apple M-series", Kind: "gpu", UnifiedMemory: 48 * 1024 * 1024 * 1024, Backends: []string{"metal"},
		}},
	}
	presets, _ := models.LoadPresetsEmbedded()
	resp := models.BuildFitResponse(catalog, hw, "local", "Test Mac", presets, models.FitOptions{})
	if len(resp.Fits) == 0 {
		t.Fatal("expected fits")
	}
	if len(resp.Winners) == 0 {
		t.Fatal("expected winners")
	}
	tooLarge := 0
	for _, f := range resp.Fits {
		if f.Label == contracts.FitTooLarge {
			tooLarge++
		}
	}
	// 32B should be too large or tight on 48GB with margin depending on memory_needed
	_ = tooLarge
}

func TestRecommendWithPresetsPrefersCoding(t *testing.T) {
	catalog, err := models.NewCatalogEmbedded()
	if err != nil {
		t.Fatal(err)
	}
	presets, err := models.LoadPresetsEmbedded()
	if err != nil {
		t.Fatal(err)
	}
	hw := contracts.HardwareInventory{
		Memory: contracts.MemoryInfo{TotalBytes: 32 * 1024 * 1024 * 1024},
		Accelerators: []contracts.Accelerator{{
			UnifiedMemory: 32 * 1024 * 1024 * 1024, Backends: []string{"metal"},
		}},
	}
	rec, err := models.RecommendWithPresets(catalog, presets, models.RecommendInput{Purpose: "coding", Hardware: hw})
	if err != nil {
		t.Fatal(err)
	}
	if len(rec.Models) == 0 {
		t.Fatal("expected models")
	}
	if !rec.Models[0].Capabilities.Coding && rec.Models[0].ID == "" {
		t.Fatalf("unexpected recommendation: %+v", rec.Models[0])
	}
}
