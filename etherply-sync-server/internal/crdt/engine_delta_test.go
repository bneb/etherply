package crdt_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/crdt"
)

func TestDeltaSync_API(t *testing.T) {
	engine, ms := setupMockEngine()
	defer ms.Close()

	workspaceID := "ws-delta"
	key := "delta-key"

	// 1. Create Initial State
	op1 := crdt.Operation{
		WorkspaceID: workspaceID,
		Key:         key,
		Value:       "v1",
		Timestamp:   time.Now().UnixMicro(),
	}
	if err := engine.ProcessOperation(op1); err != nil {
		t.Fatal(err)
	}

	// 2. Get Full State + Heads
	snapshot, err := engine.GetFullState(workspaceID)
	if err != nil {
		t.Fatal(err)
	}

	if len(snapshot.Heads) == 0 {
		t.Error("Expected heads to be present")
	}
	t.Logf("Current Heads: %v", snapshot.Heads)

	// 3. Get Changes (using empty since)
	changesBytes, err := engine.GetChanges(workspaceID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(changesBytes) == 0 {
		t.Error("Expected changes bytes")
	}
	// For MVP Stub, this is the Doc.Save() bytes
	// We can try to decode it as doc or changes?
	// Currently stub returns doc.Save().

	// 4. Get Changes (using current heads - should be empty or delta)
	// Current stub ignores 'since', so it returns full doc.
	// We just verify API runs without error.
	changesSince, err := engine.GetChanges(workspaceID, snapshot.Heads)
	if err != nil {
		t.Fatal(err)
	}
	if len(changesSince) == 0 {
		t.Error("Expected content from stub")
	}

	// 5. Verify JSON marshaling of Snapshot
	b, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Snapshot JSON: %s", string(b))
}
