package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret []byte
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

// ValidateToken checks if the provided token string is valid signed JWT.
func ValidateToken(tokenString string) error {
	trimmed := strings.TrimSpace(tokenString)
	if trimmed == "" {
		return errors.New("token is empty")
	}

	if !isAuthEnabled {
		return errors.New("server authentication is not configured (ETHERPLY_JWT_SECRET missing)")
	}

	token, err := jwt.Parse(trimmed, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return fmt.Errorf("token parse failed: %w", err)
	}

	if !token.Valid {
		return errors.New("token signature invalid")
	}

	return nil
}
