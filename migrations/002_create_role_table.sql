-- A Role is a name and the set of permissions it grants (ADR-012).
-- permissions is text[], not a join table — the permission set is small and
-- always read/written as a whole alongside its role, never queried
-- independently; a wildcard entry ('*') grants everything.
CREATE TABLE iam.role (
  name          text        PRIMARY KEY CHECK (char_length(name) BETWEEN 1 AND 50),
  permissions   text[]      NOT NULL,

  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now()
);
