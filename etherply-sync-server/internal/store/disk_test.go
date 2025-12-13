package store_test

import (
	"encoding/gob"
	"os"
	"testing"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/store"
)

// MockStruct to mimic the CRDT Operation
type MockOp struct {
	ID        string
	Value     int
	Timestamp int64
}

func TestDiskStore_Persistence_Recovery(t *testing.T) {
	// 1. Setup
	tmpFile := "test_store.aof"
	defer os.Remove(tmpFile)

	// Register the type for GOB (Checking if the system handles this or if we need to do it)
	gob.Register(MockOp{})

	// 2. Initialize Store
	ds, err := store.NewDiskStore(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create disk store: %v", err)
	}

	// 3. Write Data
	testKey := "key1"
	testOp := MockOp{ID: "uuid-123", Value: 42, Timestamp: time.Now().UnixNano()}
	
	err = ds.Set("workspace-1", testKey, testOp)
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	// 4. Verify In-Memory
	val, exists := ds.Get("workspace-1", testKey)
	if !exists {
		t.Errorf("Expected value to exist in memory")
	}
	if val.(MockOp).ID != testOp.ID {
		t.Errorf("Value mismatch in memory")
	}

	// 5. Close Store (Simulate Shutdown)
	err = ds.Close()
	if err != nil {
		t.Fatalf("Failed to close store: %v", err)
	}

	// 6. Re-Open Store (Simulate Restart/Recovery)
	ds2, err := store.NewDiskStore(tmpFile)
	if err != nil {
		t.Fatalf("Failed to re-open disk store: %v", err)
	}
	defer ds2.Close()

	// 7. Verify Data Persisted
	val2, exists2 := ds2.Get("workspace-1", testKey)
	if !exists2 {
		t.Fatalf("Expected value to be recovered from disk")
	}
	
	// Check content
	recoveredOp, ok := val2.(MockOp)
	if !ok {
		t.Fatalf("Recovered value is not of type MockOp. Got %T", val2)
	}

	if recoveredOp.ID != testOp.ID || recoveredOp.Value != testOp.Value {
		t.Errorf("Recovered data mismatch. Got %+v, want %+v", recoveredOp, testOp)
	}
}
