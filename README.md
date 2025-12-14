# EtherPly
> "The Heroku for Multiplayer."

EtherPly is a managed state synchronization engine that turns static apps into collaborative, real-time platforms. It provides a conflict-free, persistent, and scalable backend for real-time applications.

## Repository Structure

| Directory | Component | Description |
|---|---|---|
| [`etherply-sync-server/`](./etherply-sync-server) | **Core Engine** | Go-based WebSocket server with CRDT logic and Disk persistence. |
| [`packages/sdk-js/`](./packages/sdk-js) | **JS SDK** | Client library for React/Node.js. |
| [`packages/sdk-python/`](./packages/sdk-python) | **Python SDK** | Client library for Backend bots & IoT. |
| [`apps/docs/`](./apps/docs) | **Documentation** | Documentation site (Docusaurus). |

### Examples Library
| Example | Path | Tech Stack |
|---|---|---|
| **Text Editor** | [`examples/text-editor`](./examples/text-editor) | React, `useDocument` (The "Hello World") |
| **Kanban Board** | [`examples/kanban`](./examples/kanban) | React, DnD, Complex Object Sync |
| **Cursors** | [`examples/cursors`](./examples/cursors) | React, High-Freq Presence |
| **IoT Dashboard** | [`examples/iot`](./examples/iot) | Python (Simulator) + Next.js |
| **Voting App** | [`examples/voting`](./examples/voting) | Python (Bot) + Next.js |

## Quick Start (Run the Full Stack)

### Prerequisites
- **Go 1.23+**: Verify with `go version`
- **Node.js 18+**: Verify with `node -v`
- **Port 8080**: Must be free.

### 1. Start the Backend
```bash
cd etherply-sync-server
go mod tidy
go run main.go
```
*Expected Output*: `EtherPly Sync Server starting on port 8080`

### 2. Start the Frontend (Text Editor)
In a new terminal:
```bash
cd examples/text-editor
npm install
npm run dev
```
*Action*: Open http://localhost:3000 in two browser windows. Type in one, see it in the other.

## Troubleshooting

### "Address already in use"
**Symptom**: Server fails to start with `bind: address already in use`.
**Cause**: Another process is using port 8080.
**Fix**:
```bash
lsof -i :8080
kill -9 <PID>
```

### "Connection Refused"
**Symptom**: Frontend console shows WebSocket errors.
**Cause**: Backyard server is not running or crashed.
**Fix**: Ensure `go run main.go` is active in the backend terminal.

## Documentation Index
- [**Product Roadmap**](apps/docs/docs/roadmap.md) - Strategic direction & Pivot plan.
- [Commercial Due Diligence](apps/docs/docs/commercial_due_diligence.md) - Investment analysis.
- [Quality Audit](apps/docs/docs/quality_audit.md) - Frontend analysis.
- [Deployment Guide](apps/docs/docs/deployment.md) - Docker & Production setup.
- [Tech Debt Alert](docs/tech_debt.md) - Critical architectural warnings.

## Status
- **Metric:** Architecture Robustness / Correctness
- **Current Phase:** **PIVOT** (See [Roadmap](apps/docs/docs/roadmap.md)) - Paused feature work for core re-engineering.
