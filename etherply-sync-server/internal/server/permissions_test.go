package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/auth"
	"github.com/bneb/etherply/etherply-sync-server/internal/crdt"
	"github.com/gorilla/websocket"
)

func TestACL_Enforcement(t *testing.T) {
	// 1. Setup Test Server
	handler := createTestHandler()

	// Create a test server that wraps our handler logic
	// But we need to inject scopes into the context.
	// Since Middleware runs before Handler, normally Middleware does this.
	// For testing HandleWebSocket directly, we can wrap the request with context.

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock Middleware behavior: Inject Scopes from Query Param for easy testing
		scopeParam := r.URL.Query().Get("mock_scopes")
		var scopes []string
		if scopeParam != "" {
			scopes = []string{scopeParam}
		}

		ctx := auth.NewContextWithScopes(r.Context(), scopes)
		handler.HandleWebSocket(w, r.WithContext(ctx))
	}))
	defer s.Close()

	// 2. Test Case: Read-Only Token (scope="read")
	// Should connect, receive init, but FAIL to write.
	wsURL := "ws" + s.URL[4:] + "/v1/sync/ws-acl?mock_scopes=read"

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Read Init
	var initMsg map[string]interface{}
	if err := conn.ReadJSON(&initMsg); err != nil {
		t.Fatalf("Failed to read init: %v", err)
	}
	if initMsg["type"] != "init" {
		t.Errorf("Expected init message, got %v", initMsg["type"])
	}

	// Try to Write
	opMsg := map[string]interface{}{
		"type": "op",
		"payload": crdt.Operation{
			Key:       "illegal",
			Value:     "write",
			Timestamp: time.Now().UnixMicro(),
		},
	}
	if err := conn.WriteJSON(opMsg); err != nil {
		t.Fatal(err)
	}

	// Expect Error Message back
	var errMsg map[string]interface{}
	if err := conn.ReadJSON(&errMsg); err != nil {
		t.Fatalf("Failed to read back error: %v", err)
	}

	if errMsg["type"] != "error" {
		t.Errorf("Expected error message, got %v", errMsg["type"])
	}

	// 3. Test Case: Write Token (scope="write")
	// Should succeed.
	wsURL2 := "ws" + s.URL[4:] + "/v1/sync/ws-acl-write?mock_scopes=write"
	conn2, _, err := websocket.DefaultDialer.Dial(wsURL2, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn2.Close()

	// Read Init
	conn2.ReadJSON(&initMsg)

	// Write
	if err := conn2.WriteJSON(opMsg); err != nil {
		t.Fatal(err)
	}

	// Expect NO error. We can't easily check success without reading echoing or checking state.
	// But if we don't get an error immediately, it's a good sign.
	// Let's set a read deadline.
	conn2.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	var response map[string]interface{}
	// We expect the ECHO back (broadcast)
	err = conn2.ReadJSON(&response)
	if err != nil {
		t.Fatalf("Expected echo (success), got error: %v", err)
	}
	if response["type"] == "error" {
		t.Errorf("Expected success, got error response: %v", response)
	}
}
