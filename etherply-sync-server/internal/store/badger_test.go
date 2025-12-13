package store_test

import (
	"os"
	"testing"

	"github.com/bneb/etherply/etherply-sync-server/internal/store"
)

func TestBadgerStore_SetGet(t *testing.T) {
	dir, err := os.MkdirTemp("", "badger-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	bs, err := store.NewBadgerStore(dir)
	if err != nil {
		t.Fatalf("Failed to create badger store: %v", err)
	}
	defer bs.Close()

	// Test 1: Set and Get
	err = bs.Set("workspace-1", "key1", "value1")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, exists, err := bs.Get("workspace-1", "key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !exists {
		t.Error("Expected value to exist")
	}
	if val.(string) != "value1" {
		t.Errorf("Expected 'value1', got %v", val)
	}
}

func TestBadgerStore_GetAll(t *testing.T) {
	dir, err := os.MkdirTemp("", "badger-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	bs, err := store.NewBadgerStore(dir)
	if err != nil {
		t.Fatalf("Failed to create badger store: %v", err)
	}
	defer bs.Close()

	bs.Set("ws1", "k1", "v1")
	bs.Set("ws1", "k2", "v2")
	bs.Set("ws2", "k3", "v3") // Different workspace

	data, err := bs.GetAll("ws1")
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(data) != 2 {
		t.Errorf("Expected 2 items, got %d", len(data))
	}
	if data["k1"] != "v1" || data["k2"] != "v2" {
		t.Errorf("GetAll returned incorrect data: %v", data)
	}
}

func TestBadgerStore_Persistence(t *testing.T) {
	dir, err := os.MkdirTemp("", "badger-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Open, Write, Close
	bs, err := store.NewBadgerStore(dir)
	if err != nil {
		t.Fatalf("Failed to create badger store: %v", err)
	}
	err = bs.Set("ws1", "k1", "persistent")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	bs.Close()

	// Reopen
	bs2, err := store.NewBadgerStore(dir)
	if err != nil {
		t.Fatalf("Failed to reopen badger store: %v", err)
	}
	defer bs2.Close()

	val, exists, err := bs2.Get("ws1", "k1")
	if err != nil {
		t.Fatalf("Get after reopen failed: %v", err)
	}
	if !exists {
		t.Error("Expected value to persist")
	}
	if val.(string) != "persistent" {
		t.Errorf("Expected 'persistent', got %v", val)
	}
}
