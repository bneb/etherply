// Package crdt implements Conflict-Free Replicated Data Type (CRDT) logic
// for the EtherPly sync engine using the Automerge library.
//
// This replaces the previous "Last-Write-Wins" (LWW) implementation with
// a mathematically correct CRDT that ensures eventual consistency and
// automatic conflict resolution without relying on synchronized clocks.
package crdt

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/automerge/automerge-go"
	"github.com/bneb/etherply/etherply-sync-server/internal/replication"
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

// Snapshot represents a point-in-time view of the document state including vector clock heads.
type Snapshot struct {
	Data  map[string]interface{} `json:"data"`
	Heads []string               `json:"heads"`
}

// Change represents a single commit in history.
type Change struct {
	Hash      string `json:"hash"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// ReplicationCallback is invoked when changes are made locally and need to be broadcast.
// This indirection allows the Engine to remain decoupled from the transport layer.
type ReplicationCallback func(workspaceID string, changes []byte) error

type Engine struct {
	store      store.Store
	logger     *slog.Logger
	mu         sync.Mutex // Global lock for MVP. Ideally should be per-workspace.
	replicator replication.Replicator
	region     string
	serverID   string
}

func NewEngine(s store.Store) *Engine {
	// Default to JSON handler for structured output, writing to stderr
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	return &Engine{
		store:  s,
		logger: logger,
	}
}

// SetReplicator enables multi-region replication.
// When set, the engine will broadcast local changes to peer regions.
func (e *Engine) SetReplicator(r replication.Replicator, region, serverID string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.replicator = r
	e.region = region
	e.serverID = serverID
	e.logger.Info("replication_enabled",
		slog.String("region", region),
		slog.String("server_id", serverID),
	)
}

// HasReplicator returns true if multi-region replication is configured.
func (e *Engine) HasReplicator() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.replicator != nil
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
	// We use the client timestamp if provided to allow ensuring history reflects valid wall clock time from client perspective
	// This is important for "offline" edits.
	now := time.Now()
	commitOpts := automerge.CommitOptions{
		Time: &now,
	}
	if op.Timestamp > 0 {
		// Convert microsecond timestamp to time.Time
		t := time.UnixMicro(op.Timestamp)
		commitOpts.Time = &t
	}
	// "op" msg serves as commit message
	doc.Commit(fmt.Sprintf("set %s", op.Key), commitOpts)

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

	// 5. Broadcast to peer regions (if replication is enabled)
	if e.replicator != nil {
		event := replication.ChangeEvent{
			WorkspaceID:    op.WorkspaceID,
			Changes:        newBytes,
			OriginRegion:   e.region,
			OriginServerID: e.serverID,
			Timestamp:      time.Now(),
		}
		if err := e.replicator.Broadcast(context.Background(), event); err != nil {
			// Non-fatal: local operation still succeeds, log error for alerting
			e.logger.Error("replication_broadcast_failed",
				slog.String("workspace_id", op.WorkspaceID),
				slog.Any("error", err),
			)
		}
	}

	latency := time.Since(start).Milliseconds()
	e.fireSyncOperationMetric(op, latency)

	return nil
}

// GetFullState returns the materialized JSON-like view of the document and its current heads.
func (e *Engine) GetFullState(workspaceID string) (*Snapshot, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	val, exists, err := e.store.Get(workspaceID, "automerge_root")
	if err != nil {
		return nil, err
	}
	if !exists {
		return &Snapshot{
			Data:  map[string]interface{}{},
			Heads: []string{},
		}, nil
	}

	data, ok := val.([]byte)
	if !ok {
		return nil, fmt.Errorf("storage corruption: expected []byte")
	}

	doc, err := automerge.Load(data)
	if err != nil {
		return nil, err
	}

	rootVal, err := doc.Path().Get()
	if err != nil {
		return nil, err
	}

	m, err := automerge.As[map[string]interface{}](rootVal)
	if err != nil {
		return nil, fmt.Errorf("failed to convert root to map: %w", err)
	}

	// Convert heads (ChangeHash) to string slice for transport
	heads := doc.Heads()
	headStrs := make([]string, len(heads))
	for i, h := range heads {
		headStrs[i] = h.String()
	}

	return &Snapshot{
		Data:  m,
		Heads: headStrs,
	}, nil
}

// GetChanges returns the delta (changes) since a specific version vector (heads).
// If 'since' is nil/empty, it returns the full document state for bootstrapping.
func (e *Engine) GetChanges(workspaceID string, since []string) ([]byte, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	val, exists, err := e.store.Get(workspaceID, "automerge_root")
	if err != nil {
		return nil, err
	}
	if !exists {
		return []byte{}, nil
	}

	data, ok := val.([]byte)
	if !ok {
		return nil, fmt.Errorf("storage corruption: expected []byte")
	}

	doc, err := automerge.Load(data)
	if err != nil {
		return nil, err
	}

	// Efficient Delta Sync using automerge.Doc.Changes().
	// doc.Changes(since...) returns all changes that happened after the given heads.
	if len(since) > 0 {
		heads := make([]automerge.ChangeHash, 0, len(since))
		for _, h := range since {
			hash, err := automerge.NewChangeHash(h)
			if err != nil {
				// If a client sends an invalid hash, return error to signal protocol mismatch.
				e.logger.Warn("invalid_change_hash", slog.String("hash", h), slog.Any("error", err))
				return nil, fmt.Errorf("invalid change hash: %w", err)
			}
			heads = append(heads, hash)
		}

		changes, err := doc.Changes(heads...)
		if err != nil {
			return nil, fmt.Errorf("failed to get changes: %w", err)
		}

		// Serialize all changes into a single byte buffer.
		// Each change is prefixed with its length (4 bytes, big-endian) for framing.
		// This allows the client to split them back.
		var buf []byte
		for _, ch := range changes {
			chBytes := ch.Save()
			// Simple concatenation: len(4 bytes) + data
			lenBytes := make([]byte, 4)
			lenBytes[0] = byte(len(chBytes) >> 24)
			lenBytes[1] = byte(len(chBytes) >> 16)
			lenBytes[2] = byte(len(chBytes) >> 8)
			lenBytes[3] = byte(len(chBytes))
			buf = append(buf, lenBytes...)
			buf = append(buf, chBytes...)
		}
		return buf, nil
	}

	// Fallback to full state if no heads provided (initial sync)
	return doc.Save(), nil
}

// GetHistory returns the list of changes for the document.
func (e *Engine) GetHistory(workspaceID string) ([]Change, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	val, exists, err := e.store.Get(workspaceID, "automerge_root")
	if err != nil {
		return nil, err
	}
	if !exists {
		return []Change{}, nil
	}

	data, ok := val.([]byte)
	if !ok {
		return nil, fmt.Errorf("storage corruption: expected []byte")
	}

	doc, err := automerge.Load(data)
	if err != nil {
		return nil, err
	}

	changes, err := doc.Changes()
	if err != nil {
		return nil, err
	}

	history := make([]Change, 0, len(changes))
	for _, c := range changes {
		history = append(history, Change{
			Hash:      c.Hash().String(),
			Message:   c.Message(),
			Timestamp: c.Timestamp().UnixMicro(),
		})
	}

	return history, nil
}

// Stats returns aggregated metrics from the CRDT engine and underlying store.
func (e *Engine) Stats() (map[string]interface{}, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	storeStats, err := e.store.Stats()
	if err != nil {
		return nil, err
	}

	// We can add engine specific stats here if any (e.g. cache size)
	return map[string]interface{}{
		"store": storeStats,
	}, nil
}

// ApplyRemoteChanges merges changes received from peer replicas.
// This is the receiving side of multi-region replication.
// It loads the current document, merges the remote changes, and saves the result.
//
// Note: Automerge's merge operation is idempotent and commutative, meaning
// applying the same changes multiple times or in different orders will
// always result in the same final state.
func (e *Engine) ApplyRemoteChanges(workspaceID string, remoteDoc []byte) error {
	if workspaceID == "" {
		return fmt.Errorf("workspace_id is required")
	}
	if len(remoteDoc) == 0 {
		return nil // No changes to apply
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// Load the remote document
	incoming, err := automerge.Load(remoteDoc)
	if err != nil {
		return fmt.Errorf("failed to load remote document: %w", err)
	}

	// Load or create local document
	var local *automerge.Doc
	val, exists, err := e.store.Get(workspaceID, "automerge_root")
	if err != nil {
		return fmt.Errorf("failed to load local state: %w", err)
	}

	if exists {
		data, ok := val.([]byte)
		if !ok {
			// Corruption: reset to remote state
			e.logger.Warn("local_state_corruption_during_merge",
				slog.String("workspace_id", workspaceID),
			)
			local, err = incoming.Fork()
			if err != nil {
				return fmt.Errorf("failed to fork remote document: %w", err)
			}
		} else {
			local, err = automerge.Load(data)
			if err != nil {
				return fmt.Errorf("failed to load local document: %w", err)
			}
		}
	} else {
		// No local state, use remote as initial state
		local, err = incoming.Fork()
		if err != nil {
			return fmt.Errorf("failed to fork remote document: %w", err)
		}
	}

	// Merge remote changes into local document
	// Automerge guarantees convergence regardless of merge order
	_, err = local.Merge(incoming)
	if err != nil {
		return fmt.Errorf("failed to merge remote changes: %w", err)
	}

	// Save the merged result
	mergedBytes := local.Save()
	if err := e.store.Set(workspaceID, "automerge_root", mergedBytes); err != nil {
		return fmt.Errorf("failed to persist merged state: %w", err)
	}

	e.logger.Debug("remote_changes_applied",
		slog.String("workspace_id", workspaceID),
		slog.Int("merged_size_bytes", len(mergedBytes)),
	)

	return nil
}

// ToJSON helper for debugging
func (op Operation) String() string {
	b, _ := json.Marshal(op)
	return string(b)
}
