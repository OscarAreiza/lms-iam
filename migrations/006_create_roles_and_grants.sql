-- NOLOGIN: these carry permissions only. The infrastructure creates the
-- actual login user from secrets and grants it one of these roles
-- (rules/A-db-postgres.md, rule 10).
CREATE ROLE iam_reader NOLOGIN;
CREATE ROLE iam_writer NOLOGIN;

GRANT USAGE ON SCHEMA iam TO iam_reader, iam_writer;

GRANT SELECT ON iam.role, iam.role_assignment TO iam_reader;
GRANT SELECT, INSERT, UPDATE ON iam.role, iam.role_assignment TO iam_writer;

-- Ties the permission-carrying role to the actual login user — the role
-- was NOLOGIN on purpose above; the login user itself is created by
-- lms-infra-postgres's init script (rules/3-Anexo-J, J.7).
GRANT iam_writer TO iam_app;
