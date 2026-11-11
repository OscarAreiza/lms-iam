package iam

import "testing"

func TestRole_Allows_DirectPermission(t *testing.T) {
	role := &Role{Name: "LIBRARIAN", Permissions: []string{"books:write", "loans:read"}}

	if !role.Allows("books:write") {
		t.Fatal("expected role to allow a permission it directly lists")
	}
}

func TestRole_Allows_RejectsUnlistedPermission(t *testing.T) {
	role := &Role{Name: "LIBRARIAN", Permissions: []string{"books:write"}}

	if role.Allows("loans:delete") {
		t.Fatal("expected role to reject a permission it does not list")
	}
}

func TestRole_Allows_Wildcard(t *testing.T) {
	role := &Role{Name: "ADMINISTRATOR", Permissions: []string{WildcardPermission}}

	if !role.Allows("anything:at-all") {
		t.Fatal("expected the wildcard permission to allow any permission")
	}
}

func TestRole_Allows_EmptyPermissionsRejectsEverything(t *testing.T) {
	role := &Role{Name: "NONE", Permissions: nil}

	if role.Allows("books:write") {
		t.Fatal("expected a role with no permissions to allow nothing")
	}
}
