# Technical Debt Alert & Refactoring Plan

> [!WARNING]
> COMPLIANCE: REQUIRED OUTPUT FOR SERIES A TECHNICAL DUE DILIGENCE.

## Current State
The current MVP uses `store.DiskStore` with Append-Only File (AOF) persistence via `gob` encoding. This provides **MVP Durability** - data survives server restarts. However, for production scale, this single-file approach has limitations:
- Single-node bottleneck (no horizontal scaling)
- No transaction isolation between workspaces
- Recovery time increases with log size

## Critical Refactoring Steps (Next 3 Sprints)

### 1. Integrate Distributed Key-Value Store
**Target:** FoundationDB (or CockroachDB as a fallback).
**Action:** Implement the `store.StateStore` interface using the FDB Go client.
**Why:** Guarantees serializable transactions and multi-region correctness.

### 2. Implement Write-Ahead Logging (WAL)
**Target:** Kafka or Redpanda.
**Action:** Before applying to Memory/FDB, push the operation to a topic.
**Why:** Decouples ingestion from processing and ensures we can "Time Travel" / Replay state if the engine crashes.

### 3. Connection Sharding & Hashing
**Target:** Consistent Hashing Ring.
**Action:** Route `workspace_id` to specific Go pods based on a hash ring.
**Why:** To support > 100k concurrent connections, we cannot have all workspaces on one node. This ensures locality for merge operations.
