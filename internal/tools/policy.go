package tools

import (
	"fmt"
	"strings"
)

// Policy constants.
const (
	PolicyDeny            = "deny"
	PolicyAsk             = "ask"
	PolicyAllowForSession = "allow-for-session"
	PolicyAllow           = "allow"
)

// PolicyEngine resolves tool permissions.
type PolicyEngine struct {
	sessionAllowed map[string]struct{}
}

func NewPolicyEngine() *PolicyEngine {
	return &PolicyEngine{sessionAllowed: make(map[string]struct{})}
}

// Decide returns whether a tool may run and whether user approval is needed.
func (p *PolicyEngine) Decide(toolID, policy string) (allowed bool, needsPrompt bool, err error) {
	policy = strings.ToLower(strings.TrimSpace(policy))
	if policy == "" {
		policy = PolicyAsk
	}
	switch policy {
	case PolicyDeny:
		return false, false, fmt.Errorf("tool %q denied by policy", toolID)
	case PolicyAllow:
		return true, false, nil
	case PolicyAllowForSession:
		if _, ok := p.sessionAllowed[toolID]; ok {
			return true, false, nil
		}
		return false, true, nil
	case PolicyAsk:
		return false, true, nil
	default:
		return false, false, fmt.Errorf("unknown policy %q", policy)
	}
}

// AllowSession grants session-scoped permission.
func (p *PolicyEngine) AllowSession(toolID string) {
	p.sessionAllowed[toolID] = struct{}{}
}

// PolicyFor returns the policy string for a tool from profile policies.
func PolicyFor(toolID string, policies []ToolPolicy) string {
	for _, tp := range policies {
		if tp.ToolID == toolID {
			return tp.Policy
		}
	}
	return PolicyAsk
}

// ToolPolicy mirrors contracts.ToolPolicy locally.
type ToolPolicy struct {
	ToolID string `json:"tool_id"`
	Policy string `json:"policy"`
}
