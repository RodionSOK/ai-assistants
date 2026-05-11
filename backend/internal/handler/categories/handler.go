package categories

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RodionSOK/ai-assistants/internal/domain"
	"github.com/RodionSOK/ai-assistants/internal/httputil"
)

type CategoryUsecase interface {
	List(ctx context.Context) ([]domain.Category, error)
	Create(ctx context.Context, name, description string) (*domain.Category, error)
}

type Handler struct {
	uc CategoryUsecase
}

func New(uc CategoryUsecase) *Handler {
	return &Handler{uc: uc}
}

type createRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type categoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
}

type listResponse struct {
	Categories []categoryResponse `json:"categories"`
}

func toCategoryResponse(c domain.Category) categoryResponse {
	return categoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.uc.List(r.Context())
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get categories")
		return
	}

	resp := listResponse{Categories: []categoryResponse{}}
	for _, c := range categories {
		resp.Categories = append(resp.Categories, toCategoryResponse(c))
	}

	httputil.JSON(w, http.StatusOK, resp)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	if req.Name == "" {
		httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "name is required")
		return
	}

	category, err := h.uc.Create(r.Context(), req.Name, req.Description)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create category")
		return
	}

	httputil.JSON(w, http.StatusCreated, toCategoryResponse(*category))
}
