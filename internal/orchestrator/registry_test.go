package orchestrator

import (
	"testing"

	"github.com/yeixio/yggdrasil-core/internal/orchestrator/builtin/simple"
	"github.com/yeixio/yggdrasil-core/internal/orchestrator/builtin/team"
	"github.com/yeixio/yggdrasil-core/pkg/pluginapi"
)

func TestRegistryReplacesWithoutDuplicating(t *testing.T) {
	registry := NewRegistry()
	registry.Register(simple.New())
	registry.Register(team.New())
	registry.Register(simple.New())

	listed := registry.List()
	if len(listed) != 2 || listed[0].ID() != "simple" || listed[1].ID() != "team" {
		t.Fatalf("list=%v", ids(listed))
	}
	got, err := registry.Get("team")
	if err != nil || got.DisplayName() != "Team" {
		t.Fatalf("get=%v err=%v", got, err)
	}
	if _, err := registry.Get("missing"); err == nil {
		t.Fatal("missing orchestrator was returned")
	}
}

func ids(items []pluginapi.Orchestrator) []string {
	out := make([]string, len(items))
	for i, item := range items {
		out[i] = item.ID()
	}
	return out
}
