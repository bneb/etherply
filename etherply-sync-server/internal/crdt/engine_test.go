package crdt_test

import (
	"testing"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/crdt"
	"github.com/bneb/etherply/etherply-sync-server/internal/store"
)

// setupMockEngine creates an engine backed by a fresh MemoryStore
func setupMockEngine() (*crdt.Engine, store.StateStore) {
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
	op1 := crdt.Operation{
		WorkspaceID: workspaceID,
		Key:         key,
		Value:       map[string]int{"x": 10, "y": 10},
		Timestamp:   1000,
	}

	if err := engine.ProcessOperation(op1); err != nil {
		t.Fatalf("Failed to process valid op: %v", err)
	}

	// Verify state
	val, _ := ms.Get(workspaceID, key)
	storedOp := val.(crdt.Operation)
	if storedOp.Timestamp != 1000 {
		t.Errorf("Expected timestamp 1000, got %d", storedOp.Timestamp)
	}

	// T1: Newer operation (Should win)
	op2 := crdt.Operation{
		WorkspaceID: workspaceID,
		Key:         key,
		Value:       map[string]int{"x": 20, "y": 20},
		Timestamp:   2000,
	}
	if err := engine.ProcessOperation(op2); err != nil {
		t.Fatalf("Failed to process newer op: %v", err)
	}

	val, _ = ms.Get(workspaceID, key)
	storedOp = val.(crdt.Operation)
	if storedOp.Timestamp != 2000 {
		t.Errorf("Expected timestamp 2000 (winner), got %d", storedOp.Timestamp)
	}
}

func TestLWW_StaleOperation(t *testing.T) {
	engine, ms := setupMockEngine()
	defer ms.Close()

	workspaceID := "ws-1"
	key := "config"

	// Set initial state at T=5000
	initialOp := crdt.Operation{
		WorkspaceID: workspaceID,
		Key:         key,
		Value:       "latest",
		Timestamp:   5000,
	}
	engine.ProcessOperation(initialOp)

	// Attempt to process older op at T=4000 (Should be ignored)
	staleOp := crdt.Operation{
		WorkspaceID: workspaceID,
		Key:         key,
		Value:       "stale",
		Timestamp:   4000,
	}
	if err := engine.ProcessOperation(staleOp); err != nil {
		t.Fatalf("Failed to process stale op: %v", err)
	}

	// Verify state is UNCHANGED
	val, _ := ms.Get(workspaceID, key)
	storedOp := val.(crdt.Operation)
	if storedOp.Value != "latest" {
		t.Errorf("LWW Violation: Expected 'latest', got '%v'", storedOp.Value)
	}
}

func TestClockSkew_FutureTimestamp(t *testing.T) {
	engine, ms := setupMockEngine()
	defer ms.Close()

	// Operation from 10 minutes in the future
	futureTime := time.Now().Add(10 * time.Minute).UnixMicro()
	
	op := crdt.Operation{
		WorkspaceID: "ws-skew",
		Key:         "key",
		Value:       "future_value",
		Timestamp:   futureTime,
	}

	// Should not crash, should be accepted (policy is to warn, not reject currently)
	if err := engine.ProcessOperation(op); err != nil {
		t.Fatalf("Engine rejected future timestamp: %v", err)
	}

	val, exists := ms.Get("ws-skew", "key")
	if !exists {
		t.Error("Expected future op to be persisted")
	}
	storedOp := val.(crdt.Operation)
	if storedOp.Timestamp != futureTime {
		t.Error("Stored timestamp mismatch")
	}
}

func TestDataCorruption_SelfHealing(t *testing.T) {
	// "The Happy Path is a Trap" - Defensive Coding Test
	engine, ms := setupMockEngine()
	defer ms.Close()

	workspaceID := "ws-corrupt"
	key := "bad_key"

	// Manually inject bad data (not an Operation struct) directly into store
	// This simulates schema drift or bit rot
	ms.Set(workspaceID, key, "I AM NOT AN OPERATION")

	// Apply a valid operation on top of it
	recoveryOp := crdt.Operation{
		WorkspaceID: workspaceID,
		Key:         key,
		Value:       "recovered",
		Timestamp:   time.Now().UnixMicro(),
	}

	// Engine should detect type assertion failure and OVERWRITE (Self-Heal)
	if err := engine.ProcessOperation(recoveryOp); err != nil {
		t.Fatalf("Engine failed to recover from corruption: %v", err)
	}

	val, _ := ms.Get(workspaceID, key)
	// Should now be a valid Operation
	if storedOp, ok := val.(crdt.Operation); !ok {
		t.Error("Failed to self-heal: value is still not an Operation")
	} else {
		if storedOp.Value != "recovered" {
			t.Errorf("Expected 'recovered', got %v", storedOp.Value)
		}
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
