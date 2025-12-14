package sync

import (
	"encoding/json"
	"fmt"
	"time"
)

// lwwEntry represents a single key-value pair with timestamp for LWW resolution.
type lwwEntry struct {
	Value     interface{} `json:"v"`
	Timestamp int64       `json:"ts"` // Unix microseconds
}

// lwwDocument is the internal structure for LWW storage.
type lwwDocument struct {
	Entries map[string]lwwEntry `json:"entries"`
}

// LWWStrategy implements SyncStrategy using Last-Write-Wins semantics.
// The write with the highest timestamp always wins.
//
// WARNING: This strategy does NOT guarantee convergence in partitioned networks
// where clocks may be out of sync. Use Automerge for stronger guarantees.
type LWWStrategy struct{}

// NewLWWStrategy creates a new Last-Write-Wins sync strategy.
func NewLWWStrategy() *LWWStrategy {
	return &LWWStrategy{}
}

// Name returns the strategy identifier.
func (s *LWWStrategy) Name() string {
	return string(StrategyLWW)
}

// ProcessWrite applies a mutation, keeping the value with the highest timestamp.
func (s *LWWStrategy) ProcessWrite(current []byte, key string, value interface{}, ts time.Time) ([]byte, error) {
	doc := s.loadOrCreate(current)

	tsMicro := ts.UnixMicro()
	existing, exists := doc.Entries[key]

	// Only update if new timestamp is greater
	if !exists || tsMicro > existing.Timestamp {
		doc.Entries[key] = lwwEntry{
			Value:     value,
			Timestamp: tsMicro,
		}
	}

	return json.Marshal(doc)
}

// Merge combines documents by taking the highest timestamp for each key.
func (s *LWWStrategy) Merge(local, remote []byte) ([]byte, error) {
	localDoc := s.loadOrCreate(local)
	remoteDoc := s.loadOrCreate(remote)

	for key, remoteEntry := range remoteDoc.Entries {
		localEntry, exists := localDoc.Entries[key]
		if !exists || remoteEntry.Timestamp > localEntry.Timestamp {
			localDoc.Entries[key] = remoteEntry
		}
	}

	return json.Marshal(localDoc)
}

// GetState materializes the document, stripping timestamps.
func (s *LWWStrategy) GetState(doc []byte) (map[string]interface{}, error) {
	d := s.loadOrCreate(doc)
	result := make(map[string]interface{}, len(d.Entries))
	for key, entry := range d.Entries {
		result[key] = entry.Value
	}
	return result, nil
}

// GetHeads returns a single "head" based on max timestamp (for compatibility).
func (s *LWWStrategy) GetHeads(doc []byte) ([]string, error) {
	d := s.loadOrCreate(doc)
	var maxTs int64
	for _, entry := range d.Entries {
		if entry.Timestamp > maxTs {
			maxTs = entry.Timestamp
		}
	}
	if maxTs == 0 {
		return []string{}, nil
	}
	return []string{fmt.Sprintf("%d", maxTs)}, nil
}

// GetChanges returns full document (LWW does not support incremental sync).
func (s *LWWStrategy) GetChanges(doc []byte, since []string) ([]byte, error) {
	// LWW doesn't have fine-grained history, return full state
	return doc, nil
}

// GetHistory returns empty - LWW does not track history.
func (s *LWWStrategy) GetHistory(doc []byte) ([]Change, error) {
	return []Change{}, nil
}

func (s *LWWStrategy) loadOrCreate(data []byte) *lwwDocument {
	if len(data) == 0 {
		return &lwwDocument{Entries: make(map[string]lwwEntry)}
	}
	var doc lwwDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return &lwwDocument{Entries: make(map[string]lwwEntry)}
	}
	if doc.Entries == nil {
		doc.Entries = make(map[string]lwwEntry)
	}
	return &doc
}
