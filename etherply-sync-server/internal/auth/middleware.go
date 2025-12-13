package auth

import (
	"log"
	"net/http"
	"strings"
)

// Middleware performs basic bearer token validation.
// For the MVP, it accepts any token that is not empty, representing a logged-in user.
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
			log.Println("[AUTH] Rejected request: Missing Authorization token.")
			http.Error(w, "Unauthorized: Bearer token required", http.StatusUnauthorized)
			return
		}

		if err := ValidateToken(token); err != nil {
			log.Printf("[AUTH] Rejected request: Invalid token (%v)", err)
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Pass execution to the next handler
		next.ServeHTTP(w, r)
	})
}
