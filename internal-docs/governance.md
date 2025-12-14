---
slug: /governance
title: Governance
internal: true
---

# Governance & Ownership Registry

## Purpose
To eliminate the "diffusion of responsibility." Every component has a Direct Responsible Individual (DRI). The DRI has the final say on technical decisions and is accountable for the component's uptime.

## Domain Ownership

| Component | Scope | DRI | Bus Factor Risk |
|---|---|---|---|
| **Sync Engine** | `etherply-sync-server/` | **Kevin (Founder)** | 🔴 **CRITICAL** (Solo) |
| **SDK** | `pkg/go-sdk/` | Kevin | 🔴 High |
| **Frontend/Demo** | `examples/demo/` | Kevin | 🟡 Medium |
| **Documentation** | `docs/` | Kevin | 🟢 Low |
| **Design System** | `packages/ui` / `styles` | **OPEN** | 🔴 **CRITICAL** (Tragedy of the Commons) |

## Escalation Protocols

### Severity 1 (Data Loss / Full Outage)
*   **SLA:** 15 Minutes.
*   **Protocol:**
    1.  PagerDuty trigger.
    2.  Check `etherply.aof` integrity.
    3.  Rollback last deploy immediately. **Do not debug forward.** roll back first.

### Severity 2 (Performance Degradation)
*   **SLA:** 4 Hours.
*   **Protocol:**
    1.  Check PostHog metrics (`concurrent_connection_peak`).
    2.  If CPU > 80%, vertically scale Fly.io machine.

## Decision Making Rights
*   **Architectural Changes:** Require RFC + Approval from DRI.
*   **UI/Copy Changes:** Delegated to Product/Design leads.
*   **Hotfixes:** Any engineer can deploy to fix Sev1, but must file RCA within 24h.
