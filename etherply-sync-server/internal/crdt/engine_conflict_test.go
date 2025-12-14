package crdt_test

import (
	"log"
	"testing"
	"time"

	"github.com/automerge/automerge-go"
	"github.com/bneb/etherply/etherply-sync-server/internal/crdt"
)

// TestOfflineEdit_TimestampRespect validates that an edit made in the "past"
// (simulating offline mode) is recorded with the correct timestamp in history.
func TestOfflineEdit_TimestampRespect(t *testing.T) {
	engine, ms := setupMockEngine()
	defer ms.Close()

	workspaceID := "ws-conflict"
	key := "offline_doc"

	// T0: Representing "Now" (Server time)
	now := time.Now()

	// T-1Hour: Representing an offline edit made 1 hour ago
	pastTime := now.Add(-1 * time.Hour)
	pastTs := pastTime.UnixMicro()

	op := crdt.Operation{
		WorkspaceID: workspaceID,
		Key:         key,
		Value:       "offline-content",
		Timestamp:   pastTs,
	}

	if err := engine.ProcessOperation(op); err != nil {
		t.Fatalf("Failed to process offline op: %v", err)
	}

	// Verify history directly via internal store inspection to check commit timestamp
	// This is a bit "white-box" but necessary to verify the requirement.
	// Normally we'd use a History() method on the engine if it exposed one.

	val, exists, _ := ms.Get(workspaceID, "automerge_root")
	if !exists {
		t.Fatal("Store empty")
	}
	data := val.([]byte)
	doc, _ := automerge.Load(data)

	// Get changes
	changes, err := doc.Changes()
	if err != nil {
		t.Fatalf("Failed to get changes: %v", err)
	}

	if len(changes) != 1 {
		t.Fatalf("Expected 1 change, got %d", len(changes))
	}

	commit := changes[0]
	// Automerge timestamp is in milliseconds? Or depends on library.
	// automerge-go Change.Timestamp() returns time.Time? No, it returns int64 milliseconds typically in JS,
	// let's check Go API.
	// Actually, automerge-go `Change` struct has `Timestamp() time.Time`.

	commitTime := commit.Timestamp()

	// Check if it matches our Past Time (within small delta for conversion)
	// Precision might be different (micro vs milli?)
	// Let's allow 1 second delta.

	diff := commitTime.Sub(pastTime)
	if diff < 0 {
		diff = -diff
	}
	if diff > time.Second {
		t.Errorf("Commit timestamp mismatch. Expected ~%v, got %v (diff %v)", pastTime, commitTime, diff)
	} else {
		log.Printf("Verified timestamp: %v matches expected %v", commitTime, pastTime)
	}
}

// TestConflict_ConcurrentWrites verifies that we don't crash on concurrent writes
// handled by the engine (serialized by lock).
// Ideally, we'd test merge, but Engine currently just Loads->Edits->Saves.
// Merge happens if we had two DIFFERENT docs and merged them.
// Here we are just testing that we can write multiple times.
func TestSequentialWrites_History(t *testing.T) {
	engine, ms := setupMockEngine()
	defer ms.Close()

	workspaceID := "ws-seq"
	key := "seq-key"

	// Write 1
	_ = engine.ProcessOperation(crdt.Operation{
		WorkspaceID: workspaceID,
		Key:         key,
		Value:       "v1",
		Timestamp:   time.Now().UnixMicro(),
	})

	// Write 2
	_ = engine.ProcessOperation(crdt.Operation{
		WorkspaceID: workspaceID,
		Key:         key,
		Value:       "v2",
		Timestamp:   time.Now().Add(time.Second).UnixMicro(),
	})

	val, _, _ := ms.Get(workspaceID, "automerge_root")
	doc, _ := automerge.Load(val.([]byte))
	// We should see history growing
	changes, _ := doc.Changes()
	if len(changes) != 2 {
		t.Errorf("Expected 2 changes in history, got %d", len(changes))
	}
}
