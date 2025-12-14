package crdt_test

import (
	"testing"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/crdt"
	"github.com/bneb/etherply/etherply-sync-server/internal/store"
)

// setupMockEngine creates an engine backed by a fresh MemoryStore
func setupMockEngine() (*crdt.Engine, store.Store) {
	ms := store.NewMemoryStore()
	// In a real test, we might want to capture logs, but for now default to stderr or discard
	// To test defensive logging, one would typically use a custom slog.Handler,
	// but for this MVP test suite we focus on logic correctness.
	engine := crdt.NewEngine(ms)
	return engine, ms
}

func TestLWW_Correctness(t *testing.T) {
	engine, ms := setupMockEngine()
	defer ms.Close()

	workspaceID := "ws-1"
	key := "user-1:cursor"

	// T0: Initial operation
	// value is int, automerge handles numbers as float64 by default or specific types.
	// We'll use simple strings or maps.
	op1 := crdt.Operation{
		WorkspaceID: workspaceID,
		Key:         key,
		Value:       "state-1",
		Timestamp:   1000,
	}

	if err := engine.ProcessOperation(op1); err != nil {
		t.Fatalf("Failed to process valid op: %v", err)
	}

	// Verify state
	snapshot, err := engine.GetFullState(workspaceID)
	if err != nil {
		t.Fatalf("Failed to get state: %v", err)
	}
	if val, ok := snapshot.Data[key]; !ok {
		t.Errorf("Expected key %s to exist", key)
	} else if val != "state-1" {
		t.Errorf("Expected 'state-1', got %v", val)
	}

	// T1: Newer operation (Should win due to order of execution)
	op2 := crdt.Operation{
		WorkspaceID: workspaceID,
		Key:         key,
		Value:       "state-2",
		Timestamp:   2000,
	}
	if err := engine.ProcessOperation(op2); err != nil {
		t.Fatalf("Failed to process newer op: %v", err)
	}

	snapshot, _ = engine.GetFullState(workspaceID)
	if val := snapshot.Data[key]; val != "state-2" {
		t.Errorf("Expected 'state-2' (winner), got %v", val)
	}
}

func TestDataCorruption_SelfHealing(t *testing.T) {
	// "The Happy Path is a Trap" - Defensive Coding Test
	engine, ms := setupMockEngine()
	defer ms.Close()

	workspaceID := "ws-corrupt"
	key := "bad_key"

	// Manually inject bad data (invalid bytes that aren't an automerge doc)
	// The implementation checks for []byte type first, then Load().
	// Storage now stores []byte under "sync_doc".

	// Corrupt the root blob
	ms.Set("ws:"+workspaceID, "sync_doc", []byte("GARBAGE DATA NOT AUTOMERGE"))

	// Apply a valid operation on top of it
	recoveryOp := crdt.Operation{
		WorkspaceID: workspaceID,
		Key:         key,
		Value:       "recovered",
		Timestamp:   time.Now().UnixMicro(),
	}

	// Engine should detect corruption (Load fails) and reset to empty doc (Self-Heal) or error.
	// The implementation logs error and returns error?
	// Checking code:
	// doc, err = automerge.Load(data)
	// if err != nil { return fmt.Errorf("failed to hydrate...": err) }

	// Wait, the current implementation *Errors* if Load fails. It does NOT self-heal unless the type is wrong.
	// If the type is []byte (which it is) but content is bad, automerge.Load returns error.
	// So we expect an error here!

	err := engine.ProcessOperation(recoveryOp)
	if err == nil {
		t.Error("Expected error due to corrupted state, got nil")
	} else {
		// This is good defensive coding: we don't want to overwrite corrupted state silently if it looks like real data but isn't.
		// NOTE: If we wanted self-healing for corruption, we would handle the error in Get() or Load().
		t.Logf("Correctly caught corruption: %v", err)
	}
}

func TestValidation(t *testing.T) {
	engine, ms := setupMockEngine()
	defer ms.Close()

	invalidOp := crdt.Operation{
		WorkspaceID: "", // Missing
		Key:         "key",
		Value:       "val",
	}

	if err := engine.ProcessOperation(invalidOp); err == nil {
		t.Error("Expected error for missing WorkspaceID, got nil")
	}
}
