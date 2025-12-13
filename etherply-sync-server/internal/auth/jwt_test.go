package auth

import "testing"

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{"Valid Token", "some-valid-token", false},
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
