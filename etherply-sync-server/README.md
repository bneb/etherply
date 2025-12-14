# EtherPly Sync Server

The robust, Go-based heart of the EtherPly platform. Handles WebSocket connections, CRDT conflict resolution (Automerge), and durable persistence (BadgerDB).

## Quick Start

### Option 1: Go

```bash
export ETHERPLY_JWT_SECRET="your-secret-here"
go run main.go
```

### Option 2: Docker

```bash
docker build -t etherply-sync-server .
docker run -p 8080:8080 -e ETHERPLY_JWT_SECRET="your-secret" etherply-sync-server
```

### Verify Health

```bash
curl http://localhost:8080/healthz
# {"status":"ok","timestamp":"...","uptime":"..."}

curl http://localhost:8080/readyz
# {"status":"ok","timestamp":"...","checks":{"store":"ok"}}
```

## Configuration

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `ETHERPLY_JWT_SECRET` | - | **Yes** | JWT signing secret for authentication |
| `PORT` | `8080` | No | HTTP server port |
| `BADGER_PATH` | `./badger.db` | No | Path to BadgerDB data directory |
| `SHUTDOWN_TIMEOUT_SECONDS` | `30` | No | Graceful shutdown timeout |
| `WEBHOOK_URL` | - | No | URL for webhook event delivery |

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/sync/{workspace_id}` | WS | WebSocket for real-time sync |
| `/v1/presence/{workspace_id}` | GET | List users in workspace |
| `/v1/history/{workspace_id}` | GET | Document change history |
| `/v1/stats` | GET | Server metrics |
| `/healthz` | GET | Liveness probe |
| `/readyz` | GET | Readiness probe |

## Documentation

- [SPECIFICATION.md](./SPECIFICATION.md) - User stories and data contracts
- [DEPLOYMENT.md](../docs/DEPLOYMENT.md) - Docker/Kubernetes deployment guide

## Troubleshooting

### "Connection Refused"
Server not running. Check: `ps aux | grep etherply`

### "Unauthorized" (401)
`ETHERPLY_JWT_SECRET` mismatch. Verify both client and server use identical secret.

### "Address Already in Use"
Port occupied. Find PID: `lsof -i :8080` then `kill -9 <PID>`

### Persistence Errors
Check `BADGER_PATH` is writable and has disk space.

