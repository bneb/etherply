# EtherPly
> "The Heroku for Multiplayer."

EtherPly is a managed state synchronization engine that turns static apps into collaborative, real-time platforms.

## Repository Structure
This monorepo contains the following components:

| Directory | Component | Description |
|---|---|---|
| [`etherply-sync-server/`](./etherply-sync-server) | **Core Engine** | Go-based WebSocket server with CRDT logic and Disk persistence. |
| [`packages/sdk-js/`](./packages/sdk-js) | **JS SDK** | Client library for React/Node.js. |
| [`packages/sdk-python/`](./packages/sdk-python) | **Python SDK** | Client library for Backend bots & IoT. |
| [`apps/docs/`](./apps/docs) | **Documentation** | Documentation site (docusaurus). |

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

### 1. Start the Backend
```bash
cd etherply-sync-server
go mod tidy
# Ensure port 8080 is free
go run main.go
```
*Wait for output: "EtherPly Sync Server starting on port 8080"*

### 2. Start the Frontend (Text Editor)
In a new terminal:
```bash
cd examples/text-editor
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
