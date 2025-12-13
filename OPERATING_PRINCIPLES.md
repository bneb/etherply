# EtherPly Operating Principles
> "Structure is not bureaucracy; structure is the enabling constraint for velocity."

## 1. Technocratic Efficiency (The "Smart Nation" Standard)
We do not waste time. Time is the only non-renewable asset of a startup.
*   **Written over Spoken:** If it isn't written down, it didn't happen. No decisions are made in hallway chats. Write a memo.
*   **Asynchronous Default:** We prioritize deep work. Do not interrupt engineers with Slack messages unless the server is on fire.
*   **The "Two-Way Door" Rule:** Irreversible decisions (Picking a Database, Pricing Model) require a 6-page memo. Reversible decisions (CSS change, Feature flag) should be made instantly by the DRI (Directly Responsible Individual).

## 2. Safety & Alignment (The "Anthropic" Standard)
We build powerful synchronization infrastructure. It must be reliable.
*   **Defense in Depth:** "Usage Error" is a failure of system design, not the user. If a junior dev can break prod, the senior dev failed to build guardrails.
*   **Post-Mortems are Mandatory:** Every outage > 1 minute requires a Root Cause Analysis (RCA). We do not blame humans; we fix the process.
*   **Steel-Man the Counterargument:** When proposing a change, you must rigorously explain why it might *fail*. Intellectual honesty is paramount.

## 3. Completed Staff Work
*   **Don't open tickets with "Thoughts?"; open PRs with solutions.**
*   The job of the leader is to review and approve, not to do your homework. Present a recommendation, the data that supports it, and the risks involved.

## 4. The "Bus Factor" Protocol
*   **No Silos:** If only one person understands a system (e.g., the CRDT Engine), that system is an operational risk.
*   **Immediate Documentation:** Documentation is not a "nice to have" for the end of the sprint. It is part of the Definition of Done.
