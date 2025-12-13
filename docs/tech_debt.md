# Technical Debt & Architectural Failures

> [!CAUTION]
> **CRITICAL ARCHITECTURAL WARNING**
> The current system is a Proof of Concept (PoC) and must NOT be deployed to production environments requiring data durability or conflict resolution guarantees.

## Severity Levels
- **[CRITICAL]** System failure imminent at scale. Immediate remediation required.
- **[HIGH]** Significant operational risk. Remediation required before V1.0.
- **[MEDIUM]** Maintenance burden. Plan for Q3.

## 1. Persistence Layer: The "Toy Database" [CRITICAL]
*Location:* `internal/store/disk.go`

**The Failure:**
The current implementation is a naive Append-Only File (AOF) that loads **ALL** data into RAM on startup.
- **O(N) Startup Time:** Startup time scales linearly with history size.
- **O(N) Memory Usage:** RAM usage is proportional to total data ever written.
- **Single Point of Failure:** File corruption results in total data loss.
- **Global Write Lock:** `s.mu.Lock()` prevents concurrent writes, capping throughput to single-disk IOPS.

**Remediation Plan:**
- [ ] **DELETE** `internal/store/disk.go`.
- [ ] Implement `Store` interface using an embedded LSM-tree (BadgerDB/Pebble) or external DB (Postgres/Redis).
- [ ] Adopt **Write-Ahead Logging (WAL)** pattern for durability without blocking.

## 2. Synchronization Engine: LWW "Data Loss" [CRITICAL]
*Location:* `internal/crdt/engine.go`

**The Failure:**
The system uses "Last-Write-Wins" (LWW) based on timestamps.
- **Race Conditions:** Clocks are not synchronized. A user with a clock 5 seconds ahead will overwrite all other concurrent edits.
- **Intent Loss:** Simultaneous edits to the same field result in one being silently discarded.
- **Mathematically Unsound:** Does not satisfy Strong Eventual Consistency (SEC) properties for text/sequences.

**Remediation Plan:**
- [ ] Adopt a formal CRDT (Yjs/Automerge/RGA).
- [ ] Implement Vector Clocks / Lamport Timestamps for causal ordering.

## 3. Client SDK: Thread Unsafe [HIGH]
*Location:* `pkg/go-sdk/client.go`

**The Failure:**
The Go SDK client has no internal locking.
- **Panic Risk:** Concurrent calls to `SendOperation` will race on the WebSocket connection.
- **Undefined Behavior:** No guarantees on message ordering or state if sharing a client across goroutines.

**Remediation Plan:**
- [ ] Add `sync.RWMutex` around socket writes.
- [ ] Implement an operation queue for offline/reconnecting states.
