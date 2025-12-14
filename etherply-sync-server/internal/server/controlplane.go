package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/store"
)

// HandleListProjects returns all projects.
func (h *Handler) HandleListProjects(w http.ResponseWriter, r *http.Request) {
	// In a real app, we would filter by the authenticated user's Org ID.
	// For this binary, we assume single-tenant or admin access.

	// Check type assertion for the extended store interface
	// In a real generic codebase, we'd use a better interface composition.
	// Assuming h.store implements ListProjects (added in previous step via BadgerStore)

	projects, err := h.store.(*store.BadgerStore).ListProjects()
	if err != nil {
		http.Error(w, "Failed to list projects", http.StatusInternalServerError)
		return
	}

	// Defensive: Return empty list, not null
	if projects == nil {
		projects = []store.Project{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// HandleCreateProject creates a new project.
func (h *Handler) HandleCreateProject(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name   string `json:"name"`
		Region string `json:"region"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Region == "" {
		http.Error(w, "Name and Region are required", http.StatusBadRequest)
		return
	}

	p := store.Project{
		ID:        fmt.Sprintf("prj_%d", time.Now().UnixNano()), // Simple ID gen
		Name:      req.Name,
		Region:    req.Region,
		CreatedAt: time.Now(),
	}

	if err := h.store.(*store.BadgerStore).SaveProject(p); err != nil {
		http.Error(w, "Failed to save project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// HandleGetPlans returns available billing plans.
func (h *Handler) HandleGetPlans(w http.ResponseWriter, r *http.Request) {
	plans := []map[string]interface{}{
		{"id": "free", "name": "Hobby", "price": 0, "limits": "1k connections"},
		{"id": "pro", "name": "Pro", "price": 499, "limits": "100k connections"},
		{"id": "enterprise", "name": "Enterprise", "price": "Custom", "limits": "Unlimited"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plans)
}
