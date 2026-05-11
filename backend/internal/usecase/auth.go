package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/RodionSOK/ai-assistants/internal/config"
	"github.com/RodionSOK/ai-assistants/internal/domain"
	"github.com/RodionSOK/ai-assistants/internal/middleware"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
}

type AuthUsecase struct {
	userRepo UserRepository
	cfg      *config.Config
}

func NewAuthUsecase(userRepo UserRepository, cfg *config.Config) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo, cfg: cfg}
}

func (uc *AuthUsecase) DummyLogin(ctx context.Context, role domain.Role) (string, *domain.User, error) {
	var userID string
	var email string

	switch role {
	case domain.RoleAdmin:
		userID = uc.cfg.AdminUUID
		email = "admin@example.com"
	case domain.RoleUser:
		userID = uc.cfg.UserUUID
		email = "user@example.com"
	default:
		return "", nil, fmt.Errorf("invalid role: %s", role)
	}

	exists, err := uc.userRepo.ExistsByID(ctx, userID)
	if err != nil {
		return "", nil, fmt.Errorf("DummyLogin check exists: %w", err)
	}

	var user *domain.User
	if !exists {
		hashed, err := bcrypt.GenerateFromPassword([]byte("dummy"), bcrypt.DefaultCost)
		if err != nil {
			return "", nil, fmt.Errorf("DummyLogin hash: %w", err)
		}
		user, err = uc.userRepo.Create(ctx, domain.User{
			ID:       userID,
			Email:    email,
			Password: string(hashed),
			Role:     role,
		})
		if err != nil {
			return "", nil, fmt.Errorf("DummyLogin create: %w", err)
		}
	} else {
		user, err = uc.userRepo.GetByID(ctx, userID)
		if err != nil {
			return "", nil, fmt.Errorf("DummyLogin get: %w", err)
		}
	}

	token, err := uc.generateToken(user.ID, string(user.Role))
	if err != nil {
		return "", nil, fmt.Errorf("DummyLogin token: %w", err)
	}

	return token, user, nil
}

func (uc *AuthUsecase) Register(ctx context.Context, email, password string) (string, *domain.User, error) {
	existing, err := uc.userRepo.GetByEmail(ctx, email)
	if err == nil && existing != nil {
		return "", nil, fmt.Errorf("email already taken")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, fmt.Errorf("Register hash: %w", err)
	}

	user, err := uc.userRepo.Create(ctx, domain.User{
		ID:       uuid.New().String(),
		Email:    email,
		Password: string(hashed),
		Role:     domain.RoleUser,
	})
	if err != nil {
		return "", nil, fmt.Errorf("Register create: %w", err)
	}

	token, err := uc.generateToken(user.ID, string(user.Role))
	if err != nil {
		return "", nil, fmt.Errorf("Register token: %w", err)
	}

	return token, user, nil
}

func (uc *AuthUsecase) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	token, err := uc.generateToken(user.ID, string(user.Role))
	if err != nil {
		return "", nil, fmt.Errorf("Login token: %w", err)
	}

	return token, user, nil
}

func (uc *AuthUsecase) generateToken(userID, role string) (string, error) {
	claims := middleware.Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(uc.cfg.JWTExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.cfg.JWTSecret))
}
