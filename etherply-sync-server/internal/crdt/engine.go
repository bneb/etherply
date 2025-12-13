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
	"log"
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
	store store.StateStore
}

func NewEngine(s store.StateStore) *Engine {
	return &Engine{
		store: s,
	}
}

// PrepareMetricData is a helper to track operation metrics (PostHog stub)
func (e *Engine) fireSyncOperationMetric(op Operation, latencyMs int64) {
	// In a real implementation, this would send an event to PostHog.
	// Metric: sync_operation_count
	// Properties: data_bytes_transferred (approx), workspace_id, latency_ms
	log.Printf("[METRIC] sync_operation_count | workspace=%s latency=%dms", op.WorkspaceID, latencyMs)
}

// ProcessOperation handles an incoming CRDT operation.
// It implements LWW (Last-Write-Wins) conflict resolution.
func (e *Engine) ProcessOperation(op Operation) error {
	start := time.Now()

	// 1. Fetch current state to check timestamp (LWW check)
	// In a real LWW Register, we store (Value, Timestamp).
	// Here for simplicity of the "in-memory map" MVP, we might just overwrite.
	// To do LWW correctly, the Store needs to store the timestamp too.
	// For this MVP step 1, we will trust the client simply sends new data.
	// But let's add a basic check if we were using a structured object.

	// For the demo "Magic Moment", we just blindly accept the latest arrival if we assume consistent clocks
	// or if we just want to show propagation.
	// However, the mandate says "CRDT-based state synchronization".
	// So let's just save valid operations.

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
			// For MVP, we'll just log if it's > 1 minute ahead.
			if op.Timestamp > time.Now().Add(1*time.Minute).UnixMicro() {
				log.Printf("[LWW] Warning: received future timestamp for key=%s (ts=%d)", op.Key, op.Timestamp)
			}

			if op.Timestamp <= currentOp.Timestamp {
				// Incoming op is older or equal. LWW rule: discard.
				// We do not return error, as this is valid eventual consistency convergence.
				log.Printf("[LWW] Discarding stale op for key=%s (newTs=%d, currTs=%d)", op.Key, op.Timestamp, currentOp.Timestamp)
				return nil
			}
		} else {
			// CRITICAL: Type Assertion Failed.
			// This implies data corruption or schema drift.
			// We choose to OVERWRITE (Self-Heal) but we must log strictly.
			log.Printf("[CRITICAL] State Corruption Detected: Could not assert current value type for key=%s. Overwriting to self-heal.", op.Key)
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
