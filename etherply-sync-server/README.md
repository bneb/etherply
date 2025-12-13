# EtherPly Sync Server

## Overview
This is the core WebSocket engine for EtherPly. It handles real-time state synchronization, presence, and persistence.
It is written in Go to maximize concurrency (goroutines) and minimize latency (allocations).

## Directory Structure
- `cmd/`: Entry points (currently just `main.go` in root).
- `internal/crdt/`: The Conflict-Free Replicated Data Type engine. Currently implements Last-Write-Wins (LWW).
- `internal/store/`: Persistence layer. Uses `disk.go` (AOF) for durability.
- `internal/presence/`: Ephemeral user tracking.

## How to Run (Local Dev)
**Prerequisites:** Go 1.21+

1. **Install Dependencies:**
   ```bash
   go mod tidy
   ```
2. **Start the Server:**
   ```bash
   go run main.go
   ```
   *Expected Output:*
   ```text
   EtherPly Sync Server starting on port 8080
   [PERSISTENCE] Replayed X operations from disk.
   ```

## Configuration
- `PORT`: (Optional) Port to listen on. Defaults to `8080`.
- `etherply.aof`: (File) The Append-Only File created in the working directory for persistence.
- `ETHERPLY_JWT_SECRET`: (Required) Shared secret for signing/verifying JWTs. Server will not start without this.

## Troubleshooting (2 AM Guide)
- **"Failed to initialize persistence layer"**:
    - **Cause:** The process cannot write to `etherply.aof`.
    - **Fix:** Check permissions in the directory. Ensure no other process has the file locked.
- **"concurrent_connection_peak" metric missing**:
    - **Cause:** PostHog stub is just printing to stdout.
    - **Fix:** Check stdout logs. This is expected behavior in MVP.
