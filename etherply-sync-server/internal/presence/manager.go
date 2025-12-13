// Package presence provides ephemeral user presence tracking for workspaces.
// It maintains a real-time view of which users are connected to each workspace
// and their current status (online, idle, etc.). This data is not persisted
// and is rebuilt as connections are established and terminated.
//
// The Manager is thread-safe and can be accessed concurrently from multiple
// WebSocket handler goroutines.
package presence

import (
	"sync"
)

type User struct {
	UserID string `json:"user_id"`
	Status string `json:"status"` // e.g., "online", "idle"
}

type Manager struct {
	mu         sync.RWMutex
	workspaces map[string]map[string]User // workspaceID -> userID -> User
}

func NewManager() *Manager {
	return &Manager{
		workspaces: make(map[string]map[string]User),
	}
}

func (m *Manager) AddUser(workspaceID string, user User) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.workspaces[workspaceID]; !ok {
		m.workspaces[workspaceID] = make(map[string]User)
	}
	m.workspaces[workspaceID][user.UserID] = user
}

func (m *Manager) RemoveUser(workspaceID string, userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if workspace, ok := m.workspaces[workspaceID]; ok {
		delete(workspace, userID)
		if len(workspace) == 0 {
			delete(m.workspaces, workspaceID)
		}
	}
}

func (m *Manager) GetUsers(workspaceID string) []User {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var users []User
	if workspace, ok := m.workspaces[workspaceID]; ok {
		for _, u := range workspace {
			users = append(users, u)
		}
	}
	return users
}
