package repository

import (
	"context"
	"fmt"

	"github.com/RodionSOK/ai-assistants/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, user domain.User) (*domain.User, error) {
	var created domain.User
	err := r.pool.QueryRow(ctx, `
        INSERT INTO users (id, email, password, role)
        VALUES ($1, $2, $3, $4)
        RETURNING id, email, role, created_at
    `,
		user.ID, user.Email, user.Password, user.Role,
	).Scan(&created.ID, &created.Email, &created.Role, &created.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("UserRepo.Create: %w", err)
	}

	return &created, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx, `
        SELECT id, email, role, created_at
        FROM users
        WHERE id = $1
    `, id).Scan(&u.ID, &u.Email, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("UserRepo.GetByID: %w", err)
	}

	return &u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx, `
        SELECT id, email, password, role, created_at
        FROM users
        WHERE email = $1
    `, email).Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("UserRepo.GetByEmail: %w", err)
	}

	return &u, nil
}

func (r *UserRepo) ExistsByID(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
        SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)
    `, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("UserRepo.ExistsByID: %w", err)
	}

	return exists, nil
}
