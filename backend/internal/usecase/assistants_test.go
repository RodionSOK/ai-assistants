package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/RodionSOK/ai-assistants/internal/domain"
	"github.com/RodionSOK/ai-assistants/internal/repository"
	"github.com/RodionSOK/ai-assistants/internal/usecase"
)

type mockAssistantRepo struct {
	listFn    func(ctx context.Context, f repository.ListAssistantsFilter) ([]domain.Assistant, int, error)
	getByIDFn func(ctx context.Context, id string) (*domain.Assistant, error)
	createFn  func(ctx context.Context, a domain.Assistant) (*domain.Assistant, error)
	updateFn  func(ctx context.Context, a domain.Assistant) (*domain.Assistant, error)
}

func (m *mockAssistantRepo) List(ctx context.Context, f repository.ListAssistantsFilter) ([]domain.Assistant, int, error) {
	return m.listFn(ctx, f)
}
func (m *mockAssistantRepo) GetByID(ctx context.Context, id string) (*domain.Assistant, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockAssistantRepo) Create(ctx context.Context, a domain.Assistant) (*domain.Assistant, error) {
	return m.createFn(ctx, a)
}
func (m *mockAssistantRepo) Update(ctx context.Context, a domain.Assistant) (*domain.Assistant, error) {
	return m.updateFn(ctx, a)
}

func newAssistantUC(aRepo *mockAssistantRepo, cRepo *mockCategoryRepo) *usecase.AssistantUsecase {
	return usecase.NewAssistantUsecase(aRepo, cRepo)
}

func TestAssistantUsecase_Create_Success(t *testing.T) {
	catRepo := &mockCategoryRepo{
		getByIDFn: func(_ context.Context, id string) (*domain.Category, error) {
			return &domain.Category{ID: id}, nil
		},
	}
	aRepo := &mockAssistantRepo{
		createFn: func(_ context.Context, a domain.Assistant) (*domain.Assistant, error) {
			a.ID = "ast-1"
			return &a, nil
		},
	}
	uc := newAssistantUC(aRepo, catRepo)

	result, err := uc.Create(context.Background(), usecase.CreateAssistantInput{
		CategoryID:   "cat-1",
		Name:         "My Assistant",
		SystemPrompt: "You are helpful.",
		Model:        "gpt-4",
		IsActive:     true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "ast-1" {
		t.Errorf("expected ID 'ast-1', got '%s'", result.ID)
	}
}

func TestAssistantUsecase_Create_EmptySystemPrompt(t *testing.T) {
	uc := newAssistantUC(&mockAssistantRepo{}, &mockCategoryRepo{})

	_, err := uc.Create(context.Background(), usecase.CreateAssistantInput{
		CategoryID:   "cat-1",
		Name:         "My Assistant",
		SystemPrompt: "",
	})
	if err == nil {
		t.Fatal("expected error for empty system prompt, got nil")
	}
}

func TestAssistantUsecase_Create_CategoryNotFound(t *testing.T) {
	catRepo := &mockCategoryRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Category, error) {
			return nil, errors.New("not found")
		},
	}
	uc := newAssistantUC(&mockAssistantRepo{}, catRepo)

	_, err := uc.Create(context.Background(), usecase.CreateAssistantInput{
		CategoryID:   "nonexistent",
		SystemPrompt: "prompt",
	})
	if err == nil {
		t.Fatal("expected error for missing category, got nil")
	}
}

func TestAssistantUsecase_List_NonAdminCannotSeeInactive(t *testing.T) {
	aRepo := &mockAssistantRepo{
		listFn: func(_ context.Context, f repository.ListAssistantsFilter) ([]domain.Assistant, int, error) {
			// Проверяем, что несмотря на IncludeInactive=true, без IsAdmin флаг не пробрасывается
			if f.IncludeInactive {
				return nil, 0, errors.New("non-admin should not see inactive")
			}
			return []domain.Assistant{}, 0, nil
		},
	}
	uc := newAssistantUC(aRepo, &mockCategoryRepo{})

	_, err := uc.List(context.Background(), usecase.ListAssistantsInput{
		IncludeInactive: true,
		IsAdmin:         false, // не админ
		Page:            1,
		PageSize:        10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAssistantUsecase_List_FilterByCategoryID(t *testing.T) {
	wantCategoryID := "cat-42"

	aRepo := &mockAssistantRepo{
		listFn: func(_ context.Context, f repository.ListAssistantsFilter) ([]domain.Assistant, int, error) {
			if f.CategoryID != wantCategoryID {
				return nil, 0, errors.New("wrong category filter passed")
			}
			return []domain.Assistant{
				{ID: "a1", CategoryID: wantCategoryID},
				{ID: "a2", CategoryID: wantCategoryID},
			}, 2, nil
		},
	}
	uc := newAssistantUC(aRepo, &mockCategoryRepo{})

	out, err := uc.List(context.Background(), usecase.ListAssistantsInput{
		CategoryID: wantCategoryID,
		Page:       1,
		PageSize:   10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Assistants) != 2 {
		t.Errorf("expected 2 assistants, got %d", len(out.Assistants))
	}
	for _, a := range out.Assistants {
		if a.CategoryID != wantCategoryID {
			t.Errorf("assistant %s has wrong categoryID: %s", a.ID, a.CategoryID)
		}
	}
}
