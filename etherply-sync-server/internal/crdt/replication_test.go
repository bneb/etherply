// Package crdt contains tests for the ApplyRemoteChanges functionality.
// These tests validate the CRDT engine's ability to merge concurrent changes
// from different regions using Automerge's merge semantics.
package crdt

import (
	"testing"

	"github.com/bneb/etherply/etherply-sync-server/internal/store"
)

// TestApplyRemoteChanges_EmptyWorkspaceID verifies validation.
func TestApplyRemoteChanges_EmptyWorkspaceID(t *testing.T) {
	memStore := store.NewMemoryStore()
	engine := NewEngine(memStore)

	err := engine.ApplyRemoteChanges("", []byte("data"))
	if err == nil {
		t.Error("expected error for empty workspace_id")
	}
}

// TestApplyRemoteChanges_EmptyChanges verifies no-op behavior.
func TestApplyRemoteChanges_EmptyChanges(t *testing.T) {
	memStore := store.NewMemoryStore()
	engine := NewEngine(memStore)

	err := engine.ApplyRemoteChanges("test-ws", nil)
	if err != nil {
		t.Errorf("unexpected error for empty changes: %v", err)
	}
}

// TestApplyRemoteChanges_Integration verifies merging works.
func TestApplyRemoteChanges_Integration(t *testing.T) {
	memStore1 := store.NewMemoryStore()
	engine1 := NewEngine(memStore1)

	memStore2 := store.NewMemoryStore()
	engine2 := NewEngine(memStore2)

	workspaceID := "repl-test-ws"

	// Engine 1 creates a document
	err := engine1.ProcessOperation(Operation{
		WorkspaceID: workspaceID,
		Key:         "title",
		Value:       "Hello from Region 1",
	})
	if err != nil {
		t.Fatalf("engine1 failed to process operation: %v", err)
	}

	// Get engine1's state as "remote" changes
	state1, _, err := memStore1.Get(workspaceID, "automerge_root")
	if err != nil {
		t.Fatalf("failed to get engine1 state: %v", err)
	}

	// Engine 2 applies the remote changes
	err = engine2.ApplyRemoteChanges(workspaceID, state1.([]byte))
	if err != nil {
		t.Fatalf("engine2 failed to apply remote changes: %v", err)
	}

	// Verify engine2 has the same data
	snapshot2, err := engine2.GetFullState(workspaceID)
	if err != nil {
		t.Fatalf("failed to get engine2 state: %v", err)
	}

	if snapshot2.Data["title"] != "Hello from Region 1" {
		t.Errorf("expected 'Hello from Region 1', got '%v'", snapshot2.Data["title"])
	}
}

// TestApplyRemoteChanges_BidirectionalMerge tests convergence.
func TestApplyRemoteChanges_BidirectionalMerge(t *testing.T) {
	memStore1 := store.NewMemoryStore()
	engine1 := NewEngine(memStore1)

	memStore2 := store.NewMemoryStore()
	engine2 := NewEngine(memStore2)

	workspaceID := "bidir-merge-ws"

	// Engine 1 and 2 create concurrent changes
	err := engine1.ProcessOperation(Operation{
		WorkspaceID: workspaceID,
		Key:         "field_a",
		Value:       "Value A from Engine 1",
	})
	if err != nil {
		t.Fatalf("engine1 op failed: %v", err)
	}

	err = engine2.ProcessOperation(Operation{
		WorkspaceID: workspaceID,
		Key:         "field_b",
		Value:       "Value B from Engine 2",
	})
	if err != nil {
		t.Fatalf("engine2 op failed: %v", err)
	}

	// Cross-merge
	state1, _, _ := memStore1.Get(workspaceID, "automerge_root")
	state2, _, _ := memStore2.Get(workspaceID, "automerge_root")

	if err := engine1.ApplyRemoteChanges(workspaceID, state2.([]byte)); err != nil {
		t.Fatalf("engine1 merge failed: %v", err)
	}
	if err := engine2.ApplyRemoteChanges(workspaceID, state1.([]byte)); err != nil {
		t.Fatalf("engine2 merge failed: %v", err)
	}

	// Both should now have both fields
	snap1, _ := engine1.GetFullState(workspaceID)
	snap2, _ := engine2.GetFullState(workspaceID)

	if snap1.Data["field_a"] != "Value A from Engine 1" {
		t.Errorf("engine1 missing field_a")
	}
	if snap1.Data["field_b"] != "Value B from Engine 2" {
		t.Errorf("engine1 missing field_b")
	}
	if snap2.Data["field_a"] != "Value A from Engine 1" {
		t.Errorf("engine2 missing field_a")
	}
	if snap2.Data["field_b"] != "Value B from Engine 2" {
		t.Errorf("engine2 missing field_b")
	}
}
