package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/RodionSOK/ai-assistants/internal/domain"
	"github.com/RodionSOK/ai-assistants/internal/usecase"
)

type mockCategoryRepo struct {
	createFn  func(ctx context.Context, name, description string) (*domain.Category, error)
	listFn    func(ctx context.Context) ([]domain.Category, error)
	getByIDFn func(ctx context.Context, id string) (*domain.Category, error)
}

func (m *mockCategoryRepo) Create(ctx context.Context, name, description string) (*domain.Category, error) {
	return m.createFn(ctx, name, description)
}
func (m *mockCategoryRepo) List(ctx context.Context) ([]domain.Category, error) {
	return m.listFn(ctx)
}
func (m *mockCategoryRepo) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	return m.getByIDFn(ctx, id)
}

func TestCategoryUsecase_Create_Success(t *testing.T) {
	repo := &mockCategoryRepo{
		createFn: func(_ context.Context, name, description string) (*domain.Category, error) {
			return &domain.Category{ID: "cat-1", Name: name, Description: description}, nil
		},
	}
	uc := usecase.NewCategoryUsecase(repo)

	cat, err := uc.Create(context.Background(), "AI Tools", "desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cat.Name != "AI Tools" {
		t.Errorf("expected name 'AI Tools', got '%s'", cat.Name)
	}
}

func TestCategoryUsecase_Create_EmptyName(t *testing.T) {
	repo := &mockCategoryRepo{}
	uc := usecase.NewCategoryUsecase(repo)

	_, err := uc.Create(context.Background(), "", "desc")
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestCategoryUsecase_Create_RepoError(t *testing.T) {
	repo := &mockCategoryRepo{
		createFn: func(_ context.Context, _, _ string) (*domain.Category, error) {
			return nil, errors.New("db error")
		},
	}
	uc := usecase.NewCategoryUsecase(repo)

	_, err := uc.Create(context.Background(), "Valid Name", "")
	if err == nil {
		t.Fatal("expected error from repo, got nil")
	}
}
