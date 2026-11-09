-- FKs added after table creation, not inline (rules/A-db-postgres.md, rule 2).
-- RESTRICT: a role with existing assignments cannot be deleted out from
-- under them — an operator must reassign first.
ALTER TABLE iam.role_assignment
  ADD CONSTRAINT fk_role_assignment_role
  FOREIGN KEY (role_name) REFERENCES iam.role (name) ON DELETE RESTRICT;
