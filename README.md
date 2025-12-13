# EtherPly
> "The Heroku for Multiplayer."

EtherPly is a managed state synchronization engine that turns static apps into collaborative, real-time platforms.

## Repository Structure
This monorepo contains the following components:

| Directory | Component | Description |
|---|---|---|
| [`etherply-sync-server/`](./etherply-sync-server) | **Core Engine** | Go-based WebSocket server with CRDT logic and Disk persistence. |
| [`pkg/go-sdk/`](./pkg/go-sdk) | **Client SDK** | Go client library for connecting to the engine. |
| [`examples/demo/`](./examples/demo) | **Demo App** | Next.js 14 "Magic Moment" collaborative editor demo. |
| [`docs/`](./docs) | **Documentation** | Guides, API Reference, and Strategic Analysis. |

## Quick Start (Run the Full Stack)

### Prerequisites
- **Go 1.20+**: Verify with `go version`
- **Node.js 18+**: Verify with `node -v`

### 1. Start the Backend
```bash
cd etherply-sync-server
go mod tidy
# 2 AM Check: Ensure port 8080 is free before running
# lsof -i :8080 
go run main.go
```
*Wait for output: "EtherPly Sync Server starting on port 8080"*

### 2. Start the Frontend Demo
In a new terminal:
```bash
cd examples/demo
npm install
npm run dev
```
*Open http://localhost:3000 in two browser windows.*

## Documentation Index
- [**Product Roadmap**](docs/ROADMAP.md) - Strategic direction & Pivot plan.
- [Architecture Overview](docs/architecture.md) - System design and data flow diagrams.
- [API Reference](docs/api_reference.md) - WebSocket protocol and REST endpoints.
- [Integration Guide](docs/integrate.md) - 5-minute quickstart.
- [Commercial Due Diligence](docs/commercial_due_diligence.md) - Investment analysis.
- [Technical Debt Alert](docs/tech_debt.md) - Critical architectural warnings.
- [Quality Audit](docs/QUALITY_AUDIT.md) - Frontend analysis.

## Status
- **Metric:** Architecture Robustness / Correctness
- **Current Phase:** **PIVOT** (See [ROADMAP.md](docs/ROADMAP.md)) - Paused feature work for core re-engineering.
