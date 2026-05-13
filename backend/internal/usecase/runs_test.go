package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/RodionSOK/ai-assistants/internal/domain"
	"github.com/RodionSOK/ai-assistants/internal/llm"
	"github.com/RodionSOK/ai-assistants/internal/repository"
	"github.com/RodionSOK/ai-assistants/internal/usecase"
)

type mockRunRepo struct {
	createFn func(ctx context.Context, run domain.Run) (*domain.Run, error)
	updateFn func(ctx context.Context, run domain.Run) (*domain.Run, error)
	listFn   func(ctx context.Context, f repository.ListRunsFilter) ([]domain.Run, int, error)
}

func (m *mockRunRepo) Create(ctx context.Context, run domain.Run) (*domain.Run, error) {
	return m.createFn(ctx, run)
}
func (m *mockRunRepo) Update(ctx context.Context, run domain.Run) (*domain.Run, error) {
	return m.updateFn(ctx, run)
}
func (m *mockRunRepo) List(ctx context.Context, f repository.ListRunsFilter) ([]domain.Run, int, error) {
	return m.listFn(ctx, f)
}

type mockLLMProvider struct {
	completeFn func(ctx context.Context, req llm.Request) (*llm.Response, error)
}

func (m *mockLLMProvider) Complete(ctx context.Context, req llm.Request) (*llm.Response, error) {
	return m.completeFn(ctx, req)
}

func activeAssistantRepo(isActive bool) *mockAssistantRepo {
	return &mockAssistantRepo{
		getByIDFn: func(_ context.Context, id string) (*domain.Assistant, error) {
			return &domain.Assistant{
				ID:           id,
				Name:         "Test",
				Model:        "gpt-4",
				SystemPrompt: "You are helpful.",
				IsActive:     isActive,
			}, nil
		},
	}
}

func successRunRepo() *mockRunRepo {
	return &mockRunRepo{
		createFn: func(_ context.Context, run domain.Run) (*domain.Run, error) {
			run.ID = "run-1"
			return &run, nil
		},
		updateFn: func(_ context.Context, run domain.Run) (*domain.Run, error) {
			return &run, nil
		},
	}
}

func TestRunUsecase_Run_Success(t *testing.T) {
	uc := usecase.NewRunUsecase(
		successRunRepo(),
		activeAssistantRepo(true),
		&mockLLMProvider{
			completeFn: func(_ context.Context, _ llm.Request) (*llm.Response, error) {
				return &llm.Response{Output: "Hello!"}, nil
			},
		},
	)

	run, err := uc.Run(context.Background(), "ast-1", "user-1", "Hi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if run.Status != domain.RunStatusSuccess {
		t.Errorf("expected status 'success', got '%s'", run.Status)
	}
	if run.Output != "Hello!" {
		t.Errorf("expected output 'Hello!', got '%s'", run.Output)
	}
}

func TestRunUsecase_Run_AssistantNotFound(t *testing.T) {
	aRepo := &mockAssistantRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Assistant, error) {
			return nil, errors.New("not found")
		},
	}
	uc := usecase.NewRunUsecase(successRunRepo(), aRepo, &mockLLMProvider{})

	_, err := uc.Run(context.Background(), "missing", "user-1", "Hi")
	if err == nil {
		t.Fatal("expected error for missing assistant, got nil")
	}
}

func TestRunUsecase_Run_InactiveAssistant(t *testing.T) {
	uc := usecase.NewRunUsecase(
		successRunRepo(),
		activeAssistantRepo(false), // неактивен
		&mockLLMProvider{},
	)

	_, err := uc.Run(context.Background(), "ast-1", "user-1", "Hi")
	if err == nil {
		t.Fatal("expected error for inactive assistant, got nil")
	}
	if err.Error() != "assistant is inactive" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestRunUsecase_Run_LLMError(t *testing.T) {
	var capturedRun domain.Run
	runRepo := &mockRunRepo{
		createFn: func(_ context.Context, run domain.Run) (*domain.Run, error) {
			run.ID = "run-1"
			return &run, nil
		},
		updateFn: func(_ context.Context, run domain.Run) (*domain.Run, error) {
			capturedRun = run
			return &run, nil
		},
	}
	uc := usecase.NewRunUsecase(
		runRepo,
		activeAssistantRepo(true),
		&mockLLMProvider{
			completeFn: func(_ context.Context, _ llm.Request) (*llm.Response, error) {
				return nil, errors.New("timeout")
			},
		},
	)

	_, err := uc.Run(context.Background(), "ast-1", "user-1", "Hi")
	if err == nil {
		t.Fatal("expected error from LLM, got nil")
	}
	if capturedRun.Status != domain.RunStatusFailed {
		t.Errorf("expected run status 'failed', got '%s'", capturedRun.Status)
	}
}
