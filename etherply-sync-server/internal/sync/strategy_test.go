package sync_test

import (
	"testing"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/sync"
)

func TestAutomergeStrategy_BasicOperations(t *testing.T) {
	s := sync.NewAutomergeStrategy()

	if s.Name() != "automerge" {
		t.Errorf("expected name 'automerge', got '%s'", s.Name())
	}

	// Test write
	doc, err := s.ProcessWrite(nil, "key1", "value1", time.Now())
	if err != nil {
		t.Fatalf("ProcessWrite failed: %v", err)
	}

	// Test state retrieval
	state, err := s.GetState(doc)
	if err != nil {
		t.Fatalf("GetState failed: %v", err)
	}

	if state["key1"] != "value1" {
		t.Errorf("expected key1='value1', got '%v'", state["key1"])
	}

	// Test heads
	heads, err := s.GetHeads(doc)
	if err != nil {
		t.Fatalf("GetHeads failed: %v", err)
	}
	if len(heads) == 0 {
		t.Error("expected at least one head")
	}
}

func TestAutomergeStrategy_Merge(t *testing.T) {
	s := sync.NewAutomergeStrategy()

	// Create two divergent documents
	doc1, _ := s.ProcessWrite(nil, "field_a", "from_doc1", time.Now())
	doc2, _ := s.ProcessWrite(nil, "field_b", "from_doc2", time.Now())

	// Merge
	merged, err := s.Merge(doc1, doc2)
	if err != nil {
		t.Fatalf("Merge failed: %v", err)
	}

	state, _ := s.GetState(merged)
	if state["field_a"] != "from_doc1" {
		t.Errorf("merged doc missing field_a")
	}
	if state["field_b"] != "from_doc2" {
		t.Errorf("merged doc missing field_b")
	}
}

func TestLWWStrategy_BasicOperations(t *testing.T) {
	s := sync.NewLWWStrategy()

	if s.Name() != "lww" {
		t.Errorf("expected name 'lww', got '%s'", s.Name())
	}

	// Test write
	now := time.Now()
	doc, err := s.ProcessWrite(nil, "key1", "value1", now)
	if err != nil {
		t.Fatalf("ProcessWrite failed: %v", err)
	}

	// Test state
	state, err := s.GetState(doc)
	if err != nil {
		t.Fatalf("GetState failed: %v", err)
	}

	if state["key1"] != "value1" {
		t.Errorf("expected key1='value1', got '%v'", state["key1"])
	}
}

func TestLWWStrategy_TimestampWins(t *testing.T) {
	s := sync.NewLWWStrategy()

	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

	// Write with "now" timestamp
	doc, _ := s.ProcessWrite(nil, "key", "value_now", now)

	// Try to write with past timestamp - should NOT overwrite
	doc, _ = s.ProcessWrite(doc, "key", "value_past", past)

	state, _ := s.GetState(doc)
	if state["key"] != "value_now" {
		t.Errorf("LWW should keep 'value_now', got '%v'", state["key"])
	}

	// Write with future timestamp - should overwrite
	doc, _ = s.ProcessWrite(doc, "key", "value_future", future)

	state, _ = s.GetState(doc)
	if state["key"] != "value_future" {
		t.Errorf("LWW should update to 'value_future', got '%v'", state["key"])
	}
}

func TestLWWStrategy_Merge(t *testing.T) {
	s := sync.NewLWWStrategy()

	now := time.Now()
	past := now.Add(-1 * time.Hour)

	// Create doc1 with older timestamp
	doc1, _ := s.ProcessWrite(nil, "conflict_key", "older_value", past)

	// Create doc2 with newer timestamp
	doc2, _ := s.ProcessWrite(nil, "conflict_key", "newer_value", now)

	// Merge - newer should win
	merged, err := s.Merge(doc1, doc2)
	if err != nil {
		t.Fatalf("Merge failed: %v", err)
	}

	state, _ := s.GetState(merged)
	if state["conflict_key"] != "newer_value" {
		t.Errorf("LWW merge should pick newer value, got '%v'", state["conflict_key"])
	}
}

func TestLWWStrategy_NoHistory(t *testing.T) {
	s := sync.NewLWWStrategy()

	doc, _ := s.ProcessWrite(nil, "key", "value", time.Now())

	history, err := s.GetHistory(doc)
	if err != nil {
		t.Fatalf("GetHistory failed: %v", err)
	}

	if len(history) != 0 {
		t.Errorf("LWW should not track history, got %d entries", len(history))
	}
}

func TestServerAuthStrategy_BasicOperations(t *testing.T) {
	s := sync.NewServerAuthStrategy()

	if s.Name() != "server-auth" {
		t.Errorf("expected name 'server-auth', got '%s'", s.Name())
	}

	// Test write
	doc, err := s.ProcessWrite(nil, "key1", "value1", time.Now())
	if err != nil {
		t.Fatalf("ProcessWrite failed: %v", err)
	}

	state, err := s.GetState(doc)
	if err != nil {
		t.Fatalf("GetState failed: %v", err)
	}

	if state["key1"] != "value1" {
		t.Errorf("expected key1='value1', got '%v'", state["key1"])
	}
}

func TestServerAuthStrategy_LocalWins(t *testing.T) {
	s := sync.NewServerAuthStrategy()

	// Create local state
	local, _ := s.ProcessWrite(nil, "key", "local_value", time.Now())

	// Create remote state
	remote, _ := s.ProcessWrite(nil, "key", "remote_value", time.Now())

	// Merge - local (server) should win
	merged, err := s.Merge(local, remote)
	if err != nil {
		t.Fatalf("Merge failed: %v", err)
	}

	state, _ := s.GetState(merged)
	if state["key"] != "local_value" {
		t.Errorf("ServerAuth should keep local value, got '%v'", state["key"])
	}
}

func TestNewStrategy_Factory(t *testing.T) {
	tests := []struct {
		strategyType sync.StrategyType
		expectedName string
	}{
		{sync.StrategyAutomerge, "automerge"},
		{sync.StrategyLWW, "lww"},
		{sync.StrategyServerAuthoritative, "server-auth"},
		{sync.StrategyType("unknown"), "automerge"}, // Default fallback
	}

	for _, tc := range tests {
		s := sync.NewStrategy(tc.strategyType)
		if s.Name() != tc.expectedName {
			t.Errorf("NewStrategy(%s) = %s, want %s", tc.strategyType, s.Name(), tc.expectedName)
		}
	}
}
