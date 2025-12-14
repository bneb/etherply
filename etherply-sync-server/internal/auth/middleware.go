package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
)

// Middleware performs strict bearer token validation.
// It verifies the JWT signature using the server's shared secret.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow CORS for the demo
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Upgrade, Connection")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Skip auth for WebSocket upgrade handling if strictly needed,
		// but usually we pass token in query param or header.
		// For MVP simplicuty, let's just log it.

		// Strict Auth: We demand a Bearer token.
		authHeader := r.Header.Get("Authorization")

		// 1. Check Header
		var token string
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// 2. Fallback: Check Query Param (common for WebSockets)
		if token == "" {
			token = r.URL.Query().Get("token")
		}

		// 3. Validation
		if token == "" {
			// STRICT MODE: No token = 401. No more "Demo Convenience".
			slog.Warn("auth_rejected", "reason", "missing_token")
			http.Error(w, "Unauthorized: Bearer token required", http.StatusUnauthorized)
			return
		}

		claims, err := ValidateToken(token)
		if err != nil {
			slog.Warn("auth_rejected", "reason", "invalid_token", "error", err)
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// 4. Extract Scopes (Optional)
		// We add "scopes" to the request context for downstream handlers
		// "scope" claim is usually space-separated string
		// Default to empty if missing
		var scopes []string
		if scopeStr, ok := claims["scope"].(string); ok {
			scopes = strings.Split(scopeStr, " ")
		}

		// Inject into context
		ctx := r.Context()
		// We define context keys in a safer way usually, but for now strings.
		// Wait, context keys should be unexported types.
		// But handler is in different package.
		// Let's create a Helper in `auth` to get/set scopes?
		// Or just pass it?
		// Ideally `auth.ContextWithScopes`?
		ctx = NewContextWithScopes(ctx, scopes)

		// Pass execution to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Context keys
type contextKey string

const scopeContextKey contextKey = "scopes"

func NewContextWithScopes(ctx context.Context, scopes []string) context.Context {
	return context.WithValue(ctx, scopeContextKey, scopes)
}

func ScopesFromContext(ctx context.Context) []string {
	if scopes, ok := ctx.Value(scopeContextKey).([]string); ok {
		return scopes
	}
	return []string{}
}
