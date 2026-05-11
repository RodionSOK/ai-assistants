package usecase

import (
	"context"
	"fmt"

	"github.com/RodionSOK/ai-assistants/internal/domain"
	"github.com/RodionSOK/ai-assistants/internal/llm"
	"github.com/RodionSOK/ai-assistants/internal/repository"
	"github.com/google/uuid"
)

type RunRepository interface {
	Create(ctx context.Context, run domain.Run) (*domain.Run, error)
	Update(ctx context.Context, run domain.Run) (*domain.Run, error)
	List(ctx context.Context, f repository.ListRunsFilter) ([]domain.Run, int, error)
}

type LLMProvider interface {
	Complete(ctx context.Context, req llm.Request) (*llm.Response, error)
}

type RunUsecase struct {
	runRepo       RunRepository
	assistantRepo AssistantRepository
	llmProvider   LLMProvider
}

func NewRunUsecase(runRepo RunRepository, assistantRepo AssistantRepository, llmProvider LLMProvider) *RunUsecase {
	return &RunUsecase{
		runRepo:       runRepo,
		assistantRepo: assistantRepo,
		llmProvider:   llmProvider,
	}
}

type ListRunsInput struct {
	UserID      string
	AssistantID string
	Status      string
	Page        int
	PageSize    int
}

type ListRunsOutput struct {
	Runs     []domain.Run
	Total    int
	Page     int
	PageSize int
}

func (uc *RunUsecase) Run(ctx context.Context, assistantID, userID, userPrompt string) (*domain.Run, error) {
	assistant, err := uc.assistantRepo.GetByID(ctx, assistantID)
	if err != nil {
		return nil, fmt.Errorf("assistant not found")
	}

	if !assistant.IsActive {
		return nil, fmt.Errorf("assistant is inactive")
	}

	run, err := uc.runRepo.Create(ctx, domain.Run{
		ID:          uuid.New().String(),
		AssistantID: assistantID,
		UserID:      userID,
		Model:       assistant.Model,
		UserPrompt:  userPrompt,
		Status:      domain.RunStatusPending,
	})
	if err != nil {
		return nil, fmt.Errorf("RunUsecase.Run create: %w", err)
	}

	llmResp, llmErr := uc.llmProvider.Complete(ctx, llm.Request{
		Model:        assistant.Model,
		SystemPrompt: assistant.SystemPrompt,
		UserPrompt:   userPrompt,
	})

	if llmErr != nil {
		run.Status = domain.RunStatusFailed
		run.Error = llmErr.Error()
	} else {
		run.Status = domain.RunStatusSuccess
		run.Output = llmResp.Output
	}

	updated, err := uc.runRepo.Update(ctx, *run)
	if err != nil {
		return nil, fmt.Errorf("RunUsecase.Run update: %w", err)
	}

	updated.AssistantName = assistant.Name
	updated.CategoryID = assistant.CategoryID
	updated.CategoryName = assistant.CategoryName

	if llmErr != nil {
		return updated, fmt.Errorf("LLM provider error: %w", llmErr)
	}

	return updated, nil
}

func (uc *RunUsecase) My(ctx context.Context, in ListRunsInput) (*ListRunsOutput, error) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.PageSize < 1 || in.PageSize > 100 {
		in.PageSize = 20
	}

	runs, total, err := uc.runRepo.List(ctx, repository.ListRunsFilter{
		UserID:   in.UserID,
		Status:   in.Status,
		Page:     in.Page,
		PageSize: in.PageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("RunUsecase.My: %w", err)
	}

	return &ListRunsOutput{
		Runs:     runs,
		Total:    total,
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}

func (uc *RunUsecase) All(ctx context.Context, in ListRunsInput) (*ListRunsOutput, error) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.PageSize < 1 || in.PageSize > 100 {
		in.PageSize = 20
	}

	runs, total, err := uc.runRepo.List(ctx, repository.ListRunsFilter{
		AssistantID: in.AssistantID,
		Status:      in.Status,
		Page:        in.Page,
		PageSize:    in.PageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("RunUsecase.All: %w", err)
	}

	return &ListRunsOutput{
		Runs:     runs,
		Total:    total,
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}
