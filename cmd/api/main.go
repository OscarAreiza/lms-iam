// Command api is the entry point for lms-iam — identity and access
// management (ADR-012-iam-service-scope.md).
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpapi "github.com/code-corhuila/lms-iam/internal/adapter/in/httpapi"
	"github.com/code-corhuila/lms-iam/internal/adapter/in/httpapi/handler"
	"github.com/code-corhuila/lms-iam/internal/adapter/out/persistence"
	"github.com/code-corhuila/lms-iam/internal/application/usecase"
	"github.com/code-corhuila/lms-iam/internal/config"
	"github.com/code-corhuila/lms-iam/internal/infrastructure/logger"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		println("config error:", err.Error())
		return err
	}

	zapLog, err := logger.New(cfg.LogLevel)
	if err != nil {
		return err
	}
	defer func() { _ = zapLog.Sync() }()
	log := zapLog.Sugar()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := persistence.NewPool(ctx, cfg.DSN())
	if err != nil {
		log.Errorw("failed to connect to database", "error", err)
		return err
	}
	defer pool.Close()

	roles := persistence.NewRoleRepository(pool)
	assignments := persistence.NewAssignmentRepository(pool)

	authorizeUseCase := usecase.NewAuthorize(assignments, roles)
	assignRoleUseCase := usecase.NewAssignRole(assignments, roles)
	listRolesUseCase := usecase.NewListRoles(roles)

	authorizeHandler := handler.NewAuthorizeHandler(authorizeUseCase)
	assignmentHandler := handler.NewAssignmentHandler(assignRoleUseCase)
	rolesHandler := handler.NewRolesHandler(listRolesUseCase)

	router := httpapi.NewRouter(httpapi.RouterConfig{
		DB:                pool,
		JWTPublicKey:      cfg.JWTPublicKey,
		InternalJWTSecret: cfg.InternalJWTSecret,
		CORSOrigin:        cfg.CORSOrigin,
		Authorize:         authorizeHandler,
		Assignments:       assignmentHandler,
		Roles:             rolesHandler,
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Infow("lms-iam listening", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorw("server error", "error", err)
		}
	}()

	<-ctx.Done()
	log.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
