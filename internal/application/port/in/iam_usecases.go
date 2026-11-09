// Package in holds the driving (inbound) ports — what this service offers,
// as interfaces the HTTP adapter depends on instead of the concrete use
// case structs directly (rules/2-anexos/C-api-hexagonal.md, "Puertos de
// entrada").
package in

import "context"

// AuthorizeUseCase answers the one question every other service is meant
// to call before a protected write (ADR-012): does this administrator's
// role grant this permission?
type AuthorizeUseCase interface {
	Execute(ctx context.Context, administratorID, permission string) (allowed bool, err error)
}

// AssignRoleUseCase assigns a Role to an Administrator.
type AssignRoleUseCase interface {
	Execute(ctx context.Context, administratorID, roleName string) error
}

// ListRolesUseCase lists every known Role and its permissions.
type ListRolesUseCase interface {
	Execute(ctx context.Context) ([]RoleView, error)
}

// RoleView is the read shape ListRolesUseCase returns — the application
// layer's own view, not the domain.Role struct leaking out directly.
type RoleView struct {
	Name        string
	Permissions []string
}
