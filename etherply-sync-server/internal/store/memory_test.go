package store_test

import (
	"sync"
	"testing"

	"github.com/bneb/etherply/etherply-sync-server/internal/store"
)

func TestMemoryStore_SetGet(t *testing.T) {
	ms := store.NewMemoryStore()
	defer ms.Close()

	// Test 1: Set and Get
	err := ms.Set("workspace-1", "key1", "value1")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, exists, err := ms.Get("workspace-1", "key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !exists {
		t.Error("Expected value to exist")
	}
	if val != "value1" {
		t.Errorf("Expected 'value1', got %v", val)
	}
}

func TestMemoryStore_GetNonExistent(t *testing.T) {
	ms := store.NewMemoryStore()
	defer ms.Close()

	// Test: Get non-existent workspace
	_, exists, err := ms.Get("non-existent", "key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if exists {
		t.Error("Expected value to not exist for non-existent workspace")
	}

	// Test: Get non-existent key in existing workspace
	ms.Set("workspace-1", "key1", "value1")
	_, exists, err = ms.Get("workspace-1", "non-existent-key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if exists {
		t.Error("Expected value to not exist for non-existent key")
	}
}

func TestMemoryStore_GetAll(t *testing.T) {
	ms := store.NewMemoryStore()
	defer ms.Close()

	// Test 1: GetAll on empty workspace returns empty map (not nil)
	result, err := ms.GetAll("empty-workspace")
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if result == nil {
		t.Error("Expected empty map, got nil")
	}
	if len(result) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(result))
	}

	// Test 2: GetAll with data
	ms.Set("workspace-1", "k1", "v1")
	ms.Set("workspace-1", "k2", "v2")

	result, err = ms.GetAll("workspace-1")
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(result))
	}
}

func TestMemoryStore_GetAll_ReturnsCopy(t *testing.T) {
	ms := store.NewMemoryStore()
	defer ms.Close()

	ms.Set("ws", "key", "original")

	// Get the map and modify it
	result, _ := ms.GetAll("ws")
	result["key"] = "modified"

	// Original should be unchanged
	val, _, _ := ms.Get("ws", "key")
	if val != "original" {
		t.Error("GetAll should return a copy, not a reference to internal data")
	}
}

func TestMemoryStore_Overwrite(t *testing.T) {
	ms := store.NewMemoryStore()
	defer ms.Close()

	ms.Set("ws", "key", "first")
	ms.Set("ws", "key", "second")

	val, _, _ := ms.Get("ws", "key")
	if val != "second" {
		t.Errorf("Expected 'second' after overwrite, got %v", val)
	}
}

func TestMemoryStore_MultipleWorkspaces(t *testing.T) {
	ms := store.NewMemoryStore()
	defer ms.Close()

	ms.Set("ws-a", "key", "value-a")
	ms.Set("ws-b", "key", "value-b")

	valA, _, _ := ms.Get("ws-a", "key")
	valB, _, _ := ms.Get("ws-b", "key")

	if valA != "value-a" {
		t.Errorf("Expected 'value-a', got %v", valA)
	}
	if valB != "value-b" {
		t.Errorf("Expected 'value-b', got %v", valB)
	}
}

func TestMemoryStore_Concurrency(t *testing.T) {
	ms := store.NewMemoryStore()
	defer ms.Close()

	// Test concurrent access doesn't cause race conditions
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := "key"
			ms.Set("ws", key, n)
			ms.Get("ws", key)
		}(i)
	}
	wg.Wait()

	// Just verify we don't panic - exact value doesn't matter due to races
	_, exists, _ := ms.Get("ws", "key")
	if !exists {
		t.Error("Expected key to exist after concurrent writes")
	}
}

func TestMemoryStore_Close(t *testing.T) {
	ms := store.NewMemoryStore()

	// Close should return nil for MemoryStore (no resources to release)
	err := ms.Close()
	if err != nil {
		t.Errorf("Close should return nil, got %v", err)
	}
}
