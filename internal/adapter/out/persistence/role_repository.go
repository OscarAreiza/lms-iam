package persistence

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/code-corhuila/lms-iam/internal/domain/iam"
)

// RoleRepository implements out.RoleRepository against PostgreSQL — the
// only code allowed to touch the `iam.role` table
// (library-docs/09-microservices/service-boundary-rules.md).
type RoleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) FindByName(ctx context.Context, name string) (*iam.Role, error) {
	const query = `SELECT name, permissions FROM iam.role WHERE name = $1`

	var role iam.Role
	err := r.db.QueryRow(ctx, query, name).Scan(&role.Name, &role.Permissions)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, iam.ErrRoleNotFound
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) FindAll(ctx context.Context) ([]*iam.Role, error) {
	const query = `SELECT name, permissions FROM iam.role ORDER BY name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*iam.Role
	for rows.Next() {
		var role iam.Role
		if err := rows.Scan(&role.Name, &role.Permissions); err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}
	return roles, rows.Err()
}
