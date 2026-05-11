package auth

import (
    "context"
    "encoding/json"
    "net/http"

    "github.com/RodionSOK/ai-assistants/internal/domain"
    "github.com/RodionSOK/ai-assistants/internal/httputil"
)

type AuthUsecase interface {
    DummyLogin(ctx context.Context, role domain.Role) (string, *domain.User, error)
    Register(ctx context.Context, email, password string) (string, *domain.User, error)
    Login(ctx context.Context, email, password string) (string, *domain.User, error)
}

type Handler struct {
    uc AuthUsecase
}

func New(uc AuthUsecase) *Handler {
    return &Handler{uc: uc}
}

type dummyLoginRequest struct {
    Role string `json:"role"`
}

type registerRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type loginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type userResponse struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    Role      string `json:"role"`
    CreatedAt string `json:"createdAt,omitempty"`
}

type tokenResponse struct {
    Token string       `json:"token"`
    User  userResponse `json:"user"`
}

func toUserResponse(u *domain.User) userResponse {
    return userResponse{
        ID:        u.ID,
        Email:     u.Email,
        Role:      string(u.Role),
        CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
    }
}

func (h *Handler) DummyLogin(w http.ResponseWriter, r *http.Request) {
    var req dummyLoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
        return
    }

    role := domain.Role(req.Role)
    if role != domain.RoleAdmin && role != domain.RoleUser {
        httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "role must be admin or user")
        return
    }

    token, user, err := h.uc.DummyLogin(r.Context(), role)
    if err != nil {
        httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
        return
    }

    httputil.JSON(w, http.StatusOK, tokenResponse{
        Token: token,
        User:  toUserResponse(user),
    })
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    var req registerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
        return
    }

    if req.Email == "" || req.Password == "" {
        httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "email and password are required")
        return
    }

    if len(req.Password) < 8 {
        httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "password must be at least 8 characters")
        return
    }

    token, user, err := h.uc.Register(r.Context(), req.Email, req.Password)
    if err != nil {
        httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
        return
    }

    httputil.JSON(w, http.StatusCreated, tokenResponse{
        Token: token,
        User:  toUserResponse(user),
    })
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    var req loginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
        return
    }

    if req.Email == "" || req.Password == "" {
        httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "email and password are required")
        return
    }

    token, user, err := h.uc.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        httputil.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid credentials")
        return
    }

    httputil.JSON(w, http.StatusOK, tokenResponse{
        Token: token,
        User:  toUserResponse(user),
    })
}