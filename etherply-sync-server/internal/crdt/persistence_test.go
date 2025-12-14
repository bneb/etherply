package crdt_test

import (
	"os"
	"testing"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/crdt"
	"github.com/bneb/etherply/etherply-sync-server/internal/store"
)

func TestEngine_Persistence_Integration(t *testing.T) {
	// 1. Setup
	tmpDir, err := os.MkdirTemp("", "crdt_persistence_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 2. Initialize Store & Engine
	ds, err := store.NewBadgerStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create badger store: %v", err)
	}
	engine := crdt.NewEngine(ds)

	// 3. Process an Operation
	op := crdt.Operation{
		WorkspaceID: "ws-integration",
		Key:         "user:1",
		Value:       "connected",
		Timestamp:   time.Now().UnixMicro(),
	}

	if err := engine.ProcessOperation(op); err != nil {
		t.Fatalf("Failed to process operation: %v", err)
	}

	// 4. Close Store (Simulate Shutdown)
	ds.Close()

	// 5. Re-Open Store (Simulate Restart)
	// We create a NEW Engine to ensure it reads from the store freshly
	ds2, err := store.NewBadgerStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to re-open badger store: %v", err)
	}
	defer ds2.Close()
	engine2 := crdt.NewEngine(ds2)

	// 6. Verify State
	// We need to check if the engine correctly loaded the state.
	// Since ProcessOperation uses LWW check, let's try to process an OLDER operation.
	// If the state was loaded correctly (as Operation struct), the LWW check should discard this older op.
	// If the state was NOT loaded (or loaded as map and failed assertion), the LWW check might behave differently
	// or we can inspect the store directly.

	// Let's inspect store directly via GetFullState
	snapshot, err := engine2.GetFullState("ws-integration")
	if err != nil {
		t.Fatalf("Failed to get full state: %v", err)
	}

	val, ok := snapshot.Data["user:1"]
	if !ok {
		t.Fatalf("Key user:1 missing from restored state")
	}

	// CRITICAL CHECK: Type Assertion
	// The new Automerge implementation stores the value directly at the key.
	// It does NOT store the Operation struct.
	actualValue, ok := val.(string)
	if !ok {
		t.Fatalf("Restored value is not a string. Got %T: %v", val, val)
	}

	if actualValue != "connected" {
		t.Errorf("Value mismatch. Got %v, want 'connected'", actualValue)
	}
}
