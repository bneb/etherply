# EtherPly Product Roadmap

**Last Updated:** 2025-12-14  
**Status:** 🚀 Q4 2026 - "Operation Stripe" (Developer Experience)  
**Focus:** SDK, Documentation, Customer Validation

> :::note
> **Infrastructure Complete.** Q1-Q3 2026 ("Operation Ironclad") delivered production-grade infrastructure: Automerge CRDT, BadgerDB persistence, NATS replication, and on-premise options.
>
> **Now: Go-To-Market.** Q4 2026 ("Operation Stripe") focuses on developer experience—the Stripe playbook. Ship the npm SDK, world-class docs, and get our first 10 paying customers.
> :::

---

## Q1 2026: Operation Ironclad (The "Rebuild")
**Goal:** Transform EtherPly from a "naive prototype" into a "robust distributed system."
**Success Metric:** Zero data loss under partition; 10k concurrent connections with &lt;50ms latency.

### 1. Persistence Layer Overhaul (Priority: CRITICAL)
*Current State:* Append-Only File (AOF) loading entire dataset into RAM.
*Target State:* Durable, on-disk, ACID-compliant storage.
- [x] **Research & Select DB:** Evaluate embedded (SQLite/Badger) vs. external (Postgres/FoundationDB) for the core log.
- [x] **Implement WAL:** Replace `disk.go` with a Write-Ahead Log pattern tied to the new storage engine.
- [x] **Remove Global Lock:** Eliminate `s.mu.Lock()` on writes to allow concurrent processing.

### 2. Synchronization Engine Upgrade (Priority: CRITICAL)
*Current State:* Last-Write-Wins (LWW) logic (lossy conflict resolution).
*Target State:* Mathematically correct convergence.
- [x] **Integrate Real CRDTs:** Replace custom LWW logic with `Yjs` (via Y-CRDT port) or `Automerge` core.
- [x] **Vector Clocks:** Implement proper causal ordering for operations.
- [x] **Conflict Test Suite:** Automated fuzz testing for concurrent edits.

### 3. SDK Hardening
*Current State:* Not thread-safe, brittle connection handling.
*Target State:* Production-grade Go SDK.
- [x] **Thread Safety:** comprehensive locking/channel usage for `Client`.
- [x] **Resiliency:** robust reconnection logic, embrace exponential backoff.
- [x] **Queueing:** Operation buffering when offline.

---

## Q2 2026: Scale & Developer Experience
**Goal:** Enable "Stripe-like" ease of use and horizontal scalability.

### 1. Horizontal Scalability / Clustering
- [x] **Stateless Servers:** Decouple compute (WebSocket logic) from storage.
- [x] **Pub/Sub Layer:** Implement Redis/NATS for cross-server message propogation.
- [x] **Session Affinity:** Sticky sessions at the load balancer level.

### 2. Permissions & Auth
- [x] **Granular ACLs:** "Read-only" vs "Write" access tokens.
- [x] **Webhooks:** Events for `client.connect`, `doc.update`.

### 3. Developer Console (Priority: CRITICAL)
- [x] **Dashboard:** View active connections, document storage size.
- [x] **History API:** REST endpoint to retrieve document history/snapshots.

---

## Q3 2026+: Enterprise Readiness
- [x] **SLA Guarantees**
- [x] **Multi-Region Replication**
- [x] **On-Premise Deployment options**

---

## Q4 2026: Operation Stripe (The "Developer Experience" Push)
**Goal:** Best-in-class Developer Experience & Trust. Win on "It Just Works."
**Success Metric:** 50 production apps using EtherPly. 10 paying customers on Startup tier.
**Philosophy:** We are not just building a library; we are building a *database*. Developers trust Postgres with their data; they must trust EtherPly with their real-time state.

> :::important
> This phase directly addresses the Commercial Due Diligence feedback: "No JS SDK = no market" and "Get 3 paying customers."
> :::

### 1. SDK & Distribution (Priority: CRITICAL)
*Current State:* Go SDK only. TypeScript client exists but is not published.  
*Target State:* `npm install @etherply/sdk` works out of the box.

- [ ] **Publish `@etherply/sdk`:** Extract, package, and publish TypeScript SDK to npm
- [ ] **TypeScript Types:** Ship with 100% type coverage and JSDoc
- [ ] **React Hooks:** `useEtherPly()`, `usePresence()`, `useDocument()` 
- [ ] **Framework Examples:** Next.js 14, Remix, SvelteKit, Vue 3

### 2. Documentation Excellence (Priority: CRITICAL)
*Current State:* Markdown files in `/docs`. No interactive docs site.  
*Target State:* `docs.etherply.com` - Stripe/Tailwind quality.

- [ ] **Documentation Site:** Docusaurus or Mintlify hosted site
- [ ] **Interactive Quick Start:** 5-minute guided setup
- [ ] **API Reference:** Auto-generated from code with examples
- [ ] **Concepts Guide:** CRDT, LWW, sync strategies explained visually
- [ ] **Troubleshooting Guide:** Common errors + solutions
- [ ] **"5 Minute Quickstart" Video:** YouTube + landing page embed

### 3. Example Applications (Priority: HIGH)
*Current State:* 1 demo app (collaborative editor).  
*Target State:* 10+ production-grade, copy-paste examples.

- [ ] **Collaborative Text Editor** (current demo, polish)
- [ ] **Real-time Kanban Board**
- [ ] **Multiplayer Cursor Sharing** (Figma-style)
- [ ] **Live Poll/Voting App**
- [ ] **Collaborative Whiteboard**
- [ ] **Real-time Form Builder**
- [ ] **Chat Room with Presence**
- [ ] **Turn-based Game State** (Tic-Tac-Toe)
- [ ] **Inventory Sync** (e-commerce)
- [ ] **IoT Device Dashboard**

### 4. Customer Validation (Priority: CRITICAL)
*Current State:* Zero paying customers. Zero case studies.  
*Target State:* 3+ production apps, 3+ testimonials.

- [ ] **Beta Signup Page:** Collect emails, create waitlist
- [ ] **Onboard First 3 Beta Users:** Hands-on support
- [ ] **Case Study #1:** Written testimonial + logo for site
- [ ] **Case Study #2:** Video testimonial
- [ ] **Case Study #3:** Technical deep-dive blog post
- [ ] **Unit Economics Spreadsheet:** CAC, LTV, margin analysis

### 5. Observability & Ops (Priority: HIGH)
*Current State:* No metrics, basic logging.  
*Target State:* Production observability.

- [ ] **Prometheus `/metrics` Endpoint:** `etherply_connections_active`, `etherply_operations_total`, `etherply_sync_latency_ms`
- [ ] **Structured Logging:** Migrate from `log` to `slog`
- [ ] **Rate Limiting:** Token bucket per-client to prevent abuse
- [ ] **Health Check Endpoint:** `/health` and `/ready` for K8s probes
- [ ] **Grafana Dashboard Template:** Ready-to-import JSON

### 6. Pricing & Monetization
*Current State:* $29/mo (underpriced per DD report).
*Target State:* Value-based pricing aligned with market.

- [ ] **Hobby Tier:** $0/dev (Free forever for non-commercial)
- [ ] **Scale-Up Tier:** **$499/mo** (Starting Price, Volume-based)
- [ ] **Enterprise:** **Contact Sales** (VPC, SSO, Audit Logs, custom SLA)
- [ ] **Stripe Integration:** Self-serve checkout for Scale-Up

