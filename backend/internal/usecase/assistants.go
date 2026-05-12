package usecase

import (
    "context"
    "fmt"

    "github.com/RodionSOK/ai-assistants/internal/domain"
    "github.com/RodionSOK/ai-assistants/internal/repository"
)

type AssistantRepository interface {
    List(ctx context.Context, f repository.ListAssistantsFilter) ([]domain.Assistant, int, error)
    GetByID(ctx context.Context, id string) (*domain.Assistant, error)
    Create(ctx context.Context, a domain.Assistant) (*domain.Assistant, error)
    Update(ctx context.Context, a domain.Assistant) (*domain.Assistant, error)
}

type AssistantUsecase struct {
    assistantRepo AssistantRepository
    categoryRepo  CategoryRepository
}

func NewAssistantUsecase(assistantRepo AssistantRepository, categoryRepo CategoryRepository) *AssistantUsecase {
    return &AssistantUsecase{
        assistantRepo: assistantRepo,
        categoryRepo:  categoryRepo,
    }
}

type ListAssistantsInput struct {
    CategoryID      string
    Query           string
    IncludeInactive bool
    IsAdmin         bool
    Page            int
    PageSize        int
}

type ListAssistantsOutput struct {
    Assistants []domain.Assistant
    Total      int
    Page       int
    PageSize   int
}

func (uc *AssistantUsecase) List(ctx context.Context, in ListAssistantsInput) (*ListAssistantsOutput, error) {
    includeInactive := in.IncludeInactive && in.IsAdmin

    if in.Page < 1 {
        in.Page = 1
    }
    if in.PageSize < 1 || in.PageSize > 100 {
        in.PageSize = 10
    }

    assistants, total, err := uc.assistantRepo.List(ctx, repository.ListAssistantsFilter{
        CategoryID:      in.CategoryID,
        Query:           in.Query,
        IncludeInactive: includeInactive,
        Page:            in.Page,
        PageSize:        in.PageSize,
    })
    if err != nil {
        return nil, fmt.Errorf("AssistantUsecase.List: %w", err)
    }

    return &ListAssistantsOutput{
        Assistants: assistants,
        Total:      total,
        Page:       in.Page,
        PageSize:   in.PageSize,
    }, nil
}

func (uc *AssistantUsecase) GetByID(ctx context.Context, id string) (*domain.Assistant, error) {
    assistant, err := uc.assistantRepo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("AssistantUsecase.GetByID: %w", err)
    }

    return assistant, nil
}

type CreateAssistantInput struct {
    CategoryID        string
    Name              string
    Description       string
    Model             string
    SystemPrompt      string
    ExampleUserPrompt string
    IsActive          bool
}

func (uc *AssistantUsecase) Create(ctx context.Context, in CreateAssistantInput) (*domain.Assistant, error) {
    if in.SystemPrompt == "" {
        return nil, fmt.Errorf("system prompt is required")
    }

    if _, err := uc.categoryRepo.GetByID(ctx, in.CategoryID); err != nil {
        return nil, fmt.Errorf("category not found")
    }

    assistant, err := uc.assistantRepo.Create(ctx, domain.Assistant{
        CategoryID:        in.CategoryID,
        Name:              in.Name,
        Description:       in.Description,
        Model:             in.Model,
        SystemPrompt:      in.SystemPrompt,
        ExampleUserPrompt: in.ExampleUserPrompt,
        IsActive:          in.IsActive,
    })
    if err != nil {
        return nil, fmt.Errorf("AssistantUsecase.Create: %w", err)
    }

    return assistant, nil
}

type UpdateAssistantInput struct {
    ID                string
    CategoryID        string
    Name              string
    Description       string
    Model             string
    SystemPrompt      string
    ExampleUserPrompt string
    IsActive          bool
}

func (uc *AssistantUsecase) Update(ctx context.Context, in UpdateAssistantInput) (*domain.Assistant, error) {
    if in.SystemPrompt == "" {
        return nil, fmt.Errorf("system prompt is required")
    }

    if _, err := uc.assistantRepo.GetByID(ctx, in.ID); err != nil {
        return nil, fmt.Errorf("assistant not found")
    }

    if _, err := uc.categoryRepo.GetByID(ctx, in.CategoryID); err != nil {
        return nil, fmt.Errorf("category not found")
    }

    assistant, err := uc.assistantRepo.Update(ctx, domain.Assistant{
        ID:                in.ID,
        CategoryID:        in.CategoryID,
        Name:              in.Name,
        Description:       in.Description,
        Model:             in.Model,
        SystemPrompt:      in.SystemPrompt,
        ExampleUserPrompt: in.ExampleUserPrompt,
        IsActive:          in.IsActive,
    })
    if err != nil {
        return nil, fmt.Errorf("AssistantUsecase.Update: %w", err)
    }

    return assistant, nil
}