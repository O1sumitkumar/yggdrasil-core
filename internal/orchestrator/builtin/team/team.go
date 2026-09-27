package team

import (
	"context"
	"fmt"
	"strings"

	"github.com/yeixio/yggdrasil-core/internal/events"
	"github.com/yeixio/yggdrasil-core/pkg/contracts"
	"github.com/yeixio/yggdrasil-core/pkg/pluginapi"
)

const id = "team"

// Orchestrator runs coordinator → worker → reviewer pipeline.
type Orchestrator struct{}

func New() *Orchestrator { return &Orchestrator{} }

func (o *Orchestrator) ID() string          { return id }
func (o *Orchestrator) DisplayName() string { return "Team" }

func (o *Orchestrator) Capabilities() pluginapi.OrchestratorCapabilities {
	return pluginapi.OrchestratorCapabilities{
		SupportsTools: true,
		SupportsTeam:  true,
		Roles:         []string{"coordinator", "worker", "reviewer"},
	}
}

func (o *Orchestrator) ValidateProfile(ctx context.Context, profile contracts.AIProfile) error {
	required := map[string]bool{"coordinator": false, "worker": false, "reviewer": false}
	for _, r := range profile.Roles {
		if _, ok := required[r.Role]; ok && r.ModelID != "" {
			required[r.Role] = true
		}
	}
	for role, ok := range required {
		if !ok {
			return fmt.Errorf("team orchestrator requires role %q with model_id", role)
		}
	}
	return nil
}

func (o *Orchestrator) Run(
	ctx context.Context,
	task contracts.Task,
	profile contracts.AIProfile,
	env pluginapi.ExecutionEnvironment,
) (<-chan pluginapi.OrchestrationEvent, error) {
	ch := make(chan pluginapi.OrchestrationEvent, 16)
	go func() {
		defer close(ch)
		plan, err := o.runRole(ctx, env, ch, "coordinator", fmt.Sprintf("Plan how to answer: %s", task.Prompt))
		if err != nil {
			ch <- pluginapi.OrchestrationEvent{Type: "agent.error", Error: err.Error(), Done: true}
			return
		}
		work, err := o.runRole(ctx, env, ch, "worker", fmt.Sprintf("Execute plan:\n%s\n\nTask: %s", plan, task.Prompt))
		if err != nil {
			ch <- pluginapi.OrchestrationEvent{Type: "agent.error", Error: err.Error(), Done: true}
			return
		}
		review, err := o.runRole(ctx, env, ch, "reviewer", fmt.Sprintf("Review output:\n%s", work))
		if err != nil {
			ch <- pluginapi.OrchestrationEvent{Type: "agent.error", Error: err.Error(), Done: true}
			return
		}
		final := strings.TrimSpace(review)
		if final == "" {
			final = work
		}
		env.Emit(events.OrchestrationFinal, map[string]any{"content": final})
		ch <- pluginapi.OrchestrationEvent{Type: "agent.completed", Role: "coordinator", Content: final, Done: true}
	}()
	return ch, nil
}

func (o *Orchestrator) runRole(ctx context.Context, env pluginapi.ExecutionEnvironment, ch chan<- pluginapi.OrchestrationEvent, role, prompt string) (string, error) {
	ch <- pluginapi.OrchestrationEvent{Type: "agent.started", Role: role}
	nodeID, _ := env.NodeForRole(role)
	env.Emit(events.OrchestrationRole, map[string]any{"role": role, "node_id": nodeID})

	stream, err := env.Generate(ctx, role, []pluginapi.ChatMessage{{Role: "user", Content: prompt}})
	if err != nil {
		return "", err
	}
	var b strings.Builder
	var metrics *pluginapi.GenerationMetrics
	for chunk := range stream {
		if chunk.Error != "" {
			return b.String(), fmt.Errorf("%s", chunk.Error)
		}
		if chunk.Metrics != nil {
			metrics = chunk.Metrics
		}
		if chunk.Content != "" {
			b.WriteString(chunk.Content)
			ch <- pluginapi.OrchestrationEvent{Type: "agent.message", Role: role, Content: chunk.Content, NodeID: nodeID}
		}
		if chunk.Done {
			break
		}
	}
	ch <- pluginapi.OrchestrationEvent{
		Type:    "agent.completed",
		Role:    role,
		NodeID:  nodeID,
		Content: b.String(),
		Metrics: metrics,
	}
	return b.String(), nil
}
