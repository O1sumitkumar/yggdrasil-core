package scheduler

import (
	"github.com/yeixio/yggdrasil-core/pkg/contracts"
)

// ScoreInput drives deterministic node scoring.
type ScoreInput struct {
	Role          string
	ModelID       string
	Profile       contracts.AIProfile
	Nodes         []NodeCandidate
	RunningModels map[string]map[string]struct{} // nodeID -> modelIDs
	ActiveTasks   map[string]int                 // nodeID -> count
	// AvoidNodeIDs softly prefers placing this role on a different computer
	// than earlier roles in the same Team turn (automatic spread).
	AvoidNodeIDs []string
}

// NodeCandidate is a schedulable node snapshot.
type NodeCandidate struct {
	Node            contracts.Node
	InstalledModels map[string]struct{}
	FreeMemory      uint64
}

// ScoredNode includes placement score and reason fragments.
type ScoredNode struct {
	NodeID  string
	Score   int
	Reasons []string
}

// Score ranks nodes for a role placement.
func Score(input ScoreInput) []ScoredNode {
	var out []ScoredNode
	for _, n := range input.Nodes {
		s := ScoredNode{NodeID: n.Node.ID}
		if !n.Node.Paired && !n.Node.IsLocal {
			continue
		}
		if n.Node.Status != contracts.NodeStatusOnline && n.Node.Status != "" {
			s.Reasons = append(s.Reasons, "offline")
			continue
		}

		// Explicit pin.
		for _, r := range input.Profile.Roles {
			if r.Role == input.Role && r.NodeID != "" {
				if r.NodeID == n.Node.ID {
					s.Score += 1000
					s.Reasons = append(s.Reasons, "explicitly pinned")
				} else {
					s.Score -= 1000
					s.Reasons = append(s.Reasons, "pinned elsewhere")
				}
			}
		}

		if _, ok := n.InstalledModels[input.ModelID]; ok {
			s.Score += 200
			s.Reasons = append(s.Reasons, "model installed")
		} else if input.ModelID != "" {
			s.Score -= 500
			s.Reasons = append(s.Reasons, "model not installed")
		}

		if running := input.RunningModels[n.Node.ID]; running != nil {
			if _, ok := running[input.ModelID]; ok {
				s.Score += 150
				s.Reasons = append(s.Reasons, "model already running")
			}
		}

		if n.FreeMemory > 0 {
			s.Score += 50
			s.Reasons = append(s.Reasons, "has free memory")
		}

		if n.Node.IsLocal {
			if input.Profile.NodePolicy.Mode == "prefer_local" {
				s.Score += 100
				s.Reasons = append(s.Reasons, "prefer_local policy")
			} else {
				// Tiny tie-break only — easy mode should happily use remotes.
				s.Score += 5
				s.Reasons = append(s.Reasons, "local node tie-breaker")
			}
		}

		for _, avoid := range input.AvoidNodeIDs {
			if avoid != "" && avoid == n.Node.ID {
				// Only soft-spread under automatic policy; prefer_local keeps work here.
				if input.Profile.NodePolicy.Mode != "prefer_local" {
					s.Score -= 180
					s.Reasons = append(s.Reasons, "spread across computers")
				}
				break
			}
		}

		taskCount := input.ActiveTasks[n.Node.ID]
		s.Score -= taskCount * 10
		if taskCount > 0 {
			s.Reasons = append(s.Reasons, "lower task load preferred")
		}

		out = append(out, s)
	}
	return out
}
