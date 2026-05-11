package assistants

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RodionSOK/ai-assistants/internal/domain"
	"github.com/RodionSOK/ai-assistants/internal/httputil"
	mw "github.com/RodionSOK/ai-assistants/internal/middleware"
	"github.com/RodionSOK/ai-assistants/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type AssistantUsecase interface {
	List(ctx context.Context, in usecase.ListAssistantsInput) (*usecase.ListAssistantsOutput, error)
	GetByID(ctx context.Context, id string) (*domain.Assistant, error)
	Create(ctx context.Context, in usecase.CreateAssistantInput) (*domain.Assistant, error)
	Update(ctx context.Context, in usecase.UpdateAssistantInput) (*domain.Assistant, error)
}

type Handler struct {
	uc AssistantUsecase
}

func New(uc AssistantUsecase) *Handler {
	return &Handler{uc: uc}
}

type createRequest struct {
	CategoryID        string `json:"categoryId"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Model             string `json:"model"`
	SystemPrompt      string `json:"systemPrompt"`
	ExampleUserPrompt string `json:"exampleUserPrompt"`
	IsActive          *bool  `json:"isActive"`
}

type updateRequest struct {
	CategoryID        string `json:"categoryId"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Model             string `json:"model"`
	SystemPrompt      string `json:"systemPrompt"`
	ExampleUserPrompt string `json:"exampleUserPrompt"`
	IsActive          bool   `json:"isActive"`
}

type assistantResponse struct {
	ID                string `json:"id"`
	CategoryID        string `json:"categoryId"`
	CategoryName      string `json:"categoryName,omitempty"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Model             string `json:"model"`
	SystemPrompt      string `json:"systemPrompt,omitempty"`
	ExampleUserPrompt string `json:"exampleUserPrompt,omitempty"`
	IsActive          bool   `json:"isActive"`
	CreatedAt         string `json:"createdAt,omitempty"`
	UpdatedAt         string `json:"updatedAt,omitempty"`
}

type paginationResponse struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

type listResponse struct {
	Assistants []assistantResponse `json:"assistants"`
	Pagination paginationResponse  `json:"pagination"`
}

func toAssistantResponse(a domain.Assistant, isAdmin bool) assistantResponse {
	resp := assistantResponse{
		ID:                a.ID,
		CategoryID:        a.CategoryID,
		CategoryName:      a.CategoryName,
		Name:              a.Name,
		Description:       a.Description,
		Model:             a.Model,
		ExampleUserPrompt: a.ExampleUserPrompt,
		IsActive:          a.IsActive,
		CreatedAt:         a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         a.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if isAdmin {
		resp.SystemPrompt = a.SystemPrompt
	}

	return resp
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	role := mw.GetRole(r.Context())
	isAdmin := role == "admin"

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	includeInactive, _ := strconv.ParseBool(r.URL.Query().Get("includeInactive"))

	out, err := h.uc.List(r.Context(), usecase.ListAssistantsInput{
		CategoryID:      r.URL.Query().Get("categoryId"),
		Query:           r.URL.Query().Get("q"),
		IncludeInactive: includeInactive,
		IsAdmin:         isAdmin,
		Page:            page,
		PageSize:        pageSize,
	})
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get assistants")
		return
	}

	resp := listResponse{
		Assistants: []assistantResponse{},
		Pagination: paginationResponse{
			Page:     out.Page,
			PageSize: out.PageSize,
			Total:    out.Total,
		},
	}

	for _, a := range out.Assistants {
		resp.Assistants = append(resp.Assistants, toAssistantResponse(a, isAdmin))
	}

	httputil.JSON(w, http.StatusOK, resp)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "assistantId")
	role := mw.GetRole(r.Context())
	isAdmin := role == "admin"

	assistant, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "ASSISTANT_NOT_FOUND", "assistant not found")
		return
	}

	httputil.JSON(w, http.StatusOK, toAssistantResponse(*assistant, isAdmin))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	if req.Name == "" || req.Description == "" || req.Model == "" || req.SystemPrompt == "" || req.CategoryID == "" {
		httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "name, description, model, systemPrompt and categoryId are required")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	assistant, err := h.uc.Create(r.Context(), usecase.CreateAssistantInput{
		CategoryID:        req.CategoryID,
		Name:              req.Name,
		Description:       req.Description,
		Model:             req.Model,
		SystemPrompt:      req.SystemPrompt,
		ExampleUserPrompt: req.ExampleUserPrompt,
		IsActive:          isActive,
	})
	if err != nil {
		switch err.Error() {
		case "category not found":
			httputil.Error(w, http.StatusBadRequest, "CATEGORY_NOT_FOUND", "category not found")
		case "system prompt is required":
			httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "system prompt is required")
		default:
			httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create assistant")
		}
		return
	}

	httputil.JSON(w, http.StatusCreated, toAssistantResponse(*assistant, true))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "assistantId")

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	if req.Name == "" || req.Description == "" || req.Model == "" || req.SystemPrompt == "" || req.CategoryID == "" {
		httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "name, description, model, systemPrompt and categoryId are required")
		return
	}

	assistant, err := h.uc.Update(r.Context(), usecase.UpdateAssistantInput{
		ID:                id,
		CategoryID:        req.CategoryID,
		Name:              req.Name,
		Description:       req.Description,
		Model:             req.Model,
		SystemPrompt:      req.SystemPrompt,
		ExampleUserPrompt: req.ExampleUserPrompt,
		IsActive:          req.IsActive,
	})
	if err != nil {
		switch err.Error() {
		case "assistant not found":
			httputil.Error(w, http.StatusNotFound, "ASSISTANT_NOT_FOUND", "assistant not found")
		case "category not found":
			httputil.Error(w, http.StatusBadRequest, "CATEGORY_NOT_FOUND", "category not found")
		default:
			httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update assistant")
		}
		return
	}

	httputil.JSON(w, http.StatusOK, toAssistantResponse(*assistant, true))
}
