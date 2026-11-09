-- Every FK column gets an index (rules/A-db-postgres.md, rule 3) — an
-- operator listing "who holds role X" would otherwise scan the whole table.
CREATE INDEX idx_role_assignment_role_name ON iam.role_assignment (role_name);
