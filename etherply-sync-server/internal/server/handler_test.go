package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bneb/etherply/etherply-sync-server/internal/crdt"
	"github.com/bneb/etherply/etherply-sync-server/internal/presence"
	"github.com/bneb/etherply/etherply-sync-server/internal/server"
	"github.com/bneb/etherply/etherply-sync-server/internal/store"
)

// createTestHandler creates a Handler with in-memory store for testing.
func createTestHandler() *server.Handler {
	memStore := store.NewMemoryStore()
	engine := crdt.NewEngine(memStore)
	presenceManager := presence.NewManager()
	return server.NewHandler(engine, presenceManager)
}

func TestHandleGetPresence_EmptyWorkspace(t *testing.T) {
	handler := createTestHandler()

	req := httptest.NewRequest("GET", "/v1/presence/new-workspace", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetPresence(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Should return empty array, not null
	var users []presence.User
	if err := json.Unmarshal(rr.Body.Bytes(), &users); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Empty workspace should return nil/empty slice (JSON: null or [])
	// The handler returns whatever GetUsers returns
	if users != nil && len(users) > 0 {
		t.Errorf("Expected empty users list, got %d users", len(users))
	}
}

func TestHandleGetPresence_InvalidPath(t *testing.T) {
	handler := createTestHandler()

	// Path too short - missing workspace_id
	req := httptest.NewRequest("GET", "/v1/presence", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetPresence(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid path, got %d", rr.Code)
	}
}

func TestHandleGetPresence_ContentType(t *testing.T) {
	handler := createTestHandler()

	req := httptest.NewRequest("GET", "/v1/presence/test-workspace", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetPresence(rr, req)

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestHandleWebSocket_InvalidPath(t *testing.T) {
	handler := createTestHandler()

	// Path too short - no workspace_id segment
	req := httptest.NewRequest("GET", "/v1/sync", nil)
	rr := httptest.NewRecorder()

	handler.HandleWebSocket(rr, req)

	// Should return Bad Request without attempting upgrade
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid path, got %d", rr.Code)
	}
}
