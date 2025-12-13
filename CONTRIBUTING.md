# Contributing to EtherPly

## The Standard
We are building critical infrastructure. "It works on my machine" is not acceptable. We operate with the rigor of a space agency.

## 1. The Workflow (Pull Requests)
*   **Atomic Commits:** One logical change per commit.
*   **Conventional Commits:** Use the standard: `feat:`, `fix:`, `chore:`, `docs:`.
    *   *Bad:* "Fixed the thing"
    *   *Good:* `fix(server): handle nil pointer in presence manager on disconnect`
*   **No PR without Tests:** If you touch logic, you add a test. If you touch UI, you add a screenshot.

## 2. Code Style
*   **Go:** `gofmt` is strict law. No exceptions.
*   **TypeScript:** ESLint with strict mode. No `any` types allowed (unless actively migrating legacy).
*   **Comments:** "Radical Explicitness." Explain *why*, not just *what*.

## 3. Definition of Done (DoD)
A feature is not done when the code is written. It is done when:
1.  [ ] Tests pass (CI/CD green).
2.  [ ] Documentation is updated (`docs/`).
3.  [ ] Metrics are instrumented (PostHog).
4.  [ ] Operational risks are assessed (e.g., "Will this explode database CPU?").

## 4. Reporting Issues
*   **Use the Template:** Do not just say "it's broken." Provide:
    *   Reproduction Steps
    *   Expected vs Actual Behavior
    *   Logs/Screenshots
    *   Impact Assessment (Sev1/2/3)
