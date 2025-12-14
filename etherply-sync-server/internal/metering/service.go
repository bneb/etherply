package metering

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/store"
)

type MetricType string

const (
	MetricMessagesSent     MetricType = "msg_sent"
	MetricMessagesReceived MetricType = "msg_recv"
	MetricStorageBytes     MetricType = "storage_bytes"
	MetricConnections      MetricType = "connections" // Gauge-like, simpler to just sample
)

// Service defines the metering interface.
type Service interface {
	Record(workspaceID string, metric MetricType, delta int64) error
	GetUsage(workspaceID string, metric MetricType, start, end time.Time) (int64, error)
}

type BadgerMeteringService struct {
	store store.Store
}

func NewBadgerMeteringService(s store.Store) *BadgerMeteringService {
	return &BadgerMeteringService{store: s}
}

// Record increments a usage counter for a given day.
// Key format: usage:<workspace_id>:<metric>:<YYYY-MM-DD>
// Value: int64 (counter)
func (s *BadgerMeteringService) Record(workspaceID string, metric MetricType, delta int64) error {
	// In a high throughput system, we would buffer this in memory and flush periodically.
	// For "One-Shot" simplicity, we write directly to Badger's merge operator or simpler Read-Modify-Write.
	// Badger has a Merge operator for counters, but regular Get/Set is easier to implement quickly without low-level Badger structs.
	// We'll use day-bucketed keys.

	today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("usage:%s:%s:%s", workspaceID, metric, today)

	// We assume the store.Store interface exposes the underlying DB or we cast it.
	// Ideally we extend the interface. Let's try to add a specific method to Store or use a raw method.
	// Since we are inside the 'metering' package, we might not have access to internals of 'store'.
	// Let's assume we can use s.store.Update if defined, or we just rely on Get/Set (race condition possible if multiple pods, but binary is single pod per region usually).
	// Actually, NATS handles replication, but the local Badger is local state.

	// CAS loop for correctness?
	// Or assume single writer per key?
	// Let's implement a naive Read-Modify-Write for now.
	// Optimization: Add 'Increment' to Store interface later.

	// We need to cast to *store.BadgerStore to access the DB directly or use Get/Set.

	// Let's assume we can use:
	// val, exists, _ := s.store.Get(key)
	// new := val + delta
	// s.store.Set(key, new)

	// However, store.Get requires a composite key helper or just raw key?
	// The current store.Get takes (bucket, key). We can treat bucket as "metering".

	// Use "sys:metering" namespace
	// Note: The key is already "usage:wid:metric:date".
	// We treat "sys:metering" as the namespace.

	valBytes, exists, err := s.store.Get("sys:metering", key)
	if err != nil {
		return err
	}

	var current int64 = 0
	if exists && valBytes != nil {
		current = int64(binary.LittleEndian.Uint64(valBytes.([]byte)))
	}

	newVal := current + delta
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(newVal))

	return s.store.Set("sys:metering", key, buf)
}

func (s *BadgerMeteringService) GetUsage(workspaceID string, metric MetricType, start, end time.Time) (int64, error) {
	var total int64 = 0
	// Iterate specific days.
	for d := start; d.Before(end) || d.Equal(end); d = d.AddDate(0, 0, 1) {
		dayStr := d.Format("2006-01-02")
		key := fmt.Sprintf("usage:%s:%s:%s", workspaceID, metric, dayStr)

		valBytes, exists, err := s.store.Get("sys:metering", key)
		if err != nil {
			return 0, err
		}
		if exists && valBytes != nil {
			total += int64(binary.LittleEndian.Uint64(valBytes.([]byte)))
		}
	}
	return total, nil
}
