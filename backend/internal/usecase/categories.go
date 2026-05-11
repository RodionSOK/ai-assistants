package usecase

import (
	"context"
	"fmt"

	"github.com/RodionSOK/ai-assistants/internal/domain"
)

type CategoryRepository interface {
	List(ctx context.Context) ([]domain.Category, error)
	Create(ctx context.Context, name, description string) (*domain.Category, error)
	GetByID(ctx context.Context, id string) (*domain.Category, error)
}

type CategoryUsecase struct {
	categoryRepo CategoryRepository
}

func NewCategoryUsecase(categoryRepo CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{categoryRepo: categoryRepo}
}

func (uc *CategoryUsecase) List(ctx context.Context) ([]domain.Category, error) {
	categories, err := uc.categoryRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("CategoryUsecase.List: %w", err)
	}

	return categories, nil
}

func (uc *CategoryUsecase) Create(ctx context.Context, name, description string) (*domain.Category, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	category, err := uc.categoryRepo.Create(ctx, name, description)
	if err != nil {
		return nil, fmt.Errorf("CategoryUsecase.Create: %w", err)
	}

	return category, nil
}
