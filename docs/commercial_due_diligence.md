# Commercial Due Diligence: EtherPly
**Date:** 2025-12-13
**To:** Managing Partner, EtherPly Ventures
**From:** Antigravity Strategy Group
**Subject:** CONFIDENTIAL - Investment Memo & Risk Assessment

---

## 1. Executive Summary: The "Pass" Decision

**Recommendation: PASS / HOLD**
**Valuation Impact:** Pre-Product / Tech-Debt Discounted

EtherPly positions itself as "The Heroku for Multiplayer"—a managed state synchronization engine. While the value proposition (simplifying real-time infrastructure) is a valid "Painkiller" in the current market, the technical execution is fundamentally unserious. This is not a "platform"; it is a college hobby project.

The "engine" is a thin wrapper around a native Go map with rudimentary "Last-Write-Wins" logic, lacking actual Conflict-Free Replicated Data Type (CRDT) sophistication. The persistence layer is a local file. The "Moat" is non-existent.

If you pitch this to a VC today, you will not just be rejected; you will be remembered as the "AOF guy" who tried to raise Series A on a single-node text file database.

---

## 2. Technical Weaknesses Summary (The "Deal Killers")

Our technical audit revealed three fatal flaws that prevent this from being a scalable business asset:

### A. The "CRDT" Lie
The marketing claims "CRDT logic," but `etherply-sync-server/internal/crdt/engine.go` reveals a naive **Last-Write-Wins (LWW)** implementation.
- **Reality:** You are comparing timestamps on incoming JSON payloads and overwriting a key-value map.
- **Problem:** This handles basic overwrites but fails at complex concurrent edits (e.g., text insertion, list manipulation). Real production apps use Yjs, Automerge, or OT (Operational Transformation).
- **Business Risk:** You cannot charge enterprise rates for a `map[string]interface{}`.

### B. Persistence is a Toy
The backend uses `store.NewDiskStore("etherply.aof")`.
- **Reality:** All state is saved to a **local Append-Only File** on the server's disk.
- **Problem:** This limits you to **Vertical Scaling only** (one big server). You cannot load balance across regions because they can't share the file. If that one disk fills up or corrupts, you lose customer data.
- **Business Risk:** You have zero reliability guarantees. You are one `rm etherply.aof` away from bankruptcy.

### C. The Authentication Facade
The `auth.Middleware` is currently a stub that logs a warning: `"[AUTH] No token provided (Stub: Allowing for MVP Demo convenience)"`.
- **Reality:** Anyone can connect to any workspace if they guess the ID.
- **Business Risk:** Zero security. Unsellable to B2B.

---

## 3. Commercial Analysis (The "Hard" Questions)

### 1. Monetization & Unit Economics
*   **Model:** "Managed Infrastructure" (PaaS).
*   **Competitors:** Liveblocks, PartyKit, Replicache, Supabase Realtime, Firebase.
*   **Your COGS:** High. Because you use Go locking on a single process and local file I/O, your cost per concurrent user will spike faster than revenue.
*   **Verdict:** You are entering a "Red Ocean" (crowded market) with a rowboat. You have no cost advantage and vastly inferior features.

### 2. The "Moat" (Defensibility)
*   **Scenario:** Amazon releases "AWS AppSync Lite" or Vercel acquires Liveblocks.
*   **Result:** You are dead in 48 hours.
*   **Why?** You have no IP. You have no proprietary algorithms. The "Sync Server" code is less complex than a redis-server fork.
*   **Verdict:** **Negative Moat.** Your code is a liability, not an asset.

### 3. Market Sizing (TAM/SAM/SOM)
*   **TAM (Total Available Market):** Huge. Every app wants to be multiplayer.
*   **SAM (Serviceable Available Market):** Developers who want multiplayer but don't want to use Firebase/Supabase.
*   **SOM (Serviceable Obtainable Market):** Developers who don't know that Yjs exists. This is a shrinking market of "uninformed juniors."

### 4. Scalability
*   **Current Limit:** ~1 server.
*   **Growth Potential:** Zero without complete rewrite.
*   **Manual Ops:** High. You will be manually SSH-ing into servers to backup `.aof` files.

---

## 4. Remediation Plan (How to Fix This)

To upgrade this from "Project" to "Product," you must execute the following pivot immediately:

1.  **Adopt Real Standards:** Scrap your custom "LWW Engine." Wrap **Yjs** or **Automerge** in your Go server. Sell the *management* of those documents, not the *invention* of a bad algorithm.
2.  **Fix Persistence:** Replace `DiskStore` with **Redis** (for hot state) + **PostgreSQL/S3** (for cold storage).
3.  **Horizontal Scaling:** Use **Redis Pub/Sub** so multiple Go servers can talk to each other.
4.  **Kill the "Demo" Auth:** Integrate unopinionated JWT validation immediately.

---

**Final Word:**
You built a nice prototype for a CS 401 class. But you asked me if this is a business.
**It is not.**
Go back to the whiteboard, delete `engine.go`, and build a real platform.

*Sincerely,*
*Antigravity Strategy Group*
