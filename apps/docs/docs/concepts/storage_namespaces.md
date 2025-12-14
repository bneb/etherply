# Storage Namespaces

## Overview

EtherPly uses **Storage Namespaces** to provide strict isolation between different types of data within the shared BadgerDB instance. This ensures that workspace data (documents), system metadata (projects), and operational metrics (metering) never collide or leak into each other.

## Key Concepts

- **Namespace**: A logical partition of the Key-Value store, implemented as a key prefix.
- **Prefix**: The actual string prepended to the key before storage (e.g., `ws:`, `sys:projects`).
- **Isolation**: Consumers (e.g., CRDT Engine) are unaware of the global namespace; they operate within their assigned scope.

## Schema Reference

| Namespace | Prefix | Description | Consumers |
| :--- | :--- | :--- | :--- |
| **Workspace Data** | `ws:{workspace_id}` | Stores Automerge documents for a specific workspace. | `internal/crdt` |
| **Project Metadata** | `sys:projects` | Stores global list of projects and their settings. | `internal/store`, `controlplane` |
| **Metering Data** | `sys:metering` | Stores usage metrics (messages sent/received) for billing. | `internal/metering` |

## Implementation Details

The `BadgerStore` (and `MemoryStore`) implementation automatically handles prefixing.

```go
// Example: Accessing Workspace Data
// The engine requests key "sync_doc" in namespace "ws:abc-123"
// Actual Key in DB: "ws:abc-123:sync_doc"
store.Get("ws:abc-123", "sync_doc")
```

### Isolation Guarantees

1.  **Prefix Separation**: Every key is required to belong to a namespace.
2.  **No Cross-Talk**: A `Get` operation requires the namespace, making it impossible to accidentally read data from another workspace or system domain without explicit intent.
