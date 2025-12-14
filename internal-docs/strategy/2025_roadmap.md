---
title: 2025 Strategy
internal: true
---

# EtherPly 2025 Strategic Roadmap (Founding CPO)

**CONFIDENTIAL: INTERNAL ENGINEERING & LEADERSHIP ONLY**
**Last Updated:** December 14, 2025
**Owner:** Founding CPO / Technical Co-Founder

---

## 1. The Investor Memorandum (Strategic View)

### The Hook: "Sovereign Realtime Infrastructure"
We are not building another Firebase clone. We are building **The "Figma Multiplayer" Engine** that enterprises can deploy in their own VPC.
*   **For Seed Investors:** We solve the "Hard Problem" (State Convergence) with a binary you can drop into Kubernetes. We are the "Postgres for Realtime".
*   **For Acquihire (Stripe/Vercel/Cloudflare):** We have solved multi-region active-active state replication on top of standard protocols (NATS). We are a plug-and-play distributed systems team.

### The Moat: Mathematical Consistency at the Edge
*   **Tech Lock-in (Series A):** By owning the CRDT merge logic in the binary (`etherply-sync-server`), we become the authoritative source of truth for the customer's collaboration data. High switching costs.
*   **Talent Advantage (Acquihire):** Our `NATSReplicator` implementation proves we understand distributed logs, not just REST APIs.

### The Ask: Operation "Real Metrics"
To achieve the **Seed Milestone ($2M raise)**, we must prove we are not vaporware. We must move from "Mocked UI" to "Measured Traffic".

---

## 2. The Product Backlog (Tactical View)

**Metric:** Equity Value Delta (EVD).
**Threshold:** > 6.0 required to build.

| Feature | Seed (Validation) | Series A (Growth) | Acquihire (Tech) | Bootstrap (Revenue) | **Aggregated EVD** | Decision |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **Usage Metering Middleware** | 8 | 10 | 7 | 10 | **8.75** | **BUILD NOW** |
| **Real Cross-Region Demo** | 10 | 6 | 8 | 4 | **7.0** | **BUILD NOW** |
| **Remove `mocks.ts` Facade** | 9 | 10 | 10 | 5 | **8.5** | **BUILD NOW** |
| **Stripe Billing Integration** | 3 | 8 | 2 | 9 | 5.5 | **CUT** |
| **Fancy Dashboard UI** | 4 | 5 | 2 | 3 | 3.5 | **CUT** |
| **BadgerDB Namespacing** | 7 | 9 | 8 | 7 | **7.75** | **BUILD NOW** |

### The "Ignore" List (Runway Preservation)
1.  **Stripe Billing UI:** We will **not** build a self-serve checkout yet. It is high effort and low signal for a technical seed. We will use "Contact Sales" and manual invoicing based on our Metering logs.
2.  **Team Permissions (RBAC):** We will **not** build complex role management. A single API Key per Project is sufficient for the Seed stage. Focus on the *Data Plane*, not the *Control Plane*.
3.  **Client SDK Sugar:** We will **not** build React hooks for every edge case. We expose the raw WebSocket/CRDT stream. Let developers build the sugar.

---

## 3. The Technical Health Check (Execution View)

### 🚨 DEBT ALERT: The "Vaporware" Facade (`mocks.ts`)
*   **Severity:** **CRITICAL (Acquihire Killer)**
*   **Analysis:** The current `apps/web/src/lib/mocks.ts` is a liability. It makes the product look complete while the backend logic is disconnected. An investor doing code diligence will see this and walk away immediately, assuming *everything* is fake.
*   **Action:** Delete `mocks.ts`. Connect the UI to `etherply-sync-server` via a real API, even if the API is simple.

### 🏗 Scale Prep: The Revenue Engine (Metering)
*   **Requirement:** To be a "Bootstrap" capable business, we must know who to invoice.
*   **Implementation:**
    *   **NO:** Do not implement billing logic in the critical path (latency penalty).
    *   **YES:** Implement *async metering*. The `etherply-sync-server` must emit logs/metrics (e.g., `workspace_id=123 messages_in=500 storage_bytes=4096`) to a collector (Prometheus or NATS topic) that we can aggregate later for invoicing.
    *   **Why:** This proves "Unit Economics" awareness to Series A investors without blocking the "Speed" pitch to developers.

---

**Next Steps for Engineering:**
1.  **Instrument `etherply-sync-server`** with usage metrics (Connections, Messages, Storage).
2.  **Delete `apps/web/src/lib/mocks.ts`** and replace with real `fetch()` calls to the server.
3.  **Hardcode** the "Free Tier" limits in the server middleware until we have real plans.
