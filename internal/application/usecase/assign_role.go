package usecase

import (
	"context"

	"github.com/code-corhuila/lms-iam/internal/application/port/out"
	"github.com/code-corhuila/lms-iam/internal/domain/iam"
)

// AssignRole implements POST /api/v1/iam/assignments. The role must already
// exist — this use case never creates one implicitly.
type AssignRole struct {
	Assignments out.AssignmentRepository
	Roles       out.RoleRepository
}

func NewAssignRole(assignments out.AssignmentRepository, roles out.RoleRepository) *AssignRole {
	return &AssignRole{Assignments: assignments, Roles: roles}
}

func (uc *AssignRole) Execute(ctx context.Context, administratorID, roleName string) error {
	if _, err := uc.Roles.FindByName(ctx, roleName); err != nil {
		return err // iam.ErrRoleNotFound if the role doesn't exist
	}

	return uc.Assignments.Save(ctx, &iam.Assignment{
		AdministratorID: administratorID,
		RoleName:        roleName,
	})
}
