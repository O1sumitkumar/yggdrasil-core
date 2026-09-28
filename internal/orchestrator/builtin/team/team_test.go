package team

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/yeixio/yggdrasil-core/pkg/contracts"
	"github.com/yeixio/yggdrasil-core/pkg/pluginapi"
)

func TestValidateProfileRequiresEachRole(t *testing.T) {
	orch := New()
	err := orch.ValidateProfile(context.Background(), contracts.AIProfile{
		Roles: []contracts.ModelRole{
			{Role: "coordinator", ModelID: "plan"},
			{Role: "worker", ModelID: "work"},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "reviewer") {
		t.Fatalf("error=%v", err)
	}

	err = orch.ValidateProfile(context.Background(), contracts.AIProfile{
		Roles: []contracts.ModelRole{
			{Role: "coordinator", ModelID: "plan"},
			{Role: "worker", ModelID: "work"},
			{Role: "reviewer", ModelID: "review"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if orch.ID() != "team" || !orch.Capabilities().SupportsTeam {
		t.Fatalf("id=%s caps=%+v", orch.ID(), orch.Capabilities())
	}
}

func TestRunUsesReviewOrFallsBackToWork(t *testing.T) {
	env := &scriptEnv{replies: map[string]string{
		"coordinator": "plan",
		"worker":      "draft",
		"reviewer":    "  ",
	}}
	events := collect(t, New(), env)
	last := events[len(events)-1]
	if last.Type != "agent.completed" || last.Content != "draft" || !last.Done {
		t.Fatalf("final=%+v", last)
	}

	env.replies["reviewer"] = "ship it"
	events = collect(t, New(), env)
	last = events[len(events)-1]
	if last.Content != "ship it" {
		t.Fatalf("final=%+v", last)
	}
}

func TestRunStopsWhenARoleFails(t *testing.T) {
	env := &scriptEnv{
		replies: map[string]string{"coordinator": "plan"},
		fail:    "worker",
	}
	events := collect(t, New(), env)
	last := events[len(events)-1]
	if last.Type != "agent.error" || !last.Done {
		t.Fatalf("final=%+v", last)
	}
}

type scriptEnv struct {
	replies map[string]string
	fail    string
}

func (e *scriptEnv) Generate(ctx context.Context, role string, messages []pluginapi.ChatMessage) (<-chan pluginapi.ChatChunk, error) {
	if role == e.fail {
		return nil, errors.New("worker unavailable")
	}
	ch := make(chan pluginapi.ChatChunk, 1)
	ch <- pluginapi.ChatChunk{Content: e.replies[role], Done: true}
	close(ch)
	return ch, nil
}

func (e *scriptEnv) ExecuteTool(context.Context, string, map[string]any) (map[string]any, error) {
	return nil, nil
}

func (e *scriptEnv) Emit(string, map[string]any) {}

func (e *scriptEnv) NodeForRole(role string) (string, error) {
	return "node-" + role, nil
}

func collect(t *testing.T, orch *Orchestrator, env *scriptEnv) []pluginapi.OrchestrationEvent {
	t.Helper()
	stream, err := orch.Run(context.Background(), contracts.Task{Prompt: "hello"}, contracts.AIProfile{}, env)
	if err != nil {
		t.Fatal(err)
	}
	var events []pluginapi.OrchestrationEvent
	for event := range stream {
		events = append(events, event)
	}
	if len(events) == 0 {
		t.Fatal("no events")
	}
	return events
}
