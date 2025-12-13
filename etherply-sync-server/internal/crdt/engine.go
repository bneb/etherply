// Package crdt implements Conflict-Free Replicated Data Type (CRDT) logic
// for the EtherPly sync engine. Currently implements Last-Write-Wins (LWW)
// register semantics: when concurrent updates occur, the operation with the
// highest timestamp wins.
//
// Operations include a client-provided timestamp (Unix microseconds) for
// conflict resolution. The engine logs and discards stale operations that
// arrive out of order, ensuring eventual consistency across all clients.
//
// Future iterations may implement more sophisticated CRDTs such as RGA for
// collaborative text editing or OR-Set for collection operations.
package crdt


import (
	"encoding/gob"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/store"
)

func init() {
	gob.Register(Operation{})
}

// Operation represents a single mutation to the state.
// For this MVP, we use Last-Write-Wins (LWW) Key-Value logic.
type Operation struct {
	WorkspaceID string      `json:"workspace_id"`
	Key         string      `json:"key"`
	Value       interface{} `json:"value"`
	Timestamp   int64       `json:"timestamp"` // Unix Microseconds
}

type Engine struct {
	store  store.StateStore
	logger *slog.Logger
}

func NewEngine(s store.StateStore) *Engine {
	// Default to JSON handler for structured output, writing to stderr
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	return &Engine{
		store:  s,
		logger: logger,
	}
}

// PrepareMetricData is a helper to track operation metrics (PostHog stub)
func (e *Engine) fireSyncOperationMetric(op Operation, latencyMs int64) {
	// Metric: sync_operation_count
	e.logger.Info("sync_metric",
		slog.String("event", "sync_operation_count"),
		slog.String("workspace_id", op.WorkspaceID),
		slog.Int64("latency_ms", latencyMs),
		slog.Int("bytes_approx", len(fmt.Sprintf("%v", op.Value))), // Rough approximation
	)
}

// ProcessOperation handles an incoming CRDT operation.
// It implements LWW (Last-Write-Wins) conflict resolution.
func (e *Engine) ProcessOperation(op Operation) error {
	start := time.Now()

	// 0. Strict Validation (Anti-Fuzzing / Defensive Coding)
	if op.WorkspaceID == "" || op.Key == "" {
		return fmt.Errorf("invalid operation: workspace_id and key are required")
	}

	// 1. LWW Conflict Resolution (Per PRD-001 Story 3)
	// We must fetch the current state to compare timestamps.
	currentVal, exists := e.store.Get(op.WorkspaceID, op.Key)
	if exists {
		// Attempt to type assert to Operation
		if currentOp, ok := currentVal.(Operation); ok {
			// Clock Skew Protection: Reject timestamps significantly in the future?
			if op.Timestamp > time.Now().Add(1*time.Minute).UnixMicro() {
				e.logger.Warn("clock_skew_detected",
					slog.String("key", op.Key),
					slog.Int64("op_ts", op.Timestamp),
					slog.Int64("server_ts", time.Now().UnixMicro()),
				)
			}

			if op.Timestamp <= currentOp.Timestamp {
				// Incoming op is older or equal. LWW rule: discard.
				e.logger.Info("discarding_stale_op",
					slog.String("key", op.Key),
					slog.Int64("new_ts", op.Timestamp),
					slog.Int64("curr_ts", currentOp.Timestamp),
				)
				return nil
			}
		} else {
			// CRITICAL: Type Assertion Failed.
			// This implies data corruption or schema drift.
			// We choose to OVERWRITE (Self-Heal) but we must log strictly.
			e.logger.Error("state_corruption_detected",
				slog.String("key", op.Key),
				slog.String("error", "could not assert current value type"),
				slog.String("action", "overwriting_to_self_heal"),
			)
		}
	}

	err := e.store.Set(op.WorkspaceID, op.Key, op)
	if err != nil {
		return err
	}

	latency := time.Since(start).Milliseconds()
	e.fireSyncOperationMetric(op, latency)

	return nil
}

// CreateSyncMessage generates the full state message for a new client.
func (e *Engine) GetFullState(workspaceID string) (map[string]interface{}, error) {
	return e.store.GetAll(workspaceID)
}
