// Package usecase orchestrates the domain for each use case — it contains
// no business logic itself (that lives in the Aggregate), only coordination
// (library-docs/05-architecture/hexagonal-architecture.md).
package usecase

import (
	"context"

	"github.com/code-corhuila/lms-iam/internal/application/port/out"
)

// Authorize implements ADR-012's core decision. An administrator with no
// assignment row defaults to allowed=true (the backward-compatible
// default: today, every valid token can do everything, and standing this
// service up must not silently lock anyone out).
type Authorize struct {
	Assignments out.AssignmentRepository
	Roles       out.RoleRepository
}

func NewAuthorize(assignments out.AssignmentRepository, roles out.RoleRepository) *Authorize {
	return &Authorize{Assignments: assignments, Roles: roles}
}

func (uc *Authorize) Execute(ctx context.Context, administratorID, permission string) (bool, error) {
	assignment, found, err := uc.Assignments.FindByAdministratorID(ctx, administratorID)
	if err != nil {
		return false, err
	}
	if !found {
		return true, nil // ADR-012 backward-compatible default
	}

	role, err := uc.Roles.FindByName(ctx, assignment.RoleName)
	if err != nil {
		return false, err
	}

	return role.Allows(permission), nil
}
