package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/crdt"
)

func TestHandleGetHistory(t *testing.T) {
	// 1. Setup
	crdtEngine, _, handler := createTestHandlerWithComponents()

	workspaceID := "history-ws"

	// 2. Generate History (3 commits)
	ops := []crdt.Operation{
		{WorkspaceID: workspaceID, Key: "step1", Value: "start", Timestamp: time.Now().UnixMicro()},
		{WorkspaceID: workspaceID, Key: "step2", Value: "middle", Timestamp: time.Now().UnixMicro()},
		{WorkspaceID: workspaceID, Key: "step3", Value: "end", Timestamp: time.Now().UnixMicro()},
	}

	for _, op := range ops {
		if err := crdtEngine.ProcessOperation(op); err != nil {
			t.Fatalf("Failed to process op: %v", err)
		}
	}

	// 3. Call API
	req := httptest.NewRequest("GET", "/v1/history/"+workspaceID, nil)
	rr := httptest.NewRecorder()

	handler.HandleGetHistory(rr, req)

	// 4. Verify
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var history []map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &history); err != nil {
		t.Fatal(err)
	}

	// Expect 3 entries
	if len(history) != 3 {
		t.Errorf("Expected 3 history entries, got %d", len(history))
	}

	// Verify order (Automerge usually returns chronological or reverse? doc.Changes() returns chronological)
	// Let's check the messages.
	// NOTE: ProcessOperation sets commit message to "set {Key}"
	expectedMsgs := []string{"set step1", "set step2", "set step3"}

	for i, entry := range history {
		msg := entry["message"].(string)
		if msg != expectedMsgs[i] {
			t.Errorf("Index %d: Expected message %q, got %q", i, expectedMsgs[i], msg)
		}
		if entry["hash"] == "" {
			t.Error("Missing hash")
		}
	}
}
