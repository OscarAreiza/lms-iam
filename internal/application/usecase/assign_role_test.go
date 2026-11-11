package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/code-corhuila/lms-iam/internal/domain/iam"
)

func TestAssignRole_SavesAssignmentWhenRoleExists(t *testing.T) {
	assignments := newFakeAssignmentRepository()
	roles := newFakeRoleRepository(&iam.Role{Name: "LIBRARIAN", Permissions: []string{"books:write"}})
	uc := NewAssignRole(assignments, roles)

	if err := uc.Execute(context.Background(), "admin-1", "LIBRARIAN"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	saved, found := assignments.byAdministratorID["admin-1"]
	if !found {
		t.Fatal("expected the assignment to be saved")
	}
	if saved.RoleName != "LIBRARIAN" {
		t.Fatalf("expected RoleName LIBRARIAN, got %s", saved.RoleName)
	}
}

func TestAssignRole_RejectsUnknownRole(t *testing.T) {
	assignments := newFakeAssignmentRepository()
	roles := newFakeRoleRepository() // no roles seeded
	uc := NewAssignRole(assignments, roles)

	err := uc.Execute(context.Background(), "admin-1", "GHOST_ROLE")
	if !errors.Is(err, iam.ErrRoleNotFound) {
		t.Fatalf("expected iam.ErrRoleNotFound, got %v", err)
	}
	if len(assignments.saved) != 0 {
		t.Fatal("expected no assignment to be saved when the role does not exist")
	}
}

func TestAssignRole_PropagatesSaveError(t *testing.T) {
	assignments := newFakeAssignmentRepository()
	assignments.saveErr = errors.New("db unavailable")
	roles := newFakeRoleRepository(&iam.Role{Name: "LIBRARIAN", Permissions: []string{"books:write"}})
	uc := NewAssignRole(assignments, roles)

	if err := uc.Execute(context.Background(), "admin-1", "LIBRARIAN"); err == nil {
		t.Fatal("expected the repository's save error to propagate")
	}
}
