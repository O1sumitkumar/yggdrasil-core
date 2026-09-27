package orchestrator

import (
	"context"
	"testing"

	"github.com/yeixio/yggdrasil-core/internal/orchestrator/builtin/simple"
	"github.com/yeixio/yggdrasil-core/internal/orchestrator/builtin/team"
	"github.com/yeixio/yggdrasil-core/pkg/contracts"
)

func TestSimpleValidateProfile(t *testing.T) {
	o := simple.New()
	if err := o.ValidateProfile(context.Background(), contracts.AIProfile{
		Name: "x", Roles: []contracts.ModelRole{{Role: "assistant", ModelID: "m1"}},
	}); err != nil {
		t.Fatal(err)
	}
	if err := o.ValidateProfile(context.Background(), contracts.AIProfile{Name: "x"}); err == nil {
		t.Fatal("expected error")
	}
}

func TestTeamValidateProfile(t *testing.T) {
	o := team.New()
	p := contracts.AIProfile{
		Name: "team",
		Roles: []contracts.ModelRole{
			{Role: "coordinator", ModelID: "a"},
			{Role: "worker", ModelID: "b"},
			{Role: "reviewer", ModelID: "c"},
		},
	}
	if err := o.ValidateProfile(context.Background(), p); err != nil {
		t.Fatal(err)
	}
	p.Roles = p.Roles[:2]
	if err := o.ValidateProfile(context.Background(), p); err == nil {
		t.Fatal("expected missing reviewer error")
	}
}
