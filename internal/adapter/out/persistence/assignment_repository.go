package persistence

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/code-corhuila/lms-iam/internal/domain/iam"
)

// AssignmentRepository implements out.AssignmentRepository against
// PostgreSQL — the only code allowed to touch the
// `iam.role_assignment` table.
type AssignmentRepository struct {
	db *pgxpool.Pool
}

func NewAssignmentRepository(db *pgxpool.Pool) *AssignmentRepository {
	return &AssignmentRepository{db: db}
}

// FindByAdministratorID returns found=false (no error) when no assignment
// exists — a missing assignment is the expected, backward-compatible
// default (ADR-012), not a failure.
func (r *AssignmentRepository) FindByAdministratorID(ctx context.Context, administratorID string) (*iam.Assignment, bool, error) {
	const query = `SELECT administrator_id, role_name FROM iam.role_assignment WHERE administrator_id = $1`

	var a iam.Assignment
	err := r.db.QueryRow(ctx, query, administratorID).Scan(&a.AdministratorID, &a.RoleName)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &a, true, nil
}

func (r *AssignmentRepository) Save(ctx context.Context, a *iam.Assignment) error {
	const query = `
		INSERT INTO iam.role_assignment (administrator_id, role_name, created_at, updated_at)
		VALUES ($1, $2, now(), now())
		ON CONFLICT (administrator_id) DO UPDATE SET
			role_name = EXCLUDED.role_name,
			updated_at = now()`

	_, err := r.db.Exec(ctx, query, a.AdministratorID, a.RoleName)
	return err
}
