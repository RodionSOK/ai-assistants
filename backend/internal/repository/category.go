package repository

import (
	"context"
	"fmt"

	"github.com/RodionSOK/ai-assistants/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepo struct {
	pool *pgxpool.Pool
}

func NewCategoryRepo(pool *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{pool: pool}
}

func (r *CategoryRepo) List(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, COALESCE(description, ''), created_at
		FROM categories
		ORDER BY name ASC	
	`)
	if err != nil {
		return nil, fmt.Errorf("CategoryRepo.List: %w", err)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("CategoryRepo.List scan: %w", err)
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (r *CategoryRepo) Create(ctx context.Context, name, description string) (*domain.Category, error) {
	var c domain.Category
	err := r.pool.QueryRow(ctx, `
        INSERT INTO categories (name, description)
        VALUES ($1, $2)
        RETURNING id, name, COALESCE(description, ''), created_at
    `, name, description).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("CategoryRepo.Create: %w", err)
	}

	return &c, nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	var c domain.Category
	err := r.pool.QueryRow(ctx, `
        SELECT id, name, COALESCE(description, ''), created_at
        FROM categories
        WHERE id = $1
    `, id).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("CategoryRepo.GetByID: %w", err)
	}

	return &c, nil
}
