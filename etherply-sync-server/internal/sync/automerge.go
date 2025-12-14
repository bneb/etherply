package sync

import (
	"fmt"
	"time"

	"github.com/automerge/automerge-go"
)

// AutomergeStrategy implements SyncStrategy using the Automerge CRDT library.
// This provides mathematically correct merge semantics with eventual consistency.
type AutomergeStrategy struct{}

// NewAutomergeStrategy creates a new Automerge-based sync strategy.
func NewAutomergeStrategy() *AutomergeStrategy {
	return &AutomergeStrategy{}
}

// Name returns the strategy identifier.
func (s *AutomergeStrategy) Name() string {
	return string(StrategyAutomerge)
}

// ProcessWrite applies a key-value mutation using Automerge.
func (s *AutomergeStrategy) ProcessWrite(current []byte, key string, value interface{}, ts time.Time) ([]byte, error) {
	var doc *automerge.Doc
	var err error

	if len(current) == 0 {
		doc = automerge.New()
	} else {
		doc, err = automerge.Load(current)
		if err != nil {
			return nil, fmt.Errorf("failed to load automerge doc: %w", err)
		}
	}

	if err := doc.Path(key).Set(value); err != nil {
		return nil, fmt.Errorf("failed to set key %q: %w", key, err)
	}

	commitOpts := automerge.CommitOptions{Time: &ts}
	doc.Commit(fmt.Sprintf("set %s", key), commitOpts)

	return doc.Save(), nil
}

// Merge combines local and remote documents using Automerge's CRDT merge.
// This is idempotent, commutative, and associative.
func (s *AutomergeStrategy) Merge(local, remote []byte) ([]byte, error) {
	if len(remote) == 0 {
		return local, nil
	}
	if len(local) == 0 {
		return remote, nil
	}

	localDoc, err := automerge.Load(local)
	if err != nil {
		return nil, fmt.Errorf("failed to load local doc: %w", err)
	}

	remoteDoc, err := automerge.Load(remote)
	if err != nil {
		return nil, fmt.Errorf("failed to load remote doc: %w", err)
	}

	if _, err := localDoc.Merge(remoteDoc); err != nil {
		return nil, fmt.Errorf("failed to merge: %w", err)
	}

	return localDoc.Save(), nil
}

// GetState materializes the document as a map.
func (s *AutomergeStrategy) GetState(doc []byte) (map[string]interface{}, error) {
	if len(doc) == 0 {
		return map[string]interface{}{}, nil
	}

	d, err := automerge.Load(doc)
	if err != nil {
		return nil, err
	}

	rootVal, err := d.Path().Get()
	if err != nil {
		return nil, err
	}

	m, err := automerge.As[map[string]interface{}](rootVal)
	if err != nil {
		return nil, fmt.Errorf("failed to convert root to map: %w", err)
	}

	return m, nil
}

// GetHeads returns the current vector clock heads.
func (s *AutomergeStrategy) GetHeads(doc []byte) ([]string, error) {
	if len(doc) == 0 {
		return []string{}, nil
	}

	d, err := automerge.Load(doc)
	if err != nil {
		return nil, err
	}

	heads := d.Heads()
	result := make([]string, len(heads))
	for i, h := range heads {
		result[i] = h.String()
	}
	return result, nil
}

// GetChanges returns changes since the given heads.
func (s *AutomergeStrategy) GetChanges(doc []byte, since []string) ([]byte, error) {
	if len(doc) == 0 {
		return []byte{}, nil
	}

	d, err := automerge.Load(doc)
	if err != nil {
		return nil, err
	}

	if len(since) == 0 {
		return doc, nil // Full sync
	}

	heads := make([]automerge.ChangeHash, 0, len(since))
	for _, h := range since {
		hash, err := automerge.NewChangeHash(h)
		if err != nil {
			return nil, fmt.Errorf("invalid change hash %q: %w", h, err)
		}
		heads = append(heads, hash)
	}

	changes, err := d.Changes(heads...)
	if err != nil {
		return nil, fmt.Errorf("failed to get changes: %w", err)
	}

	// Frame changes with length prefix
	var buf []byte
	for _, ch := range changes {
		chBytes := ch.Save()
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

// GetHistory returns the commit history.
func (s *AutomergeStrategy) GetHistory(doc []byte) ([]Change, error) {
	if len(doc) == 0 {
		return []Change{}, nil
	}

	d, err := automerge.Load(doc)
	if err != nil {
		return nil, err
	}

	changes, err := d.Changes()
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
