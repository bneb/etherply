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

### 1. Start the Backend
```bash
cd etherply-sync-server
go mod tidy
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
- [API Reference](docs/api_reference.md) - WebSocket and REST API documentation
- [Integration Guide](docs/integrate.md) - 5-minute quickstart for developers
- [Commercial Due Diligence](docs/commercial_due_diligence.md) - Investment readiness assessment
- [Technical Debt Alert](docs/tech_debt.md) - Known limitations and remediation plan
- [Quality Audit](docs/QUALITY_AUDIT.md) - Frontend code quality analysis

## Status
- **Metric:** Weekly Active Teams (WATs)
- **Current Phase:** Series A Preparation / MVP hardening.
