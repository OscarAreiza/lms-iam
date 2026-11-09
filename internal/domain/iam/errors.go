package iam

import "errors"

// ErrRoleNotFound is the port-level "not found" error for a Role lookup —
// use cases depend on this, never on an infrastructure package's own error
// type (library-docs/05-architecture/hexagonal-architecture.md, Dependency
// Rule).
var ErrRoleNotFound = errors.New("role not found")
