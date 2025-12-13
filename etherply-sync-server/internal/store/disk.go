package store

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// DiskStore implements StateStore with a simple Append-Only File (AOF) via gob.
// This is a naive implementation for MVP persistence without external dependencies.
// In production, this would be FoundationDB or CockroachDB.
type DiskStore struct {
	mu      sync.RWMutex
	data    map[string]map[string]interface{}
	aofFile *os.File
	encoder *gob.Encoder
}

// OpLogEntry represents a single mutation logged to disk.
type OpLogEntry struct {
	WorkspaceID string
	Key         string
	Value       interface{}
	Timestamp   int64
}

func NewDiskStore(filePath string) (*DiskStore, error) {
	// Open file in Append|Create|ReadWrite mode
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open persistence file: %w", err)
	}

	ds := &DiskStore{
		data:    make(map[string]map[string]interface{}),
		aofFile: f,
	}

	// Replay existing log
	if err := ds.replay(f); err != nil {
		return nil, fmt.Errorf("failed to replay persistence log: %w", err)
	}

	// Setup encoder for appending
	ds.encoder = gob.NewEncoder(f)

	return ds, nil
}

func (s *DiskStore) replay(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	count := 0
	for {
		var entry OpLogEntry
		err := decoder.Decode(&entry)
		if err == io.EOF {
			break
		}
		if err != nil {
			// Corrupt entry or partial write at end? We log warning and stop replay.
			// "Happy Path is a Trap" -> defensive handling.
			log.Printf("[PERSISTENCE] Warning: corrupt log entry encountered after %d ops: %v", count, err)
			break
		}

		// Apply to memory map
		if _, ok := s.data[entry.WorkspaceID]; !ok {
			s.data[entry.WorkspaceID] = make(map[string]interface{})
		}
		s.data[entry.WorkspaceID][entry.Key] = entry.Value
		count++
	}
	log.Printf("[PERSISTENCE] Replayed %d operations from disk.", count)
	return nil
}

func (s *DiskStore) Get(workspaceID string, key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	workspace, ok := s.data[workspaceID]
	if !ok {
		return nil, false
	}
	val, ok := workspace[key]
	return val, ok
}

func (s *DiskStore) Set(workspaceID string, key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Update Memory
	if _, ok := s.data[workspaceID]; !ok {
		s.data[workspaceID] = make(map[string]interface{})
	}
	s.data[workspaceID][key] = value

	// 2. Append to Disk (AOF)
	entry := OpLogEntry{
		WorkspaceID: workspaceID,
		Key:         key,
		Value:       value,
		Timestamp:   time.Now().UnixMicro(),
	}
	
	// Note: We are holding the lock during disk write here.
	// This is a bottleneck for throughput but ensures durability.
	// Optimization: Buffered Channel for async persistence (Risk: data loss on crash).
	// Decision: Synchronous for now to satisfy "Durability" requirement strictly.
	if err := s.encoder.Encode(entry); err != nil {
		return fmt.Errorf("failed to persist op: %w", err)
	}

	return nil
}

func (s *DiskStore) GetAll(workspaceID string) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	workspace, ok := s.data[workspaceID]
	if !ok {
		return make(map[string]interface{}), nil
	}

	result := make(map[string]interface{})
	for k, v := range workspace {
		result[k] = v
	}
	return result, nil
}

// Close ensures file handles are cleaned up.
func (s *DiskStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.aofFile.Close()
}
