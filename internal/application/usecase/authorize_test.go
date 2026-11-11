package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/code-corhuila/lms-iam/internal/domain/iam"
)

func TestAuthorize_NoAssignment_DefaultsAllowed(t *testing.T) {
	assignments := newFakeAssignmentRepository()
	roles := newFakeRoleRepository()
	uc := NewAuthorize(assignments, roles)

	allowed, err := uc.Execute(context.Background(), "admin-1", "books:write")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("expected the ADR-012 backward-compatible default (no assignment = allowed)")
	}
}

func TestAuthorize_AssignedRoleGrantsPermission(t *testing.T) {
	assignments := newFakeAssignmentRepository()
	assignments.byAdministratorID["admin-1"] = &iam.Assignment{AdministratorID: "admin-1", RoleName: "LIBRARIAN"}
	roles := newFakeRoleRepository(&iam.Role{Name: "LIBRARIAN", Permissions: []string{"books:write"}})
	uc := NewAuthorize(assignments, roles)

	allowed, err := uc.Execute(context.Background(), "admin-1", "books:write")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("expected the assigned role to grant this permission")
	}
}

func TestAuthorize_AssignedRoleDeniesPermission(t *testing.T) {
	assignments := newFakeAssignmentRepository()
	assignments.byAdministratorID["admin-1"] = &iam.Assignment{AdministratorID: "admin-1", RoleName: "LIBRARIAN"}
	roles := newFakeRoleRepository(&iam.Role{Name: "LIBRARIAN", Permissions: []string{"books:write"}})
	uc := NewAuthorize(assignments, roles)

	allowed, err := uc.Execute(context.Background(), "admin-1", "loans:delete")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("expected the assigned role to deny a permission it does not list")
	}
}

func TestAuthorize_AssignmentLookupError_PropagatesAsDenied(t *testing.T) {
	assignments := newFakeAssignmentRepository()
	assignments.findErr = errors.New("db unavailable")
	roles := newFakeRoleRepository()
	uc := NewAuthorize(assignments, roles)

	allowed, err := uc.Execute(context.Background(), "admin-1", "books:write")
	if err == nil {
		t.Fatal("expected the assignment lookup error to propagate")
	}
	if allowed {
		t.Fatal("expected allowed=false alongside a propagated error")
	}
}

func TestAuthorize_AssignedToUnknownRole_PropagatesError(t *testing.T) {
	assignments := newFakeAssignmentRepository()
	assignments.byAdministratorID["admin-1"] = &iam.Assignment{AdministratorID: "admin-1", RoleName: "GHOST_ROLE"}
	roles := newFakeRoleRepository() // GHOST_ROLE was never seeded
	uc := NewAuthorize(assignments, roles)

	_, err := uc.Execute(context.Background(), "admin-1", "books:write")
	if !errors.Is(err, iam.ErrRoleNotFound) {
		t.Fatalf("expected iam.ErrRoleNotFound, got %v", err)
	}
}
