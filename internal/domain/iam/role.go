// Package iam implements the identity-and-access-management bounded context
// (ADR-012-iam-service-scope.md): Role definitions and the assignment of a
// Role to an Administrator. It does not model Administrator itself — that
// stays lms-access-api's aggregate; this package only ever refers to an
// administrator by its ID (rules/A-db-postgres.md, rule 8 — a cross-domain
// reference is not a DB-level FK).
package iam

// WildcardPermission grants every permission — used by the seeded
// ADMINISTRATOR role and by the backward-compatible "no assignment" default
// (ADR-012) so standing this service up never locks an existing
// Administrator out.
const WildcardPermission = "*"

// Role is a name and the set of permissions it grants.
type Role struct {
	Name        string
	Permissions []string
}

// Allows reports whether this role grants the given permission — either
// directly, or via the wildcard.
func (r *Role) Allows(permission string) bool {
	for _, p := range r.Permissions {
		if p == WildcardPermission || p == permission {
			return true
		}
	}
	return false
}

// Assignment maps one Administrator (by ID) to exactly one Role.
type Assignment struct {
	AdministratorID string
	RoleName        string
}
