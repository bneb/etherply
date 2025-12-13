package auth

import (
	"errors"
	"strings"
)

// ValidateToken checks if the provided token string is valid.
// For this stage of the "Commercial Due Diligence" remediation, we are
// strictly enforcing that a token exists and is non-empty.
//
// TODO(security): Replace this with real JWT signature verification (e.g. usage of github.com/golang-jwt/jwt/v5)
// once an Identity Provider (Auth0, Supabase) is selected.
func ValidateToken(tokenString string) error {
	trimmed := strings.TrimSpace(tokenString)
	if trimmed == "" {
		return errors.New("token is empty")
	}

	// Future: Parse and verify JWT signature here.
	// if !token.Valid { return errors.New("invalid signature") }

	// For now, we assume if the client sends a Bearer token, they are "authorized"
	// effectively moving from "Open Public" to "Shared Secret" model if we were to check a specific value,
	// but here we just ensure the protocol is respected (Authentication Facade Removal).
	return nil
}
