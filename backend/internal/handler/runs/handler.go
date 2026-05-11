package runs

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

type RunUsecase interface {
	Run(ctx context.Context, assistantID, userID, userPrompt string) (*domain.Run, error)
	My(ctx context.Context, in usecase.ListRunsInput) (*usecase.ListRunsOutput, error)
	All(ctx context.Context, in usecase.ListRunsInput) (*usecase.ListRunsOutput, error)
}

type Handler struct {
	uc RunUsecase
}

func New(uc RunUsecase) *Handler {
	return &Handler{uc: uc}
}

type runRequest struct {
	UserPrompt string `json:"userPrompt"`
}

type runResponse struct {
	ID            string `json:"id"`
	AssistantID   string `json:"assistantId"`
	AssistantName string `json:"assistantName,omitempty"`
	CategoryID    string `json:"categoryId,omitempty"`
	CategoryName  string `json:"categoryName,omitempty"`
	UserID        string `json:"userId"`
	Model         string `json:"model"`
	UserPrompt    string `json:"userPrompt"`
	Output        string `json:"output,omitempty"`
	Status        string `json:"status"`
	Error         string `json:"error,omitempty"`
	CreatedAt     string `json:"createdAt,omitempty"`
}

type paginationResponse struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

type listResponse struct {
	Runs       []runResponse      `json:"runs"`
	Pagination paginationResponse `json:"pagination"`
}

func toRunResponse(r domain.Run) runResponse {
	return runResponse{
		ID:            r.ID,
		AssistantID:   r.AssistantID,
		AssistantName: r.AssistantName,
		CategoryID:    r.CategoryID,
		CategoryName:  r.CategoryName,
		UserID:        r.UserID,
		Model:         r.Model,
		UserPrompt:    r.UserPrompt,
		Output:        r.Output,
		Status:        string(r.Status),
		Error:         r.Error,
		CreatedAt:     r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *Handler) Run(w http.ResponseWriter, r *http.Request) {
	assistantID := chi.URLParam(r, "assistantId")
	userID := mw.GetUserID(r.Context())

	var req runRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	if req.UserPrompt == "" {
		httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "userPrompt is required")
		return
	}

	run, err := h.uc.Run(r.Context(), assistantID, userID, req.UserPrompt)
	if err != nil {
		switch err.Error() {
		case "assistant not found":
			httputil.Error(w, http.StatusNotFound, "ASSISTANT_NOT_FOUND", "assistant not found")
		case "assistant is inactive":
			httputil.Error(w, http.StatusConflict, "ASSISTANT_INACTIVE", "assistant is inactive")
		default:
			if run != nil {
				// LLM ошибка — запуск сохранён со статусом failed
				httputil.JSON(w, http.StatusBadGateway, map[string]any{
					"error": map[string]any{
						"code":    "LLM_PROVIDER_ERROR",
						"message": err.Error(),
					},
				})
				return
			}
			httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to run assistant")
		}
		return
	}

	httputil.JSON(w, http.StatusCreated, toRunResponse(*run))
}

func (h *Handler) My(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r.Context())

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	out, err := h.uc.My(r.Context(), usecase.ListRunsInput{
		UserID:   userID,
		Status:   r.URL.Query().Get("status"),
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get runs")
		return
	}

	resp := listResponse{
		Runs: []runResponse{},
		Pagination: paginationResponse{
			Page:     out.Page,
			PageSize: out.PageSize,
			Total:    out.Total,
		},
	}
	for _, run := range out.Runs {
		resp.Runs = append(resp.Runs, toRunResponse(run))
	}

	httputil.JSON(w, http.StatusOK, resp)
}

func (h *Handler) All(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	out, err := h.uc.All(r.Context(), usecase.ListRunsInput{
		AssistantID: r.URL.Query().Get("assistantId"),
		Status:      r.URL.Query().Get("status"),
		Page:        page,
		PageSize:    pageSize,
	})
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get runs")
		return
	}

	resp := listResponse{
		Runs: []runResponse{},
		Pagination: paginationResponse{
			Page:     out.Page,
			PageSize: out.PageSize,
			Total:    out.Total,
		},
	}
	for _, run := range out.Runs {
		resp.Runs = append(resp.Runs, toRunResponse(run))
	}

	httputil.JSON(w, http.StatusOK, resp)
}
