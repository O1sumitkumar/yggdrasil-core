package external

import (
	"context"
	"fmt"

	"github.com/yeixio/yggdrasil-core/internal/runtimes/llamacpp"
	"github.com/yeixio/yggdrasil-core/pkg/pluginapi"
)

const runtimeID = "external-openai"

// Config holds external OpenAI-compatible endpoint settings.
type Config struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key,omitempty"`
}

// Runtime adapts a remote OpenAI-compatible server.
type Runtime struct {
	cfg Config
}

// New creates an external OpenAI runtime adapter.
func New(cfg Config) *Runtime {
	return &Runtime{cfg: cfg}
}

func (r *Runtime) ID() string          { return runtimeID }
func (r *Runtime) DisplayName() string { return "External OpenAI-compatible" }

func (r *Runtime) Detect(ctx context.Context) (pluginapi.RuntimeDetection, error) {
	if r.cfg.BaseURL == "" {
		return pluginapi.RuntimeDetection{
			Installed: false,
			Message:   "Configure base_url in settings to use an external OpenAI-compatible endpoint",
		}, nil
	}
	return pluginapi.RuntimeDetection{
		Installed: true,
		Path:      r.cfg.BaseURL,
		Version:   "remote",
	}, nil
}

func (r *Runtime) Install(ctx context.Context, opts pluginapi.InstallOptions) error {
	return fmt.Errorf("external runtime cannot be installed; configure base_url instead")
}

func (r *Runtime) Update(ctx context.Context) error {
	return nil
}

func (r *Runtime) Capabilities(ctx context.Context) (pluginapi.RuntimeCapabilities, error) {
	return pluginapi.RuntimeCapabilities{
		Backends:          []string{"remote"},
		SupportsStreaming: true,
		SupportsTools:     true,
		SupportsGPU:       false,
	}, nil
}

func (r *Runtime) StartModel(ctx context.Context, cfg pluginapi.ModelStartConfig) (pluginapi.RunningModel, error) {
	det, err := r.Detect(ctx)
	if err != nil {
		return pluginapi.RunningModel{}, err
	}
	if !det.Installed {
		return pluginapi.RunningModel{}, fmt.Errorf("external runtime not configured")
	}
	return pluginapi.RunningModel{
		ID:        cfg.ModelID,
		ModelID:   cfg.ModelID,
		Endpoint:  r.cfg.BaseURL,
		Status:    "remote",
		RuntimeID: runtimeID,
	}, nil
}

func (r *Runtime) StopModel(ctx context.Context, id string) error { return nil }

func (r *Runtime) ListRunning(ctx context.Context) ([]pluginapi.RunningModel, error) {
	return nil, nil
}

func (r *Runtime) Health(ctx context.Context) error {
	det, err := r.Detect(ctx)
	if err != nil {
		return err
	}
	if !det.Installed {
		return fmt.Errorf("external runtime not configured")
	}
	return nil
}

// Client wraps llamacpp client for OpenAI-compatible remote endpoints.
type Client struct {
	inner *llamacpp.Client
}

func NewClient() *Client {
	return &Client{inner: llamacpp.NewClient()}
}

func (c *Client) Chat(ctx context.Context, req pluginapi.ChatRequest) (<-chan pluginapi.ChatChunk, error) {
	return c.inner.Chat(ctx, req)
}
