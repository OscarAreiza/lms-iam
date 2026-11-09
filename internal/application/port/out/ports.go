// Package out holds the driven (outbound) ports every use case depends on —
// what the application needs from the outside world, never how it's
// implemented (rules/2-anexos/C-api-hexagonal.md, "Puertos de salida").
package out

import (
	"context"

	"github.com/code-corhuila/lms-iam/internal/domain/iam"
)

// RoleRepository is the driven port for Role persistence.
type RoleRepository interface {
	FindByName(ctx context.Context, name string) (*iam.Role, error)
	FindAll(ctx context.Context) ([]*iam.Role, error)
}

// AssignmentRepository is the driven port for Administrator-to-Role
// assignments.
type AssignmentRepository interface {
	// FindByAdministratorID returns iam.ErrRoleNotFound-free: a missing
	// assignment is a normal, expected outcome (the backward-compatible
	// default, ADR-012), not an error — callers check the returned bool.
	FindByAdministratorID(ctx context.Context, administratorID string) (*iam.Assignment, bool, error)
	Save(ctx context.Context, a *iam.Assignment) error
}
