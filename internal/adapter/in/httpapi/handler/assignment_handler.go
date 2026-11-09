package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/code-corhuila/lms-iam/internal/adapter/in/httpapi/middleware"
	"github.com/code-corhuila/lms-iam/internal/adapter/in/httpapi/response"
	"github.com/code-corhuila/lms-iam/internal/application/port/in"
	"github.com/code-corhuila/lms-iam/internal/domain/iam"
)

// AssignmentHandler implements POST /api/v1/iam/assignments.
type AssignmentHandler struct {
	assignRole in.AssignRoleUseCase
}

func NewAssignmentHandler(assignRole in.AssignRoleUseCase) *AssignmentHandler {
	return &AssignmentHandler{assignRole: assignRole}
}

type assignRoleRequest struct {
	AdministratorID string `json:"administratorId"`
	Role            string `json:"role"`
}

func (h *AssignmentHandler) Assign(w http.ResponseWriter, r *http.Request) {
	traceID := middleware.FromContext(r.Context())

	var req assignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", traceID)
		return
	}
	if req.AdministratorID == "" || req.Role == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "administratorId and role are both required", traceID)
		return
	}

	if err := h.assignRole.Execute(r.Context(), req.AdministratorID, req.Role); err != nil {
		if errors.Is(err, iam.ErrRoleNotFound) {
			response.Error(w, http.StatusNotFound, "ROLE_NOT_FOUND", "No role with that name exists", traceID)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", traceID)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"administratorId": req.AdministratorID,
		"role":            req.Role,
	})
}
