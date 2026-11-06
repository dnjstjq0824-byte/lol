package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

type Claims struct {
	UserID uuid.UUID `json:"userId"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uuid.UUID, email, role, secret string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "bn-riftstats",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				// Anonymous user
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err == nil && token.Valid {
				ctx := context.WithValue(r.Context(), UserContextKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(UserContextKey).(*Claims)
		if !ok || claims == nil {
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"type":"https://api.riftstats.bn.io/errors/unauthorized","title":"Unauthorized","status":401,"detail":"Valid Bearer token required"}`))
			return
		}
		next(w, r)
	}
}

// RequireRole enforces RBAC: the caller must be authenticated and hold one
// of the given roles (e.g. SUPER_ADMIN, ADMIN).
func RequireRole(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*Claims)
			if !ok || claims == nil {
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"type":"https://api.riftstats.bn.io/errors/unauthorized","title":"Unauthorized","status":401,"detail":"Valid Bearer [REDACTED] required"}`))
				return
			}
			if _, ok := allowed[claims.Role]; !ok {
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"type":"https://api.riftstats.bn.io/errors/forbidden","title":"Forbidden","status":403,"detail":"This endpoint requires an administrator role"}`))
				return
			}
			next(w, r)
		}
	}
}

func GetUserFromContext(ctx context.Context) *Claims {
	if c, ok := ctx.Value(UserContextKey).(*Claims); ok {
		return c
	}
	return nil
}
