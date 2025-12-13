package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestValidateToken(t *testing.T) {
	// 1. Setup
	secret := "test-super-secret-key"
	Init(secret)

	// Helper to generate token
	sign := func(valid bool) string {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": "user-123",
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		
		key := []byte(secret)
		if !valid {
			key = []byte("wrong-key")
		}
		
		s, _ := token.SignedString(key)
		return s
	}

	validToken := sign(true)
	invalidSigToken := sign(false)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{"Valid Token", validToken, false},
		{"Invalid Signature", invalidSigToken, true},
		{"Garbage Token", "not.a.jwt", true},
		{"Empty Token", "", true},
		{"Whitespace Token", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateToken(tt.token); (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
