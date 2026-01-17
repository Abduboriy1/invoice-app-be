// internal/interfaces/http/middleware/auth.go
package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/invoice-app-be/internal/infrastructure/auth"
)

type contextKey string

const userIDKey contextKey = "userID"

type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{jwtManager: jwtManager}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("[AUTH] ERROR: Missing authorization header")
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			log.Printf("[AUTH] ERROR: Expected 2 parts, got %d. Header: '%s'", len(parts), authHeader)
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		if parts[0] != "Bearer" {
			log.Printf("[AUTH] ERROR: Expected 'Bearer', got '%s'", parts[0])
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		claims, err := m.jwtManager.Verify(token)
		if err != nil {
			log.Printf("[AUTH] ERROR: Token verification failed: %v", err)
			log.Printf("[AUTH] Token that failed: %s... (length: %d)", truncate(token, 20), len(token))
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper function to safely truncate strings for logging
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

func GetUserIDFromContext(ctx context.Context) uuid.UUID {
	userID, _ := ctx.Value(userIDKey).(uuid.UUID)
	return userID
}
