package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/RodionSOK/ai-assistants/internal/config"
	"github.com/RodionSOK/ai-assistants/internal/db"
	handlerAssistants "github.com/RodionSOK/ai-assistants/internal/handler/assistants"
	handlerAuth "github.com/RodionSOK/ai-assistants/internal/handler/auth"
	handlerCategories "github.com/RodionSOK/ai-assistants/internal/handler/categories"
	handlerRuns "github.com/RodionSOK/ai-assistants/internal/handler/runs"
	"github.com/RodionSOK/ai-assistants/internal/llm"
	mw "github.com/RodionSOK/ai-assistants/internal/middleware"
	"github.com/RodionSOK/ai-assistants/internal/repository"
	"github.com/RodionSOK/ai-assistants/internal/usecase"
)

func main() {
	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}
	log.Println("connected to database")

	if err := db.RunMigrations(cfg.PostgresDSN); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("migrations applied")

	userRepo := repository.NewUserRepo(pool)
	categoryRepo := repository.NewCategoryRepo(pool)
	assistantRepo := repository.NewAssistantRepo(pool)
	runRepo := repository.NewRunRepo(pool)

	llmProvider := llm.NewProvider(cfg)

	authUC := usecase.NewAuthUsecase(userRepo, cfg)
	categoryUC := usecase.NewCategoryUsecase(categoryRepo)
	assistantUC := usecase.NewAssistantUsecase(assistantRepo, categoryRepo)
	runUC := usecase.NewRunUsecase(runRepo, assistantRepo, llmProvider)

	authHandler := handlerAuth.New(authUC)
	categoryHandler := handlerCategories.New(categoryUC)
	assistantHandler := handlerAssistants.New(assistantUC)
	runHandler := handlerRuns.New(runUC)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.CORSAllowedOrigins},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Get("/_info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Post("/dummyLogin", authHandler.DummyLogin)
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(mw.Auth(cfg.JWTSecret))

		r.Get("/categories", categoryHandler.List)
		r.With(mw.RequireRole("admin")).Post("/categories", categoryHandler.Create)

		r.Get("/assistants", assistantHandler.List)
		r.Get("/assistants/{assistantId}", assistantHandler.Get)
		r.With(mw.RequireRole("admin")).Post("/assistants", assistantHandler.Create)
		r.With(mw.RequireRole("admin")).Put("/assistants/{assistantId}", assistantHandler.Update)

		r.Post("/assistants/{assistantId}/run", runHandler.Run)
		r.Get("/runs/my", runHandler.My)

		r.With(mw.RequireRole("admin")).Get("/admin/runs", runHandler.All)
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exited")
}
