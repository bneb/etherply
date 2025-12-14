# Request for Comments (RFC)

**RFC ID:** [YYYY-MM-DD-Title]
**Status:** [DRAFT | REVIEW | APPROVED | REJECTED]
**DRI:** [Author Name]

## 1. The Problem Statement
> "A problem well stated is a problem half-solved." - Charles Kettering

Describe the **business or technical problem** we are solving. Do not describe the solution yet.
*   **Context:** Why now?
*   **Impact:** What happens if we do nothing?

## 2. Proposed Solution
Describe the architectural change.
*   **High-Level Design:** (Mermaid Diagram preferred)
*   **API Changes:** (Protobuf / JSON Schema)
*   **Database Changes:** (Schema migrations)

## 3. The "Steel Man" (Risk Assessment)
> Rigorously argue *against* your own proposal.

*   **Failure Modes:** How will this break?
*   **Operational Burden:** Does this increase on-call load?
*   **Performance:** What is the latency impact?
*   **Security:** New attack vectors?

## 4. Operational Readiness
*   [ ] Metrics (PostHog/Prometheus) defined?
*   [ ] Rollback plan documented?
*   [ ] Feature Flag defined?

## 5. Alternatives Considered
Why is this better than the other options?
