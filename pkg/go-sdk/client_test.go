package etherply_test

import (
	"testing"

	etherply "github.com/bneb/etherply/pkg/go-sdk"
)

func TestNewClient_Fields(t *testing.T) {
	client := etherply.NewClient("ws://localhost:8080", "test-token")

	if client.BaseURL != "ws://localhost:8080" {
		t.Errorf("BaseURL mismatch: got %s, want ws://localhost:8080", client.BaseURL)
	}
	if client.Token != "test-token" {
		t.Errorf("Token mismatch: got %s, want test-token", client.Token)
	}
	if client.Conn != nil {
		t.Error("Conn should be nil before Connect() is called")
	}
}

func TestClient_Close_NilConnection(t *testing.T) {
	// Close should not panic when connection is nil
	client := etherply.NewClient("ws://localhost:8080", "test-token")

	// This should not panic or error
	err := client.Close()
	if err != nil {
		t.Errorf("Close on nil connection should return nil, got %v", err)
	}
}

func TestClient_SendOperation_NotConnected(t *testing.T) {
	client := etherply.NewClient("ws://localhost:8080", "test-token")

	// SendOperation on nil connection should return nil (queued for offline support)
	// This verifies the "Offline Support" feature.
	err := client.SendOperation("key", "value")
	if err != nil {
		t.Errorf("Expected nil error (queued) when calling SendOperation without connection, got: %v", err)
	}
}

// contains checks if s contains substr (simple helper for tests)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
