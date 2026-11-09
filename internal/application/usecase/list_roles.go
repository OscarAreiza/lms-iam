package usecase

import (
	"context"

	"github.com/code-corhuila/lms-iam/internal/application/port/in"
	"github.com/code-corhuila/lms-iam/internal/application/port/out"
)

// ListRoles implements GET /api/v1/iam/roles — operator visibility into
// every known role and what it grants.
type ListRoles struct {
	Roles out.RoleRepository
}

func NewListRoles(roles out.RoleRepository) *ListRoles {
	return &ListRoles{Roles: roles}
}

func (uc *ListRoles) Execute(ctx context.Context) ([]in.RoleView, error) {
	roles, err := uc.Roles.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	views := make([]in.RoleView, 0, len(roles))
	for _, r := range roles {
		views = append(views, in.RoleView{Name: r.Name, Permissions: r.Permissions})
	}
	return views, nil
}
