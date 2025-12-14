// Package crdt implements the synchronization engine for EtherPly.
//
// The engine is now strategy-agnostic and supports multiple synchronization backends:
//   - Automerge (CRDT): Full conflict-free merge with convergence guarantees
//   - LWW: Last-Write-Wins, simpler semantics for non-collaborative data
//   - Server-Authoritative: Server state always wins
//
// Strategy selection is done via EngineOption at construction time.
package crdt

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	gosync "sync"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/replication"
	"github.com/bneb/etherply/etherply-sync-server/internal/store"
	"github.com/bneb/etherply/etherply-sync-server/internal/sync"
)

// docKey is the reserved storage key for document blobs.
const docKey = "sync_doc"

// Operation represents a verified request to mutate the document state.
// Timestamp is strictly ordered by the server (Unix Microseconds).
type Operation struct {
	WorkspaceID string      `json:"workspace_id"`
	Key         string      `json:"key"`
	Value       interface{} `json:"value"`
	Timestamp   int64       `json:"timestamp"` // Unix Microseconds
}

// Snapshot represents a point-in-time view of the document state including version heads.
type Snapshot struct {
	Data  map[string]interface{} `json:"data"`
	Heads []string               `json:"heads"`
}

// Change represents a single commit in history (re-exported from sync package).
type Change = sync.Change

// ReplicationCallback is invoked when changes are made locally and need to be broadcast.
type ReplicationCallback func(workspaceID string, changes []byte) error

// Engine orchestrates document synchronization using a pluggable strategy.
type Engine struct {
	store    store.Store
	strategy sync.SyncStrategy
	logger   *slog.Logger
	// mu protects the engine state.
	// Note: Granular locking per-workspace is planned for future optimization
	// when contention on the global lock becomes a bottleneck (>100k OPS).
	mu         gosync.Mutex
	replicator replication.Replicator
	region     string
	serverID   string
}

// EngineConfig holds engine configuration.
type EngineConfig struct {
	Strategy sync.SyncStrategy
	Logger   *slog.Logger
}

// EngineOption configures the engine.
type EngineOption func(*EngineConfig)

// WithStrategy sets the synchronization strategy.
func WithStrategy(s sync.SyncStrategy) EngineOption {
	return func(cfg *EngineConfig) {
		cfg.Strategy = s
	}
}

// WithLogger sets a custom logger.
func WithLogger(l *slog.Logger) EngineOption {
	return func(cfg *EngineConfig) {
		cfg.Logger = l
	}
}

// NewEngine creates a new sync engine with the given store and options.
// Defaults to Automerge strategy for backward compatibility.
func NewEngine(s store.Store, opts ...EngineOption) *Engine {
	cfg := &EngineConfig{
		Strategy: sync.NewAutomergeStrategy(),
		Logger:   slog.New(slog.NewJSONHandler(os.Stderr, nil)),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	cfg.Logger.Info("engine_initialized",
		slog.String("strategy", cfg.Strategy.Name()),
	)

	return &Engine{
		store:    s,
		strategy: cfg.Strategy,
		logger:   cfg.Logger,
	}
}

// Strategy returns the current sync strategy name.
func (e *Engine) Strategy() string {
	return e.strategy.Name()
}

// SetReplicator enables multi-region replication.
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

// fireSyncOperationMetric tracks operation metrics.
func (e *Engine) fireSyncOperationMetric(op Operation, latencyMs int64) {
	e.logger.Info("sync_metric",
		slog.String("event", "sync_operation_count"),
		slog.String("workspace_id", op.WorkspaceID),
		slog.Int64("latency_ms", latencyMs),
		slog.String("strategy", e.strategy.Name()),
	)
}

// ProcessOperation handles an incoming mutation using the configured strategy.
func (e *Engine) ProcessOperation(op Operation) error {
	start := time.Now()

	// Strict Validation
	if op.WorkspaceID == "" {
		return fmt.Errorf("invalid operation: workspace_id is missing")
	}
	if op.Key == "" {
		return fmt.Errorf("invalid operation: key is missing")
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// 1. Load current document
	current, err := e.loadDoc(op.WorkspaceID)
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	// 2. Determine timestamp
	ts := time.Now()
	if op.Timestamp > 0 {
		ts = time.UnixMicro(op.Timestamp)
	}

	// 3. Apply mutation via strategy
	newDoc, err := e.strategy.ProcessWrite(current, op.Key, op.Value, ts)
	if err != nil {
		return fmt.Errorf("failed to apply operation: %w", err)
	}

	// 4. Persist
	if err := e.saveDoc(op.WorkspaceID, newDoc); err != nil {
		return fmt.Errorf("failed to persist state: %w", err)
	}

	// 5. Broadcast to peer regions (if replication enabled)
	if e.replicator != nil {
		event := replication.ChangeEvent{
			WorkspaceID:    op.WorkspaceID,
			Changes:        newDoc,
			OriginRegion:   e.region,
			OriginServerID: e.serverID,
			Timestamp:      time.Now(),
		}
		if err := e.replicator.Broadcast(context.Background(), event); err != nil {
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

// GetFullState returns the materialized view of the document and its current heads.
func (e *Engine) GetFullState(workspaceID string) (*Snapshot, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	doc, err := e.loadDoc(workspaceID)
	if err != nil {
		return nil, err
	}

	data, err := e.strategy.GetState(doc)
	if err != nil {
		return nil, err
	}

	heads, err := e.strategy.GetHeads(doc)
	if err != nil {
		return nil, err
	}

	return &Snapshot{
		Data:  data,
		Heads: heads,
	}, nil
}

// GetChanges returns the delta since a specific version vector.
func (e *Engine) GetChanges(workspaceID string, since []string) ([]byte, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	doc, err := e.loadDoc(workspaceID)
	if err != nil {
		return nil, err
	}

	return e.strategy.GetChanges(doc, since)
}

// GetHistory returns the list of changes for the document.
func (e *Engine) GetHistory(workspaceID string) ([]Change, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	doc, err := e.loadDoc(workspaceID)
	if err != nil {
		return nil, err
	}

	return e.strategy.GetHistory(doc)
}

// Stats returns aggregated metrics from the engine and store.
func (e *Engine) Stats() (map[string]interface{}, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	storeStats, err := e.store.Stats()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"store":    storeStats,
		"strategy": e.strategy.Name(),
	}, nil
}

// ApplyRemoteChanges merges changes received from peer replicas.
func (e *Engine) ApplyRemoteChanges(workspaceID string, remoteDoc []byte) error {
	if workspaceID == "" {
		return fmt.Errorf("workspace_id is required")
	}
	if len(remoteDoc) == 0 {
		return nil
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	local, err := e.loadDoc(workspaceID)
	if err != nil {
		return fmt.Errorf("failed to load local state: %w", err)
	}

	merged, err := e.strategy.Merge(local, remoteDoc)
	if err != nil {
		return fmt.Errorf("failed to merge: %w", err)
	}

	if err := e.saveDoc(workspaceID, merged); err != nil {
		return fmt.Errorf("failed to persist merged state: %w", err)
	}

	e.logger.Debug("remote_changes_applied",
		slog.String("workspace_id", workspaceID),
		slog.Int("merged_size_bytes", len(merged)),
		slog.String("strategy", e.strategy.Name()),
	)

	return nil
}

// loadDoc retrieves document bytes from the store.
func (e *Engine) loadDoc(workspaceID string) ([]byte, error) {
	val, exists, err := e.store.Get(workspaceID, docKey)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	data, ok := val.([]byte)
	if !ok {
		e.logger.Warn("store_type_mismatch",
			slog.String("workspace_id", workspaceID),
			slog.String("expected", "[]byte"),
			slog.String("got", fmt.Sprintf("%T", val)),
		)
		return nil, nil // Treat as empty, strategy will initialize
	}
	return data, nil
}

// saveDoc persists document bytes to the store.
func (e *Engine) saveDoc(workspaceID string, doc []byte) error {
	return e.store.Set(workspaceID, docKey, doc)
}

// String helper for Operation debugging.
func (op Operation) String() string {
	b, _ := json.Marshal(op)
	return string(b)
}
