# EtherPly Product Roadmap

**Last Updated:** 2025-12-13
**Status:** ACTIVE - PIVOT PHASE
**Focus:** Reliability, Correctness, and Scalability ("Operation Ironclad")

> [!IMPORTANT]
> **Strategic Pivot:** Following a comprehensive Commercial Due Diligence audit, verified by the Technical Architect, we have identified critical scalability limits in the current MVP architecture. We are pausing new feature development to re-engineer the core sync and persistence engine.

---

## Q1 2026: Operation Ironclad (The "Rebuild")
**Goal:** Transform EtherPly from a "naive prototype" into a "robust distributed system."
**Success Metric:** Zero data loss under partition; 10k concurrent connections with <50ms latency.

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
- [ ] **Vector Clocks:** Implement proper causal ordering for operations.
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
- [ ] **Session Affinity:** Sticky sessions at the load balancer level.

### 2. Permissions & Auth
- [x] **Granular ACLs:** "Read-only" vs "Write" access tokens.
- [x] **Webhooks:** Events for `client.connect`, `doc.update`.

### 3. Developer Console
- [x] **Dashboard:** View active connections, document storage size.
- [x] **History API:** REST endpoint to retrieve document history/snapshots.

---

## Q3 2026+: Enterprise Readiness
- [ ] **SLA Guarantees**
- [ ] **Multi-Region Replication**
- [ ] **On-Premise Deployment options**
