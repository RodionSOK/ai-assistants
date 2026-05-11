package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/RodionSOK/ai-assistants/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssistantRepo struct {
	pool *pgxpool.Pool
}

func NewAssistantRepo(pool *pgxpool.Pool) *AssistantRepo {
	return &AssistantRepo{pool: pool}
}

type ListAssistantsFilter struct {
	CategoryID      string
	Query           string
	IncludeInactive bool
	Page            int
	PageSize        int
}

func (r *AssistantRepo) List(ctx context.Context, f ListAssistantsFilter) ([]domain.Assistant, int, error) {
	conditions := []string{}
	args := []any{}
	i := 1

	if !f.IncludeInactive {
		conditions = append(conditions, fmt.Sprintf("a.is_active = $%d", i))
		args = append(args, true)
		i++
	}

	if f.CategoryID != "" {
		conditions = append(conditions, fmt.Sprintf("a.category_id = $%d", i))
		args = append(args, f.CategoryID)
		i++
	}

	if f.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(a.name ILIKE $%d OR a.description ILIKE $%d)", i, i))
		args = append(args, "%"+f.Query+"%")
		i++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM assistants a %s`, where)
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("AssistantRepo.List count: %w", err)
	}

	offset := (f.Page - 1) * f.PageSize
	args = append(args, f.PageSize, offset)

	query := fmt.Sprintf(`
        SELECT
            a.id, a.category_id, COALESCE(c.name, ''),
            a.name, a.description, a.model,
            a.system_prompt, COALESCE(a.example_user_prompt, ''),
            a.is_active, a.created_at, a.updated_at
        FROM assistants a
        LEFT JOIN categories c ON c.id = a.category_id
        %s
        ORDER BY a.created_at DESC
        LIMIT $%d OFFSET $%d
    `, where, i, i+1)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("AssistantRepo.List: %w", err)
	}
	defer rows.Close()

	var assistants []domain.Assistant
	for rows.Next() {
		var a domain.Assistant
		if err := rows.Scan(
			&a.ID, &a.CategoryID, &a.CategoryName,
			&a.Name, &a.Description, &a.Model,
			&a.SystemPrompt, &a.ExampleUserPrompt,
			&a.IsActive, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("AssistantRepo.List scan: %w", err)
		}
		assistants = append(assistants, a)
	}

	return assistants, total, nil
}

func (r *AssistantRepo) GetByID(ctx context.Context, id string) (*domain.Assistant, error) {
	var a domain.Assistant
	err := r.pool.QueryRow(ctx, `
        SELECT
            a.id, a.category_id, COALESCE(c.name, ''),
            a.name, a.description, a.model,
            a.system_prompt, COALESCE(a.example_user_prompt, ''),
            a.is_active, a.created_at, a.updated_at
        FROM assistants a
        LEFT JOIN categories c ON c.id = a.category_id
        WHERE a.id = $1
    `, id).Scan(
		&a.ID, &a.CategoryID, &a.CategoryName,
		&a.Name, &a.Description, &a.Model,
		&a.SystemPrompt, &a.ExampleUserPrompt,
		&a.IsActive, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("AssistantRepo.GetByID: %w", err)
	}

	return &a, nil
}

func (r *AssistantRepo) Create(ctx context.Context, a domain.Assistant) (*domain.Assistant, error) {
	var created domain.Assistant
	err := r.pool.QueryRow(ctx, `
        INSERT INTO assistants
            (category_id, name, description, model, system_prompt, example_user_prompt, is_active)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING
            id, category_id, name, description, model,
            system_prompt, COALESCE(example_user_prompt, ''),
            is_active, created_at, updated_at
    `,
		a.CategoryID, a.Name, a.Description, a.Model,
		a.SystemPrompt, a.ExampleUserPrompt, a.IsActive,
	).Scan(
		&created.ID, &created.CategoryID, &created.Name,
		&created.Description, &created.Model, &created.SystemPrompt,
		&created.ExampleUserPrompt, &created.IsActive,
		&created.CreatedAt, &created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("AssistantRepo.Create: %w", err)
	}

	created.CategoryID = a.CategoryID
	return &created, nil
}

func (r *AssistantRepo) Update(ctx context.Context, a domain.Assistant) (*domain.Assistant, error) {
	var updated domain.Assistant
	err := r.pool.QueryRow(ctx, `
        UPDATE assistants
        SET
            category_id         = $1,
            name                = $2,
            description         = $3,
            model               = $4,
            system_prompt       = $5,
            example_user_prompt = $6,
            is_active           = $7,
            updated_at          = NOW()
        WHERE id = $8
        RETURNING
            id, category_id, name, description, model,
            system_prompt, COALESCE(example_user_prompt, ''),
            is_active, created_at, updated_at
    `,
		a.CategoryID, a.Name, a.Description, a.Model,
		a.SystemPrompt, a.ExampleUserPrompt, a.IsActive, a.ID,
	).Scan(
		&updated.ID, &updated.CategoryID, &updated.Name,
		&updated.Description, &updated.Model, &updated.SystemPrompt,
		&updated.ExampleUserPrompt, &updated.IsActive,
		&updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("AssistantRepo.Update: %w", err)
	}

	return &updated, nil
}
