package app

import (
	"testing"

	"github.com/yeixio/yggdrasil-core/internal/profiles"
	"github.com/yeixio/yggdrasil-core/pkg/contracts"
)

func TestWithChatModelPreservesTeamRolesAndPins(t *testing.T) {
	p := profiles.Profile{
		ID:             "programming",
		OrchestratorID: "team",
		NodePolicy:     contracts.NodePolicy{Mode: "prefer_local"},
		Roles: []contracts.ModelRole{
			{Role: "coordinator", ModelID: "", NodeID: "node-a"},
			{Role: "worker", ModelID: "special-worker", NodeID: "node-b"},
			{Role: "reviewer"},
		},
	}
	got := withChatModel(p, "ui-model")
	if got.OrchestratorID != "team" {
		t.Fatalf("orchestrator=%q, want team", got.OrchestratorID)
	}
	if got.NodePolicy.Mode != "prefer_local" {
		t.Fatalf("node policy=%q, want prefer_local (preserved)", got.NodePolicy.Mode)
	}
	byRole := map[string]contracts.ModelRole{}
	for _, r := range got.Roles {
		byRole[r.Role] = r
	}
	if byRole["coordinator"].ModelID != "ui-model" {
		t.Fatalf("coordinator model=%q, want ui-model", byRole["coordinator"].ModelID)
	}
	if byRole["worker"].ModelID != "special-worker" {
		t.Fatalf("worker model=%q, want special-worker", byRole["worker"].ModelID)
	}
	if byRole["worker"].NodeID != "node-b" {
		t.Fatalf("worker node=%q, want node-b", byRole["worker"].NodeID)
	}
	if byRole["reviewer"].ModelID != "ui-model" {
		t.Fatalf("reviewer model=%q, want ui-model", byRole["reviewer"].ModelID)
	}
	if byRole["coordinator"].NodeID != "node-a" {
		t.Fatalf("coordinator node=%q, want node-a", byRole["coordinator"].NodeID)
	}
}

func TestWithChatModelCollapsesNonTeam(t *testing.T) {
	p := profiles.Profile{
		ID:             "general-assistant",
		OrchestratorID: "simple",
		Roles: []contracts.ModelRole{
			{Role: "assistant", ModelID: "old"},
		},
	}
	got := withChatModel(p, "ui-model")
	if got.OrchestratorID != "simple" {
		t.Fatalf("orchestrator=%q, want simple", got.OrchestratorID)
	}
	if len(got.Roles) != 1 || got.Roles[0].Role != "assistant" || got.Roles[0].ModelID != "ui-model" {
		t.Fatalf("roles=%v, want single assistant/ui-model", got.Roles)
	}
}

func TestProfileNeedsModelFill(t *testing.T) {
	teamEmpty := profiles.Profile{
		OrchestratorID: "team",
		Roles: []contracts.ModelRole{
			{Role: "coordinator"},
			{Role: "worker"},
			{Role: "reviewer"},
		},
	}
	if !profileNeedsModelFill(teamEmpty) {
		t.Fatal("team with empty roles should need fill")
	}
	teamFull := profiles.Profile{
		OrchestratorID: "team",
		Roles: []contracts.ModelRole{
			{Role: "coordinator", ModelID: "a"},
			{Role: "worker", ModelID: "b"},
			{Role: "reviewer", ModelID: "c"},
		},
	}
	if profileNeedsModelFill(teamFull) {
		t.Fatal("team with all models should not need fill")
	}
	teamPartial := profiles.Profile{
		OrchestratorID: "team",
		Roles: []contracts.ModelRole{
			{Role: "coordinator", ModelID: "a"},
			{Role: "worker"},
			{Role: "reviewer", ModelID: "c"},
		},
	}
	if !profileNeedsModelFill(teamPartial) {
		t.Fatal("team missing worker should need fill")
	}
	simpleEmpty := profiles.Profile{
		OrchestratorID: "simple",
		Roles:          []contracts.ModelRole{{Role: "assistant"}},
	}
	if !profileNeedsModelFill(simpleEmpty) {
		t.Fatal("simple empty should need fill")
	}
	simpleOK := profiles.Profile{
		OrchestratorID: "simple",
		Roles:          []contracts.ModelRole{{Role: "assistant", ModelID: "m"}},
	}
	if profileNeedsModelFill(simpleOK) {
		t.Fatal("simple with model should not need fill")
	}
}

func TestChatExecEnvModelForRolePrefersPinnedModel(t *testing.T) {
	env := &chatExecEnv{
		modelOverride: "ui-model",
		profile: profiles.Profile{
			Roles: []contracts.ModelRole{
				{Role: "coordinator", ModelID: "ui-model"},
				{Role: "worker", ModelID: "worker-model"},
			},
		},
	}
	if got := env.modelForRole("worker"); got != "worker-model" {
		t.Fatalf("worker=%q, want worker-model", got)
	}
	if got := env.modelForRole("reviewer"); got != "ui-model" {
		t.Fatalf("missing role fallback=%q, want ui-model", got)
	}
}
