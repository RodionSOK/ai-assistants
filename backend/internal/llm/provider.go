package llm

import (
	"context"

	"github.com/RodionSOK/ai-assistants/internal/config"
)

type Request struct {
	Model        string
	SystemPrompt string
	UserPrompt   string
}

type Response struct {
	Output string
}

type Provider interface {
	Complete(ctx context.Context, req Request) (*Response, error)
}

func NewProvider(cfg *config.Config) Provider {
	switch cfg.LLMProvider {
	case "openai":
		return NewOpenAIProvider(cfg)
	default:
		return NewMockProvider()
	}
}
