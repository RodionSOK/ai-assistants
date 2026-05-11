package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/RodionSOK/ai-assistants/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RunRepo struct {
	pool *pgxpool.Pool
}

func NewRunRepo(pool *pgxpool.Pool) *RunRepo {
	return &RunRepo{pool: pool}
}

type ListRunsFilter struct {
	UserID      string
	AssistantID string
	Status      string
	Page        int
	PageSize    int
}

func (r *RunRepo) Create(ctx context.Context, run domain.Run) (*domain.Run, error) {
	var created domain.Run
	err := r.pool.QueryRow(ctx, `
        INSERT INTO runs (assistant_id, user_id, model, user_prompt, status)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, assistant_id, user_id, model, user_prompt,
                  COALESCE(output, ''), status, COALESCE(error, ''), created_at
    `,
		run.AssistantID, run.UserID, run.Model, run.UserPrompt, run.Status,
	).Scan(
		&created.ID, &created.AssistantID, &created.UserID,
		&created.Model, &created.UserPrompt, &created.Output,
		&created.Status, &created.Error, &created.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("RunRepo.Create: %w", err)
	}

	return &created, nil
}

func (r *RunRepo) Update(ctx context.Context, run domain.Run) (*domain.Run, error) {
	var updated domain.Run
	err := r.pool.QueryRow(ctx, `
        UPDATE runs
        SET
            output = $1,
            status = $2,
            error  = $3
        WHERE id = $4
        RETURNING id, assistant_id, user_id, model, user_prompt,
                  COALESCE(output, ''), status, COALESCE(error, ''), created_at
    `,
		run.Output, run.Status, run.Error, run.ID,
	).Scan(
		&updated.ID, &updated.AssistantID, &updated.UserID,
		&updated.Model, &updated.UserPrompt, &updated.Output,
		&updated.Status, &updated.Error, &updated.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("RunRepo.Update: %w", err)
	}

	return &updated, nil
}

func (r *RunRepo) List(ctx context.Context, f ListRunsFilter) ([]domain.Run, int, error) {
	conditions := []string{}
	args := []any{}
	i := 1

	if f.UserID != "" {
		conditions = append(conditions, fmt.Sprintf("r.user_id = $%d", i))
		args = append(args, f.UserID)
		i++
	}

	if f.AssistantID != "" {
		conditions = append(conditions, fmt.Sprintf("r.assistant_id = $%d", i))
		args = append(args, f.AssistantID)
		i++
	}

	if f.Status != "" {
		conditions = append(conditions, fmt.Sprintf("r.status = $%d", i))
		args = append(args, f.Status)
		i++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM runs r %s`, where)
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("RunRepo.List count: %w", err)
	}

	offset := (f.Page - 1) * f.PageSize
	args = append(args, f.PageSize, offset)

	query := fmt.Sprintf(`
        SELECT
            r.id, r.assistant_id, COALESCE(a.name, ''),
            COALESCE(a.category_id::text, ''), COALESCE(c.name, ''),
            r.user_id, r.model, r.user_prompt,
            COALESCE(r.output, ''), r.status, COALESCE(r.error, ''),
            r.created_at
        FROM runs r
        LEFT JOIN assistants a ON a.id = r.assistant_id
        LEFT JOIN categories c ON c.id = a.category_id
        %s
        ORDER BY r.created_at DESC
        LIMIT $%d OFFSET $%d
    `, where, i, i+1)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("RunRepo.List: %w", err)
	}
	defer rows.Close()

	var runs []domain.Run
	for rows.Next() {
		var run domain.Run
		if err := rows.Scan(
			&run.ID, &run.AssistantID, &run.AssistantName,
			&run.CategoryID, &run.CategoryName,
			&run.UserID, &run.Model, &run.UserPrompt,
			&run.Output, &run.Status, &run.Error,
			&run.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("RunRepo.List scan: %w", err)
		}
		runs = append(runs, run)
	}

	return runs, total, nil
}
