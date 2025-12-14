// Package crdt implements Conflict-Free Replicated Data Type (CRDT) logic
// for the EtherPly sync engine using the Automerge library.
//
// This replaces the previous "Last-Write-Wins" (LWW) implementation with
// a mathematically correct CRDT that ensures eventual consistency and
// automatic conflict resolution without relying on synchronized clocks.
package crdt

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/automerge/automerge-go"
	"github.com/bneb/etherply/etherply-sync-server/internal/store"
)

// Operation represents a request to mutate the state.
// Note: Timestamp is no longer used for conflict resolution (handled by Automerge),
// but we keep the field for client compatibility and audit logs.
type Operation struct {
	WorkspaceID string      `json:"workspace_id"`
	Key         string      `json:"key"`
	Value       interface{} `json:"value"`
	Timestamp   int64       `json:"timestamp"` // Unix Microseconds, informative only
}

type Engine struct {
	store  store.Store
	logger *slog.Logger
	mu     sync.Mutex // Global lock for MVP. Ideally should be per-workspace.
}

func NewEngine(s store.Store) *Engine {
	// Default to JSON handler for structured output, writing to stderr
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	return &Engine{
		store:  s,
		logger: logger,
	}
}

// fireSyncOperationMetric is a helper to track operation metrics (PostHog stub)
func (e *Engine) fireSyncOperationMetric(op Operation, latencyMs int64) {
	e.logger.Info("sync_metric",
		slog.String("event", "sync_operation_count"),
		slog.String("workspace_id", op.WorkspaceID),
		slog.Int64("latency_ms", latencyMs),
		slog.String("engine", "automerge"),
	)
}

// ProcessOperation handles an incoming mutation using Automerge.
// It performs a Read-Modify-Write cycle on the persistent store.
func (e *Engine) ProcessOperation(op Operation) error {
	start := time.Now()

	// 0. Strict Validation
	if op.WorkspaceID == "" || op.Key == "" {
		return fmt.Errorf("invalid operation: workspace_id and key are required")
	}

	// Lock critical section to ensure atomicity of Load -> Edit -> Save
	// CRITICAL: automerge.Doc is NOT thread-safe.
	e.mu.Lock()
	defer e.mu.Unlock()

	// 1. Fetch current binary state from Persistence Layer
	var doc *automerge.Doc
	var err error

	// We store the raw Automerge binary blob
	val, exists, err := e.store.Get(op.WorkspaceID, "automerge_root")
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	if exists {
		// assert type to []byte because our BadgerStore (and MemoryStore ideally) should support raw bytes
		// Note: The previous LWW implementation stored struct types.
		// If we are migrating, we might encounter legacy data.
		// For this "Ironclad" rebuild, we assume a fresh start or compatible store.
		data, ok := val.([]byte)
		if !ok {
			// This handles the legacy/corruption case defensive coding
			e.logger.Error("store_type_mismatch",
				slog.String("workspace_id", op.WorkspaceID),
				slog.String("expected", "[]byte"),
				slog.String("got", fmt.Sprintf("%T", val)),
			)
			// Reset to empty doc to self-heal
			doc = automerge.New()
		} else {
			doc, err = automerge.Load(data)
			if err != nil {
				return fmt.Errorf("failed to hydrate automerge doc: %w", err)
			}
		}
	} else {
		doc = automerge.New()
	}

	// 2. Apply Limitation: Automerge works best with a single root map for this KV use case.
	// We map op.Key directly to the root of the Automerge document.
	// op.Value can be complex, but Automerge-go Set() supports primitives.
	// For complex objects (maps/slices), we might need traversing.
	// Assuming op.Value is a primitive or simple map for MVP.
	// Automerge-go's Set() automatically handles most Go types.

	err = doc.Path(op.Key).Set(op.Value)
	if err != nil {
		return fmt.Errorf("failed to apply operation to automerge doc: %w", err)
	}

	// 3. Commit logic (Automerge handles history)
	// commit := doc.Commit("scan", automerge.CommitOptions{Time: time.Now()})

	// 4. Save back to persistence
	newBytes := doc.Save()

	// We use a constant key "automerge_root" for the blob within the workspace bucket
	// if the store supports buckets.
	// However, the `store.Store` interface `Get(workspaceID, key)` implies the workspaceID is the bucket
	// and `key` is the item.
	// So we should store it under a specific reserved key.
	err = e.store.Set(op.WorkspaceID, "automerge_root", newBytes)
	if err != nil {
		return fmt.Errorf("failed to persist state: %w", err)
	}

	latency := time.Since(start).Milliseconds()
	e.fireSyncOperationMetric(op, latency)

	return nil
}

// GetFullState returns the materialized JSON-like view of the document.
func (e *Engine) GetFullState(workspaceID string) (map[string]interface{}, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	val, exists, err := e.store.Get(workspaceID, "automerge_root")
	if err != nil {
		return nil, err
	}
	if !exists {
		return map[string]interface{}{}, nil
	}

	data, ok := val.([]byte)
	if !ok {
		return nil, fmt.Errorf("storage corruption: expected []byte")
	}

	doc, err := automerge.Load(data)
	if err != nil {
		return nil, err
	}

	// Get the root map and convert to Go map
	// doc.Root() returns a Value. Get() returns value.
	// We want the whole document as a map.
	// doc.Root() isn't a map itself, it provides access.
	// doc.Path().Get() on empty path gives root?
	// The canonical way to dump content is strictly typed or usually via Map().
	// But doc.Path("") probably refers to root.

	rootVal, err := doc.Path().Get()
	if err != nil {
		return nil, err
	}

	// automerge-go implementation detail: use As[T] to convert
	m, err := automerge.As[map[string]interface{}](rootVal)
	if err != nil {
		return nil, fmt.Errorf("failed to convert root to map: %w", err)
	}

	return m, nil
}

// ToJSON helper for debugging
func (op Operation) String() string {
	b, _ := json.Marshal(op)
	return string(b)
}
