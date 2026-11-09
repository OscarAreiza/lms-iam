package handler

import (
	"net/http"

	"github.com/code-corhuila/lms-iam/internal/adapter/in/httpapi/middleware"
	"github.com/code-corhuila/lms-iam/internal/adapter/in/httpapi/response"
	"github.com/code-corhuila/lms-iam/internal/application/port/in"
)

// RolesHandler implements GET /api/v1/iam/roles.
type RolesHandler struct {
	listRoles in.ListRolesUseCase
}

func NewRolesHandler(listRoles in.ListRolesUseCase) *RolesHandler {
	return &RolesHandler{listRoles: listRoles}
}

type roleResponse struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

func (h *RolesHandler) List(w http.ResponseWriter, r *http.Request) {
	traceID := middleware.FromContext(r.Context())

	roles, err := h.listRoles.Execute(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", traceID)
		return
	}

	items := make([]roleResponse, 0, len(roles))
	for _, rl := range roles {
		items = append(items, roleResponse{Name: rl.Name, Permissions: rl.Permissions})
	}
	response.JSON(w, http.StatusOK, items)
}
