// Package persistence holds the secondary (driven) adapters that implement
// this service's repository ports against PostgreSQL — the `iam` schema on
// the shared instance (Anexo J), migrated by this repo's own migrations/
// folder (ADR-012, same no-separate-`-db`-repo precedent as lms-workflow).
package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool opens a connection pool and verifies connectivity with a bounded timeout.
func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("creating connection pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return pool, nil
}
