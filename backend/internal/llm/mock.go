package llm

import (
	"context"
	"fmt"
)

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) Complete(ctx context.Context, req Request) (*Response, error) {
	output := fmt.Sprintf(
		"[ Mock LLM | model: %s ]\n\nСистемный промпт: %s\n\nОтвет на запрос \"%s\":\nЭто детерминированный ответ мок-провайдера. В реальном приложении здесь был бы ответ языковой модели.",
		req.Model,
		req.SystemPrompt,
		req.UserPrompt,
	)

	return &Response{Output: output}, nil
}
