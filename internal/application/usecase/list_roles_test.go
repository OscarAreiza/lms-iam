package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/code-corhuila/lms-iam/internal/domain/iam"
)

func TestListRoles_ReturnsEveryRoleAsAView(t *testing.T) {
	roles := newFakeRoleRepository(
		&iam.Role{Name: "LIBRARIAN", Permissions: []string{"books:write"}},
		&iam.Role{Name: "ADMINISTRATOR", Permissions: []string{iam.WildcardPermission}},
	)
	uc := NewListRoles(roles)

	views, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(views) != 2 {
		t.Fatalf("expected 2 views, got %d", len(views))
	}
}

func TestListRoles_EmptyRepositoryReturnsEmptySlice(t *testing.T) {
	roles := newFakeRoleRepository()
	uc := NewListRoles(roles)

	views, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(views) != 0 {
		t.Fatalf("expected 0 views, got %d", len(views))
	}
}

func TestListRoles_PropagatesRepositoryError(t *testing.T) {
	roles := newFakeRoleRepository()
	roles.allErr = errors.New("db unavailable")
	uc := NewListRoles(roles)

	if _, err := uc.Execute(context.Background()); err == nil {
		t.Fatal("expected the repository's error to propagate")
	}
}
