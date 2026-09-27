package tools

import "context"

// Tool executes a mediated capability.
type Tool interface {
	ID() string
	DisplayName() string
	Description() string
	Execute(ctx context.Context, args map[string]any) (map[string]any, error)
}
