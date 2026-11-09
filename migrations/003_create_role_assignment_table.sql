-- Maps one Administrator (by ID) to exactly one Role. administrator_id is a
-- real UUID (not text) because it originates from a Postgres UUID column
-- (access.administrator.id) on the same shared instance — unlike
-- lms-workflow's loan_id, which comes from MongoDB and must stay text.
--
-- A missing row here is the expected, backward-compatible default (ADR-012:
-- no assignment = full ADMINISTRATOR access), not an error condition — this
-- table only ever holds administrators that were explicitly, deliberately
-- assigned a role.
--
-- role_name -> iam.role(name) is a real foreign key (not a cross-domain
-- reference, rules/A-db-postgres.md rule 8 does not apply here — both
-- tables belong to this same iam domain).
CREATE TABLE iam.role_assignment (
  administrator_id  UUID        PRIMARY KEY,
  role_name         text        NOT NULL,

  created_at        timestamptz NOT NULL DEFAULT now(),
  updated_at        timestamptz NOT NULL DEFAULT now()
);
