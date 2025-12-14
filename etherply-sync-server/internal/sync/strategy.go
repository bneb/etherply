// Package sync provides pluggable synchronization strategies for the EtherPly engine.
//
// CRDT (Automerge) is just one option. Users can choose based on their use case:
//   - Automerge: Full CRDT with mathematical convergence guarantees
//   - LWW: Last-Write-Wins, simpler but no offline convergence
//   - ServerAuthoritative: Server state is always canonical
package sync

import (
	"time"
)

// SyncStrategy defines the contract for document synchronization algorithms.
// Implementations must be stateless - all state is stored in the document bytes.
type SyncStrategy interface {
	// ProcessWrite applies a key-value mutation to the document.
	// Returns the updated document bytes.
	//
	// Parameters:
	//   - current: existing document bytes (nil/empty for new documents)
	//   - key: the field to mutate
	//   - value: the new value (must be JSON-serializable)
	//   - ts: timestamp of the mutation (for audit/ordering)
	//
	// Returns updated document bytes or error.
	ProcessWrite(current []byte, key string, value interface{}, ts time.Time) ([]byte, error)

	// Merge combines two document states (local and remote).
	// This is the core conflict resolution mechanism.
	//
	// Merge contracts:
	//   - Idempotent: Merge(a, a) == a
	//   - Commutative: Merge(a, b) == Merge(b, a)
	//   - Associative: Merge(Merge(a, b), c) == Merge(a, Merge(b, c))
	//
	// Not all strategies guarantee all properties (e.g., LWW is not commutative).
	Merge(local, remote []byte) ([]byte, error)

	// GetState materializes the document as a JSON-like map for API responses.
	GetState(doc []byte) (map[string]interface{}, error)

	// GetHeads returns version identifiers for delta sync.
	// Returns nil/empty for strategies that don't support incremental sync.
	GetHeads(doc []byte) ([]string, error)

	// GetChanges returns changes since given version heads.
	// If since is nil/empty, returns full document for bootstrapping.
	GetChanges(doc []byte, since []string) ([]byte, error)

	// GetHistory returns the change log for the document.
	// Not all strategies support history (LWW does not).
	GetHistory(doc []byte) ([]Change, error)

	// Name returns the strategy identifier for logging and metrics.
	Name() string
}

// Change represents a single mutation in document history.
type Change struct {
	Hash      string `json:"hash"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// StrategyType identifies the sync strategy.
type StrategyType string

const (
	// StrategyAutomerge uses Automerge CRDT for mathematically correct merges.
	StrategyAutomerge StrategyType = "automerge"

	// StrategyLWW uses Last-Write-Wins based on timestamps.
	// WARNING: Does not guarantee convergence in partitioned networks.
	StrategyLWW StrategyType = "lww"

	// StrategyServerAuthoritative server state always wins.
	StrategyServerAuthoritative StrategyType = "server-auth"
)

// NewStrategy creates a SyncStrategy based on the given type.
// Defaults to Automerge if unknown.
func NewStrategy(t StrategyType) SyncStrategy {
	switch t {
	case StrategyLWW:
		return NewLWWStrategy()
	case StrategyServerAuthoritative:
		return NewServerAuthStrategy()
	case StrategyAutomerge:
		fallthrough
	default:
		return NewAutomergeStrategy()
	}
}
