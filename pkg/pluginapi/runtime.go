package pluginapi

import "context"

// RuntimeDetection reports what a runtime found on the host.
type RuntimeDetection struct {
	Installed bool   `json:"installed"`
	Version   string `json:"version,omitempty"`
	Path      string `json:"path,omitempty"`
	Message   string `json:"message,omitempty"`
}

// InstallOptions configures runtime installation.
type InstallOptions struct {
	Force bool `json:"force,omitempty"`
}

// RuntimeCapabilities describes what the runtime can do.
type RuntimeCapabilities struct {
	Backends          []string `json:"backends"`
	SupportsStreaming bool     `json:"supports_streaming"`
	SupportsTools     bool     `json:"supports_tools"`
	SupportsGPU       bool     `json:"supports_gpu"`
}

// ModelStartConfig starts a model under a runtime.
type ModelStartConfig struct {
	ModelID   string `json:"model_id"`
	ModelPath string `json:"model_path"`
	Context   int    `json:"context,omitempty"`
	GPULayers int    `json:"gpu_layers,omitempty"`
	Port      int    `json:"port,omitempty"`
}

// RunningModel is a loaded model instance.
type RunningModel struct {
	ID        string `json:"id"`
	ModelID   string `json:"model_id"`
	Endpoint  string `json:"endpoint"`
	Status    string `json:"status"`
	RuntimeID string `json:"runtime_id"`
}

// Runtime is the replaceable inference runtime adapter.
type Runtime interface {
	ID() string
	DisplayName() string

	Detect(ctx context.Context) (RuntimeDetection, error)
	Install(ctx context.Context, opts InstallOptions) error
	Update(ctx context.Context) error

	Capabilities(ctx context.Context) (RuntimeCapabilities, error)

	StartModel(ctx context.Context, cfg ModelStartConfig) (RunningModel, error)
	StopModel(ctx context.Context, id string) error
	ListRunning(ctx context.Context) ([]RunningModel, error)

	Health(ctx context.Context) error
}

// ChatRequest is a generation request to a running model.
type ChatRequest struct {
	ModelEndpoint string           `json:"model_endpoint"`
	Messages      []ChatMessage    `json:"messages"`
	Temperature   float64          `json:"temperature,omitempty"`
	MaxTokens     int              `json:"max_tokens,omitempty"`
	Stream        bool             `json:"stream"`
	Tools         []map[string]any `json:"tools,omitempty"`
}

// ChatMessage is a single chat turn.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatChunk is a streaming generation fragment.
type ChatChunk struct {
	Content string             `json:"content,omitempty"`
	Done    bool               `json:"done"`
	Error   string             `json:"error,omitempty"`
	Metrics *GenerationMetrics `json:"metrics,omitempty"`
}

// GenerationMetrics mirrors contracts.GenerationMetrics for runtime adapters.
type GenerationMetrics struct {
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens"`
	TTFTMs           float64 `json:"ttft_ms"`
	PromptMs         float64 `json:"prompt_ms"`
	EvalMs           float64 `json:"eval_ms"`
	TotalMs          float64 `json:"total_ms"`
	PromptTokPerSec  float64 `json:"prompt_tok_per_sec"`
	EvalTokPerSec    float64 `json:"eval_tok_per_sec"`
}

// Generator can stream chat completions against a running model.
type Generator interface {
	Chat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
}
