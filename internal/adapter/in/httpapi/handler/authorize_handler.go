package handler

import (
	"net/http"

	"github.com/code-corhuila/lms-iam/internal/adapter/in/httpapi/middleware"
	"github.com/code-corhuila/lms-iam/internal/adapter/in/httpapi/response"
	"github.com/code-corhuila/lms-iam/internal/application/port/in"
)

// AuthorizeHandler implements GET /api/v1/iam/authorize (ADR-012) — the one
// endpoint every other service is meant to call before a protected write.
type AuthorizeHandler struct {
	authorize in.AuthorizeUseCase
}

func NewAuthorizeHandler(authorize in.AuthorizeUseCase) *AuthorizeHandler {
	return &AuthorizeHandler{authorize: authorize}
}

func (h *AuthorizeHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	traceID := middleware.FromContext(r.Context())

	administratorID := r.URL.Query().Get("administratorId")
	permission := r.URL.Query().Get("permission")
	if administratorID == "" || permission == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "administratorId and permission are both required", traceID)
		return
	}

	allowed, err := h.authorize.Execute(r.Context(), administratorID, permission)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", traceID)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{"allowed": allowed})
}
