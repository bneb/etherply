---
sidebar_position: 10
title: Commercial Due Diligence
description: Assessment of Venture Viability (BCG Partner Audit)
---

# Commercial Due Diligence: EtherPly

**Date:** December 14, 2025
**To:** Kevin (Founder)
**From:** [Redacted], Partner, BCG
**Subject:** Assessment of Venture Viability

Per your request, I’ve torn apart your business case. I ignored the cleaner code (which is admittedly good) and focused purely on the money and the market.

Here is the cold, hard integration of your tech stack with market reality.

---

## 1. Executive Summary

**Verdict:** **Investable Seed Round** (Series A risk remains).

You have built a robust *technical* product that solves a hard problem (scaling real-time state). However, your *commercial* model is naive. You are trying to sell "infrastructure" in a market rapidly moving towards "platforms."

- **The Good:** The architecture (Go + NATS + Automerge) is legitimate "Scale-Up" tech. It’s not an MVP toy.
- **The Bad:** Your CAC assumptions are fantasy. In B2B DevTools, blended CAC is rarely under $1k.
- **The Ugly:** The "Real-time" feature market is becoming commoditized by Supabase, Convex, and Vercel.

---

## 2. Technical Audit (The "Product" Reality)

I had my team (me) look at the internals.

### Scalability (Moat Potential: High)
Using **NATS JetStream** for replication (`internal/replication/nats.go`) was a brilliant choice. It gives you multi-region "for free" compared to rolling your own Redis pub/sub mesh. This is defensible because it allows you to promise <100ms latency globally, which Firebase cannot do easily.

### Conflict Resolution (Moat Potential: Medium)
The `crdt` engine wrapping **Automerge** (`internal/crdt/engine.go`) allows you to draft off the industry standard. You didn't invent the algo, which is smart (less risk), but it means your IP is thin here.

### Observability (Enterprise Readiness: High)
The migration to `slog` and Prometheus metrics (`internal/metrics`) is what separates "Project" from "Product." You can actually sell this to a CTO now.

---

## 3. Commercial Critique

### Total Addressable Market (TAM)
- **The Pitch:** "Every app is becoming collaborative."
- **The Reality:** True, but most apps only need *simple* collaboration (online indicators), which they get for free from standardized frameworks.
- **Your Wedge:** You need to target **"Complex State Apps"**—Kanban boards, Design tools (Figma clones), IDEs.
- **Sizing:** The market for "General Realtime" is huge but crowded (Pusher, PubNub). The market for "State Synchronization" (Replicache, Liveblocks) is smaller (~$500M) but growing fast.

### Unit Economics (The "BCG" View)

Your `unit_economics.md` sheet is optimistic.

| Metric | Your Claim | Likely Reality | Implication |
| :--- | :--- | :--- | :--- |
| **CAC** | $350 | **$1,200+** | You cannot survive on $49/mo customers. You need to push Up-Market immediately. |
| **Churn** | ~4%/mo | **20%/mo** early on | Dev tools have high churn until deep integration. |
| **Gross Margin** | 85% | **65%** | Multi-region NATS clusters are expensive to operate reliably. |

**Strategic Pivot Required:**
Stop targeting "Startups" ($49/mo). They churn. Target "Scale-ups" ($499/mo starting). Your "Pro" tier is too cheap.

---

## 4. Defensibility (The Moat)

Why shouldn't I just use **Supabase Realtime**?
- *Supabase* does Postgres-to-Client pushing. It doesn't handle **Client-side CRDT conflict resolution** well (yet).
- **Your Moat:** "We handle the merge conflicts so you don't have to."
- **Risk:** If Supabase acquires a CRDT engine (like Yjs bindings), you are dead in the water. You need to race to "Enterprise Features" (SSO, Audit Logs, On-Prem) before they catch up.

---

## 5. The "Go / No-Go" Recommendation

**If I were a VC:**
I would give you a **Term Sheet for $2M Seed** on a $12M Post-money valuation.

**Conditions:**
1.  **Hire a Sales Founder:** You are an engineer. You need someone to grind the $20k ACV deals.
2.  **Raise Prices:** Kill the $49 plan. Make it Free or $199. Middle ground is death valley.
3.  **Positioning:** Don't sell "Realtime SDK." Sell "The Sync Engine for Pro SaaS."

Good luck, Kevin. You built a Ferrari engine. Now go build the rest of the car.
