# EtherPly Sync Server

The robust, Go-based heart of the EtherPly platform. Handles WebSocket connections, CRDT conflict resolution, and persistence.

## Quick Start

### 1. Build and Run
```bash
# From this directory
go mod tidy
go run main.go
```

### 2. Verify Health
Run the following command to check if the server is accepting connections:

```bash
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/v1/presence/health-check
```
**Expected Output:** `200`

## Configuration
Configuration is handled via Environment Variables.

| Variable | Description | Default | Required? |
|---|---|---|---|
| `PORT` | The HTTP port to listen on. | `8080` | No |
| `ETHERPLY_JWT_SECRET` | Secret key for validating auth tokens. | **NONE** | **YES (CRITICAL)** |
| `STORAGE_PATH` | Path to the AOF persistence file. | `etherply.aof` | No |

## Specification (Source of Truth)
For the definitive guide on User Stories, State Machines, and Data Contracts, refer to [SPECIFICATION.md](./SPECIFICATION.md).


## Troubleshooting (The "2 AM" Guide)

### "Connection Refused"
**Symptoms:** `dial tcp [::1]:8080: connect: connection refused`
**Diagnosis:** The server process is not running.
**Fix:**
1. Check if the process crashed: `ps aux | grep main.go`
2. Restart it: `go run main.go`

### "Unauthorized" / 401 Errors
**Symptoms:** Client receives `401 Unauthorized`.
**Diagnosis:** The `ETHERPLY_JWT_SECRET` env var doesn't match the token signature.
**Fix:**
1. Check the server's secret: `echo $ETHERPLY_JWT_SECRET`
2. Ensure the client is signing with the *exact* same string.

### "Panic: Address already in use"
**Symptoms:** `listen tcp :8080: bind: address already in use`
**Diagnosis:** Another instance is running, or a zombie process is holding the port.
**Fix:**
1. Find the PID: `lsof -i :8080`
2. Kill it: `kill -9 <PID>`
