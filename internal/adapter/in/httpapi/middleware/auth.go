package middleware

import (
	"context"
	"crypto/rsa"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type administratorIDKey struct{}

// RequireAuth accepts exactly two algorithms, each with its own dedicated
// key — never "whatever alg the token declares" against one shared secret
// (rules/2-anexos/C-api-hexagonal.md, numeral 5.3.7): RS256 (a real
// Administrator session, verified with lms-access-api's public key) or
// HS256 (an internal service-to-service call — every other service calling
// GET /api/v1/iam/authorize uses this). Same middleware every sibling Go
// service already uses, ported here unchanged.
func RequireAuth(publicKey *rsa.PublicKey, internalSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				writeUnauthorized(w)
				return
			}

			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(parts[1], claims, func(t *jwt.Token) (interface{}, error) {
				switch t.Method.(type) {
				case *jwt.SigningMethodRSA:
					return publicKey, nil
				case *jwt.SigningMethodHMAC:
					return []byte(internalSecret), nil
				default:
					return nil, jwt.ErrTokenSignatureInvalid
				}
			}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg(), jwt.SigningMethodHS256.Alg()}))
			if err != nil || !token.Valid {
				writeUnauthorized(w)
				return
			}

			administratorID, _ := claims["sub"].(string)
			ctx := context.WithValue(r.Context(), administratorIDKey{}, administratorID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CallerID retrieves the authenticated caller's sub (an Administrator ID
// for an RS256 token, or a service name for an internal HS256 token) from context.
func CallerID(ctx context.Context) string {
	id, _ := ctx.Value(administratorIDKey{}).(string)
	return id
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":"UNAUTHORIZED","message":"Authentication token required"}`))
}
