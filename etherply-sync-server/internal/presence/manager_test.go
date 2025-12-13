package presence

import (
	"sync"
	"testing"
)

func TestManager_AddRemoveUser(t *testing.T) {
	m := NewManager()

	// Test 1: Add user to workspace
	user := User{UserID: "user-1", Status: "online"}
	m.AddUser("workspace-1", user)

	users := m.GetUsers("workspace-1")
	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}
	if users[0].UserID != "user-1" {
		t.Errorf("Expected user-1, got %s", users[0].UserID)
	}

	// Test 2: Add second user
	user2 := User{UserID: "user-2", Status: "idle"}
	m.AddUser("workspace-1", user2)

	users = m.GetUsers("workspace-1")
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	// Test 3: Remove user
	m.RemoveUser("workspace-1", "user-1")
	users = m.GetUsers("workspace-1")
	if len(users) != 1 {
		t.Errorf("Expected 1 user after removal, got %d", len(users))
	}
	if users[0].UserID != "user-2" {
		t.Errorf("Expected user-2, got %s", users[0].UserID)
	}

	// Test 4: Remove last user (workspace should be cleaned up)
	m.RemoveUser("workspace-1", "user-2")
	users = m.GetUsers("workspace-1")
	if len(users) != 0 {
		t.Errorf("Expected 0 users after all removals, got %d", len(users))
	}
}

func TestManager_GetUsers_EmptyWorkspace(t *testing.T) {
	m := NewManager()

	// Getting users from non-existent workspace should return nil or empty slice
	// In Go, nil slice is functionally equivalent to empty slice for range/len operations
	users := m.GetUsers("non-existent")
	if len(users) != 0 {
		t.Errorf("Expected 0 users, got %d", len(users))
	}
}

func TestManager_MultipleWorkspaces(t *testing.T) {
	m := NewManager()

	// Add users to different workspaces
	m.AddUser("ws-a", User{UserID: "alice", Status: "online"})
	m.AddUser("ws-b", User{UserID: "bob", Status: "online"})
	m.AddUser("ws-a", User{UserID: "charlie", Status: "idle"})

	usersA := m.GetUsers("ws-a")
	usersB := m.GetUsers("ws-b")

	if len(usersA) != 2 {
		t.Errorf("Expected 2 users in ws-a, got %d", len(usersA))
	}
	if len(usersB) != 1 {
		t.Errorf("Expected 1 user in ws-b, got %d", len(usersB))
	}
}

func TestManager_UpdateUserStatus(t *testing.T) {
	m := NewManager()

	// Adding a user with same ID should update (overwrite)
	m.AddUser("ws-1", User{UserID: "user-1", Status: "online"})
	m.AddUser("ws-1", User{UserID: "user-1", Status: "idle"})

	users := m.GetUsers("ws-1")
	if len(users) != 1 {
		t.Errorf("Expected 1 user (update, not duplicate), got %d", len(users))
	}
	if users[0].Status != "idle" {
		t.Errorf("Expected status 'idle', got %s", users[0].Status)
	}
}

func TestManager_Concurrency(t *testing.T) {
	m := NewManager()

	// Test concurrent access doesn't cause race conditions
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			user := User{UserID: "user", Status: "online"}
			m.AddUser("ws", user)
			m.GetUsers("ws")
			m.RemoveUser("ws", "user")
		}(i)
	}
	wg.Wait()

	// After all goroutines complete, workspace should be empty or have one user
	// (depending on timing). Just verify no panic occurred.
	users := m.GetUsers("ws")
	if users != nil && len(users) > 1 {
		t.Errorf("Expected 0 or 1 user after concurrent ops, got %d", len(users))
	}
}
