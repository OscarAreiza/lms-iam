# lms-iam

> Identity and access management: role definitions and role assignments

Part of the **LMS Library** distributed system — team `lms-library`, Grupo 2.
Governance and documentation live in [`library-docs`](https://github.com/code-corhuila/library-docs).

Go REST API implementing role-based authorization for the Administrator role
(`ADR-012-iam-service-scope.md`). Closes **T2/SEC-04** from the week-13 security threat model:
before this service existed, any valid token — including a 1-minute internal one — could call
every protected route of every service, because no roles or scopes existed anywhere.

This service does **not** re-implement login (that stays `lms-access-api`'s job per `ADR-011`),
does **not** issue the end-user session JWT, and does **not** touch the internal
service-to-service HS256 token scheme (`ADR-004`'s accepted gap, tracked separately as SEC-03).
It owns exactly one thing: *does this administrator's role allow this permission*.

## Structure

```
cmd/api/                       → entry point (main.go)
migrations/                    → this domain's own schema migrations (see "Migration scope" below)
internal/
├── domain/iam/                 → Role, Assignment — the aggregate and its invariants
├── application/
│   ├── port/in/                  → driving ports (AuthorizeUseCase, AssignRoleUseCase, ListRolesUseCase)
│   ├── port/out/                  → driven ports (RoleRepository, AssignmentRepository)
│   └── usecase/                    → Authorize, AssignRole, ListRoles
├── adapter/
│   ├── in/httpapi/                 → chi router, middleware (CORS, correlation id, RequireAuth), handlers
│   └── out/persistence/             → RoleRepository, AssignmentRepository (PostgreSQL)
├── config/                      → environment variable loading
└── infrastructure/logger/        → structured (zap) logger
```

## Key endpoints

- `GET /api/v1/iam/authorize?administratorId={id}&permission={perm}` → `200 {"allowed": bool}` —
  the one endpoint every other service is meant to call before a protected write
- `POST /api/v1/iam/assignments` → assign a role to an administrator
- `GET /api/v1/iam/roles` → list known roles and their permissions

**Backward-compatible default (ADR-012):** an administrator with no assignment row is treated as
the seeded `ADMINISTRATOR` role, wildcard permission `*` — always allowed. This preserves
today's real behavior (every valid token can do everything) so standing this service up never
locks anyone out; tightening it is a deliberate future decision, not a side effect of this one.

**Auth:** same dual-algorithm `RequireAuth` middleware every sibling Go service already uses —
RS256 (a real Administrator session, verified with `lms-access-api`'s public key) or HS256 (an
internal service-to-service call). This service never holds the RS256 private key.

## Internal caller allow-list (ADR-013)

`GET /api/v1/iam/authorize` requires `sub` in `{membership-service, catalog-service,
circulation-service}` (`middleware.RequireCaller`) — the three services expected to call it.
Closes **SEC-03** from the week-13 threat model (T1/T2) for this route.

## Known integration gap

**As of 2026-11-10, no service calling this on `develop` yet.** Each caller's own middleware PR
is still open (`lms-membership-api` #17/`feat/iam-authorization`, `lms-catalog-api` #15, same
branch name, `lms-circulation-api` #12) — the allow-list above anticipates their real `sub`
values (verified by reading those branches directly), but until those PRs merge, nothing on
`develop` actually calls `GET /api/v1/iam/authorize`. Flagged here so it isn't silently assumed
already wired up end-to-end.

## Migration scope

**From scratch** — the course never allocated this repo a scope note the way the original 18
repos have one; built per `ADR-012` once the repo existed, same situation `lms-workflow` was in
under `ADR-011`.

This service's `iam.role`/`iam.role_assignment` tables live in this repo's own `migrations/`
folder, targeting an `iam` schema on the same shared PostgreSQL instance every other domain's
schema lives on (Anexo J) — not a new database server, just a new schema. No separate
`lms-iam-db` repo exists; the course's fixed 18-repo topology never allocated one, same reasoning
`ADR-011` already used for `lms-workflow`.

## Development

```bash
go mod download
go run ./cmd/api/...
```

## Tests

```bash
go test ./...
```

---

## Branching

Three permanent branches. **None of them accepts a direct commit** — you enter through a child
branch and leave through a Pull Request.

```
develop  <--PR--  feat/... fix/... chore/...
qa       <--PR--  qa/...
main     <--PR--  release/...  hotfix/...
```

Promotion happens **by re-application** (`git cherry-pick -x`), never by merging one permanent
branch into another: `merge develop -> qa` and `merge qa -> main` do not exist in this model.

`main` requires **1 approval from `ariel5253`**. On `develop` and `qa` the team sets its own review
rule.

Full policy: `00-governance/branching-policy.md` in `library-docs`.

## Correlations

* Scope and design → `library-docs/05-architecture/decisions/records/ADR-012-iam-service-scope.md`
* The gap this closes → `sistemas-distribuidos-2026-b-g2` (W13 fork), `THREAT-MODEL.md`, T2/SEC-04
* Repo decomposition → `library-docs/05-architecture/decisions/records/ADR-006-repo-per-context-decomposition.md`
