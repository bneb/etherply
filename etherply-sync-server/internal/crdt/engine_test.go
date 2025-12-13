package crdt

import (
	"testing"
)

// MockStore for testing
type MockStore struct {
	data map[string]interface{}
}

func NewMockStore() *MockStore {
	return &MockStore{data: make(map[string]interface{})}
}

func (m *MockStore) Get(workspaceID, key string) (interface{}, bool) {
	val, ok := m.data[key]
	return val, ok
}

func (m *MockStore) GetAll(workspaceID string) (map[string]interface{}, error) {
	return m.data, nil
}

func (m *MockStore) Set(workspaceID, key string, value interface{}) error {
	m.data[key] = value
	return nil
}

func (m *MockStore) Close() error { return nil }

func TestEngine_ProcessOperation(t *testing.T) {
	mockStore := NewMockStore()
	engine := NewEngine(mockStore)

	// Test 1: Validation
	err := engine.ProcessOperation(Operation{WorkspaceID: "", Key: "foo", Value: "bar"})
	if err == nil {
		t.Errorf("Expected error for empty workspace_id")
	}

	// Test 2: Valid Operation
	op1 := Operation{
		WorkspaceID: "ws1",
		Key:         "doc1",
		Value:       "Hello",
		Timestamp:   1000,
	}
	if err := engine.ProcessOperation(op1); err != nil {
		t.Errorf("ProcessOperation failed: %v", err)
	}

	// Test 3: LWW (Older timestamp should be ignored)
	opOld := Operation{
		WorkspaceID: "ws1",
		Key:         "doc1",
		Value:       "Old",
		Timestamp:   900,
	}
	if err := engine.ProcessOperation(opOld); err != nil {
		t.Errorf("ProcessOperation failed: %v", err)
	}
	// Verify value hasn't changed
	val, _ := mockStore.Get("ws1", "doc1")
	if val.(Operation).Value != "Hello" {
		t.Errorf("LWW failed: Value overwritten by older op")
	}

	// Test 4: LWW (Newer timestamp should overwrite)
	opNew := Operation{
		WorkspaceID: "ws1",
		Key:         "doc1",
		Value:       "World",
		Timestamp:   2000,
	}
	if err := engine.ProcessOperation(opNew); err != nil {
		t.Errorf("ProcessOperation failed: %v", err)
	}
	val, _ = mockStore.Get("ws1", "doc1")
	if val.(Operation).Value != "World" {
		t.Errorf("LWW failed: Value did not update")
	}
}
