package profiles

import "github.com/yeixio/yggdrasil-core/pkg/contracts"

// Profile is an alias for the public contract.
type Profile = contracts.AIProfile

// Validate checks basic profile invariants.
func Validate(p Profile) error {
	if p.Name == "" {
		return ErrInvalidProfile("name is required")
	}
	if p.OrchestratorID == "" {
		return ErrInvalidProfile("orchestrator_id is required")
	}
	for _, r := range p.Roles {
		if r.Role == "" {
			return ErrInvalidProfile("role name is required")
		}
		if r.Required && r.ModelID == "" {
			return ErrInvalidProfile("required role " + r.Role + " missing model_id")
		}
	}
	return nil
}

// ErrInvalidProfile indicates profile validation failure.
type ErrInvalidProfile string

func (e ErrInvalidProfile) Error() string { return string(e) }
