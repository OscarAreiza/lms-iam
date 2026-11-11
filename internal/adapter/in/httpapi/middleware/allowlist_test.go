package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func withCaller(sub string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), administratorIDKey{}, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestRequireCaller_Allowed(t *testing.T) {
	allowed := RequireCaller("membership-service", "catalog-service", "circulation-service")

	for _, sub := range []string{"membership-service", "catalog-service", "circulation-service"} {
		handler := withCaller(sub)(allowed(okHandler()))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/iam/authorize", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("sub=%s: expected 200, got %d", sub, rec.Code)
		}
	}
}

func TestRequireCaller_Denied(t *testing.T) {
	handler := withCaller("lms-worker")(RequireCaller("membership-service", "catalog-service", "circulation-service")(okHandler()))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/iam/authorize", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

func TestRequireCaller_EmptySubDenied(t *testing.T) {
	handler := RequireCaller("membership-service")(okHandler())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/iam/authorize", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}
