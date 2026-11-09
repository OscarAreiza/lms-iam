package httpapi

import (
	"crypto/rsa"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/code-corhuila/lms-iam/internal/adapter/in/httpapi/handler"
	"github.com/code-corhuila/lms-iam/internal/adapter/in/httpapi/middleware"
)

// RouterConfig carries what the router needs to wire itself.
type RouterConfig struct {
	DB                *pgxpool.Pool
	JWTPublicKey      *rsa.PublicKey
	InternalJWTSecret string
	CORSOrigin        string
	Authorize         *handler.AuthorizeHandler
	Assignments       *handler.AssignmentHandler
	Roles             *handler.RolesHandler
}

// NewRouter builds the chi router with the base middleware stack, health
// endpoints, and the IAM bounded context's /iam routes (ADR-012).
func NewRouter(cfg RouterConfig) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CorrelationID)
	r.Use(middleware.CORS(cfg.CORSOrigin))

	health := handler.NewHealthHandler(cfg.DB)
	r.Get("/health", health.Liveness)
	r.Get("/health/ready", health.Readiness)

	r.Route("/api/v1/iam", func(api chi.Router) {
		api.Use(middleware.RequireAuth(cfg.JWTPublicKey, cfg.InternalJWTSecret))

		api.Get("/authorize", cfg.Authorize.Authorize)
		api.Post("/assignments", cfg.Assignments.Assign)
		api.Get("/roles", cfg.Roles.List)
	})

	return r
}
