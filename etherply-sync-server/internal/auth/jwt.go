// Package auth provides JWT-based authentication for the EtherPly sync server.
// It implements a strict "Fail Secure" policy: if ETHERPLY_JWT_SECRET is not
// configured, the server refuses to start. All requests (except CORS preflight)
// require a valid Bearer token in the Authorization header or 'token' query param.
//
// Usage:
//
//	auth.Init(os.Getenv("ETHERPLY_JWT_SECRET"))
//	handler := auth.Middleware(myHandler)
package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret     []byte
	isAuthEnabled bool
)

// Init configures the authentication layer.
// logic: If secret is empty, we default to blocking everything (Fail Secure),
// unless we explicitly want a "Dev Mode" (which we might add later, but for now strict).
func Init(secret string) {
	if secret != "" {
		jwtSecret = []byte(secret)
		isAuthEnabled = true
	}
}

// ValidateToken checks if the provided token string is valid signed JWT and returns the claims.
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	trimmed := strings.TrimSpace(tokenString)
	if trimmed == "" {
		return nil, errors.New("token is empty")
	}

	if !isAuthEnabled {
		return nil, errors.New("server authentication is not configured (ETHERPLY_JWT_SECRET missing)")
	}

	token, err := jwt.Parse(trimmed, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parse failed: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("token signature invalid")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid claims structure")
}
