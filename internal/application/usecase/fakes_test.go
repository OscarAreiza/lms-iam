package usecase

import (
	"context"

	"github.com/code-corhuila/lms-iam/internal/domain/iam"
)

// fakeAssignmentRepository is an in-memory out.AssignmentRepository double
// (tdd-guide.md's "Fakes" convention — no mocking library).
type fakeAssignmentRepository struct {
	byAdministratorID map[string]*iam.Assignment
	findErr           error
	saveErr           error
	saved             []*iam.Assignment
}

func newFakeAssignmentRepository() *fakeAssignmentRepository {
	return &fakeAssignmentRepository{byAdministratorID: make(map[string]*iam.Assignment)}
}

func (f *fakeAssignmentRepository) FindByAdministratorID(_ context.Context, administratorID string) (*iam.Assignment, bool, error) {
	if f.findErr != nil {
		return nil, false, f.findErr
	}
	a, found := f.byAdministratorID[administratorID]
	return a, found, nil
}

func (f *fakeAssignmentRepository) Save(_ context.Context, a *iam.Assignment) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.saved = append(f.saved, a)
	f.byAdministratorID[a.AdministratorID] = a
	return nil
}

// fakeRoleRepository is an in-memory out.RoleRepository double.
type fakeRoleRepository struct {
	byName  map[string]*iam.Role
	findErr error
	allErr  error
}

func newFakeRoleRepository(roles ...*iam.Role) *fakeRoleRepository {
	f := &fakeRoleRepository{byName: make(map[string]*iam.Role)}
	for _, r := range roles {
		f.byName[r.Name] = r
	}
	return f
}

func (f *fakeRoleRepository) FindByName(_ context.Context, name string) (*iam.Role, error) {
	if f.findErr != nil {
		return nil, f.findErr
	}
	r, ok := f.byName[name]
	if !ok {
		return nil, iam.ErrRoleNotFound
	}
	return r, nil
}

func (f *fakeRoleRepository) FindAll(_ context.Context) ([]*iam.Role, error) {
	if f.allErr != nil {
		return nil, f.allErr
	}
	roles := make([]*iam.Role, 0, len(f.byName))
	for _, r := range f.byName {
		roles = append(roles, r)
	}
	return roles, nil
}
