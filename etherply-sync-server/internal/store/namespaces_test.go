package store_test

import (
	"testing"

	"github.com/bneb/etherply/etherply-sync-server/internal/store"
)

// TestNamespaceIsolation ensures that identical keys in different namespaces do not collide.
func TestNamespaceIsolation(t *testing.T) {
	s := store.NewMemoryStore()

	key := "config"
	val1 := "user-data"
	val2 := "system-data"

	// Write to Namespace A
	if err := s.Set("ws:customer-1", key, val1); err != nil {
		t.Fatalf("Failed to write to ns A: %v", err)
	}

	// Write same key to Namespace B
	if err := s.Set("sys:projects", key, val2); err != nil {
		t.Fatalf("Failed to write to ns B: %v", err)
	}

	// Read from Namespace A
	readVal1, exists, err := s.Get("ws:customer-1", key)
	if err != nil || !exists {
		t.Fatalf("Failed to read from ns A: %v", err)
	}
	if readVal1 != val1 {
		t.Errorf("Namespace A contamination! Got %v, want %v", readVal1, val1)
	}

	// Read from Namespace B
	readVal2, exists, err := s.Get("sys:projects", key)
	if err != nil || !exists {
		t.Fatalf("Failed to read from ns B: %v", err)
	}
	if readVal2 != val2 {
		t.Errorf("Namespace B contamination! Got %v, want %v", readVal2, val2)
	}
}
