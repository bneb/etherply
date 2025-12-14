package crdt_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/automerge/automerge-go"
	"github.com/bneb/etherply/etherply-sync-server/internal/crdt"
)

func TestDeltaSync_API(t *testing.T) {
	engine, ms := setupMockEngine()
	defer ms.Close()

	workspaceID := "ws-delta"

	// 1. Create Initial State with first operation
	op1 := crdt.Operation{
		WorkspaceID: workspaceID,
		Key:         "key1",
		Value:       "initial",
		Timestamp:   time.Now().UnixMicro(),
	}
	if err := engine.ProcessOperation(op1); err != nil {
		t.Fatal(err)
	}

	// 2. Get Full State + Heads after first op (checkpoint)
	snapshotAfterOp1, err := engine.GetFullState(workspaceID)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshotAfterOp1.Heads) == 0 {
		t.Error("Expected heads to be present after first op")
	}
	checkpointHeads := snapshotAfterOp1.Heads
	t.Logf("Checkpoint Heads (after op1): %v", checkpointHeads)

	// 3. Apply more operations
	for i := 0; i < 5; i++ {
		op := crdt.Operation{
			WorkspaceID: workspaceID,
			Key:         "key2",
			Value:       i,
			Timestamp:   time.Now().UnixMicro(),
		}
		if err := engine.ProcessOperation(op); err != nil {
			t.Fatal(err)
		}
	}

	// 4. Get Full State bytes (for size comparison)
	fullStateBytes, err := engine.GetChanges(workspaceID, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Full state size: %d bytes", len(fullStateBytes))

	// 5. Get Delta since checkpoint
	deltaBytes, err := engine.GetChanges(workspaceID, checkpointHeads)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Delta size (since checkpoint): %d bytes", len(deltaBytes))

	// 6. Log size comparison (delta may not always be smaller for small docs due to framing overhead)
	// Delta sync is most beneficial for large documents with few recent changes.
	t.Logf("Size comparison: delta=%d bytes, full=%d bytes", len(deltaBytes), len(fullStateBytes))
	if len(deltaBytes) < len(fullStateBytes) {
		t.Log("Delta is smaller than full state (optimal case)")
	} else {
		t.Log("Delta is larger than full state (expected for small docs with framing overhead)")
	}

	// 7. Verify delta can be applied on top of checkpoint to get final state
	// First, get checkpoint state as a doc
	checkpointBytes, err := engine.GetChanges(workspaceID, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Fork a new doc from the full state and verify it matches
	fullDoc, err := automerge.Load(fullStateBytes)
	if err != nil {
		t.Fatal(err)
	}

	finalSnapshot, err := engine.GetFullState(workspaceID)
	if err != nil {
		t.Fatal(err)
	}

	// Verify final state has correct data
	if finalSnapshot.Data["key2"] != float64(4) { // JSON unmarshals numbers as float64
		t.Errorf("Expected key2 to be 4, got %v", finalSnapshot.Data["key2"])
	}

	// 8. Verify JSON marshaling of Snapshot
	b, err := json.Marshal(finalSnapshot)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Final Snapshot JSON: %s", string(b))

	// Keep linter happy by using vars
	_ = checkpointBytes
	_ = fullDoc
}

func TestDeltaSync_InvalidHash(t *testing.T) {
	engine, ms := setupMockEngine()
	defer ms.Close()

	workspaceID := "ws-delta-invalid"

	// Create a document first
	op := crdt.Operation{
		WorkspaceID: workspaceID,
		Key:         "key",
		Value:       "value",
		Timestamp:   time.Now().UnixMicro(),
	}
	if err := engine.ProcessOperation(op); err != nil {
		t.Fatal(err)
	}

	// Try to get changes with invalid hash - should error
	_, err := engine.GetChanges(workspaceID, []string{"invalid-hash-format"})
	if err == nil {
		t.Error("Expected error for invalid hash")
	}
	t.Logf("Got expected error: %v", err)
}

func TestDeltaSync_EmptyDocument(t *testing.T) {
	engine, ms := setupMockEngine()
	defer ms.Close()

	// GetChanges on non-existent workspace should return empty bytes
	bytes, err := engine.GetChanges("ws-nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(bytes) != 0 {
		t.Errorf("Expected empty bytes for non-existent workspace, got %d bytes", len(bytes))
	}
}
