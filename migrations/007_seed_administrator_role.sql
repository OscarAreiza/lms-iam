-- Seed data (DML), not DDL — kept in its own numbered file so it's obviously
-- a separate changeset. Seeds the one role that already implicitly exists
-- today per ADR-012's backward-compatible default: every Administrator,
-- assigned or not, currently has full access. Explicitly assigning an
-- administrator to this role is a no-op in practice (the wildcard default
-- already grants everything), but it exists so AssignRole has a real role
-- to point at from day one.
INSERT INTO iam.role (name, permissions) VALUES ('ADMINISTRATOR', ARRAY['*']);
