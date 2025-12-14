package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
)

func TestSessionAffinity_HeaderPresent(t *testing.T) {
	// Test that WebSocket connections receive a unique session ID in headers
	handler := createTestHandler()

	// Create test server
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.HandleWebSocket(w, r)
	}))
	defer s.Close()

	// Connect via WebSocket
	wsURL := "ws" + s.URL[4:] + "/v1/sync/ws-session-test"

	// Store the response header
	var responseHeader http.Header

	dialer := websocket.Dialer{}
	conn, resp, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	responseHeader = resp.Header

	// Verify X-Session-ID header is present
	sessionID := responseHeader.Get("X-Session-ID")
	if sessionID == "" {
		t.Error("Expected X-Session-ID header to be present")
	}

	// Verify it looks like a UUID (36 chars with dashes)
	if len(sessionID) != 36 {
		t.Errorf("Session ID should be 36 chars (UUID format), got %d chars: %s", len(sessionID), sessionID)
	}

	t.Logf("Session ID: %s", sessionID)
}

func TestSessionAffinity_UniqueSessions(t *testing.T) {
	// Test that different connections get different session IDs
	handler := createTestHandler()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.HandleWebSocket(w, r)
	}))
	defer s.Close()

	wsURL := "ws" + s.URL[4:] + "/v1/sync/ws-session-unique"

	// Connect multiple times and collect session IDs
	sessionIDs := make(map[string]bool)
	const numConnections = 10

	for i := 0; i < numConnections; i++ {
		conn, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("Connection %d failed: %v", i, err)
		}

		sessionID := resp.Header.Get("X-Session-ID")
		if sessionID == "" {
			t.Errorf("Connection %d: missing session ID", i)
			conn.Close()
			continue
		}

		if sessionIDs[sessionID] {
			t.Errorf("Duplicate session ID detected: %s", sessionID)
		}
		sessionIDs[sessionID] = true

		conn.Close()
	}

	if len(sessionIDs) != numConnections {
		t.Errorf("Expected %d unique session IDs, got %d", numConnections, len(sessionIDs))
	}

	t.Logf("Generated %d unique session IDs", len(sessionIDs))
}
