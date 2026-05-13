//go:build e2e

package usecase_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/RodionSOK/ai-assistants/internal/domain"
	"github.com/RodionSOK/ai-assistants/internal/llm"
	"github.com/RodionSOK/ai-assistants/internal/repository"
	"github.com/RodionSOK/ai-assistants/internal/usecase"
)

func setupPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		t.Skip("TEST_DSN not set, skipping E2E test")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	t.Cleanup(func() { pool.Close() })
	return pool
}

func cleanupE2E(t *testing.T, pool *pgxpool.Pool, runID, assistantID, categoryID, userID string) {
	t.Helper()
	ctx := context.Background()

	if runID != "" {
		pool.Exec(ctx, `DELETE FROM runs WHERE id = $1`, runID)
	}
	if assistantID != "" {
		pool.Exec(ctx, `DELETE FROM assistants WHERE id = $1`, assistantID)
	}
	if categoryID != "" {
		pool.Exec(ctx, `DELETE FROM categories WHERE id = $1`, categoryID)
	}
	if userID != "" {
		pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	}
}

// admin создал категорию → admin создал ассистента → user запустил ассистента.
func TestE2E_AdminCreateCategory_AdminCreateAssistant_UserRunAssistant(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()

	// репозитории
	categoryRepo := repository.NewCategoryRepo(pool)
	assistantRepo := repository.NewAssistantRepo(pool)
	runRepo := repository.NewRunRepo(pool)
	userRepo := repository.NewUserRepo(pool)

	// юзкейсы
	categoryUC := usecase.NewCategoryUsecase(categoryRepo)
	assistantUC := usecase.NewAssistantUsecase(assistantRepo, categoryRepo)
	runUC := usecase.NewRunUsecase(runRepo, assistantRepo, llm.NewMockProvider())

	// трекинг созданных ID для cleanup
	var runID, assistantID, categoryID, userID string
	t.Cleanup(func() {
		cleanupE2E(t, pool, runID, assistantID, categoryID, userID)
	})

	// 1. Создаём тестового пользователя (user)
	hashed, err := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	user, err := userRepo.Create(ctx, domain.User{
		ID:       uuid.New().String(),
		Email:    "e2e-user-" + uuid.New().String()[:8] + "@test.com",
		Password: string(hashed),
		Role:     domain.RoleUser,
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	userID = user.ID

	// 2. Admin создаёт категорию
	category, err := categoryUC.Create(ctx, "E2E Category "+uuid.New().String()[:8], "test")
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	categoryID = category.ID

	// 3. Admin создаёт ассистента в этой категории
	assistant, err := assistantUC.Create(ctx, usecase.CreateAssistantInput{
		CategoryID:   categoryID,
		Name:         "E2E Assistant",
		Description:  "test assistant",
		Model:        "mock",
		SystemPrompt: "You are a helpful assistant.",
		IsActive:     true,
	})
	if err != nil {
		t.Fatalf("create assistant: %v", err)
	}
	assistantID = assistant.ID

	// 4. User запускает ассистента
	run, err := runUC.Run(ctx, assistantID, userID, "Hello!")
	if err != nil {
		t.Fatalf("run assistant: %v", err)
	}
	runID = run.ID

	// 5. Проверяем результат
	if run.Status != domain.RunStatusSuccess {
		t.Errorf("expected status 'success', got '%s'", run.Status)
	}
	if run.Output == "" {
		t.Error("expected non-empty output")
	}
	if run.AssistantID != assistantID {
		t.Errorf("expected assistantID '%s', got '%s'", assistantID, run.AssistantID)
	}
	if run.UserID != userID {
		t.Errorf("expected userID '%s', got '%s'", userID, run.UserID)
	}

	t.Logf("E2E passed: category=%s assistant=%s run=%s status=%s",
		categoryID, assistantID, runID, run.Status)
}
