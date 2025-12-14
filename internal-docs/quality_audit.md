---
title: Quality Audit
internal: true
---

# 🕵️ Quality Gatekeeper Audit: CollaborativeEditor.tsx

**Date:** 2025-12-13
**Auditor:** Principal Design Engineer / QA Architect
**Component:** `CollaborativeEditor.tsx`

---

## 1. The Visual Forensic Report (The "What")

### The Glitch List
*   **Arbitrary Geometry (`h-96`)**: You have hardcoded the editor height to `h-96` (384px). On a modern 27" monitor (2560x1440), this looks like a postage stamp. On a mobile device with keyboard open, it might push content off-screen.
*   **Loading State UX Disaster**: The absolute overlay `bg-white/50 backdrop-blur-sm` is a nice idea, but you've placed it *inside* a container relative to the `textarea`. If the status is `IDLE` or `CONNECTING`, the user sees a "disabled" looking box.
*   **Placeholder Typos**: "Start typing to initialize workspace..." vs "Connecting...". This text shift causes a layout repaint and a visual "jump" if the font rendering matches poorly.
*   **Scrollbar Jank**: `resize-none` is good, but without `overflow-hidden` or custom scrollbar styling, the native browser scrollbar will appear/disappear as content grows, causing a layout shift (CLS) for the text content.
*   **Focus Ring Overlap**: `focus-within:ring-2` applies a ring to the *container*. If the `textarea` also has default browser focus styles (which `outline-none` tries to remove, but not essentially guaranteed across all user agent stylesheets without reset), you might get double focus indicators.

### Responsiveness Stress Test
*   **Mobile (375px)**: The 384px height is fixed. With a virtual keyboard taking up ~40-50% of the viewport, the top of the editor might be scrolled out of view, or the user is editing blindly.
*   **Ultrawide (4k)**: The `max-w` isn't constrained in this file (presumably in parent), but assuming full width, a line length of 200+ characters makes reading painful. Typography 101: 65-75 chars per line. The font is `font-mono`, complicating this further.

### The "Jank" Index
*   **Cursor Jitter**: The "Naive cursor restoration" logic uses `requestAnimationFrame` *after* a state update. This guarantees a single-frame "flash" where the cursor jumps to the end (default React behavior on value change) and then snaps back. It looks broken.
*   **Typing Lag**: You are calling `client.sendOperation` on *every single keystroke* (line 86) without debouncing. On a 120Hz display with a fast typist, you are blocking the main thread and flooding the network queue.

---

## 2. The Root Cause Analysis (The "Why")

### Code Pathology
*   **Unreachable Cleanup Logic (Line 64)**:
    ```typescript
    return () => { if (unsubscribe) unsubscribe(); };
    return () => { if (unsubscribe) unsubscribe(); }; // 💀 Dead Code
    ```
    You have two return statements in the same `useEffect`. The second one is unreachable. While harmless at runtime, it indicates sloppy copy-pasting.
*   **The "Never-Ready" State (Line 9 & 74)**:
    `const [isInitialized, setIsInitialized] = useState(false);`
    You check `!isInitialized` on line 74, but **you verify explicitly that `setIsInitialized(true)` is NEVER called anywhere in the file.** This means the initialization logic (if you intended to have one) is permanently broken or the variable is useless noise.
*   **Network Flood (Line 86)**:
    `client.sendOperation` is called synchronously in `handleChange`. This creates a 1:1 packet-to-keystroke ratio. In a persistent connection, this is a denial-of-service attack on your own backend.

### State Desynchronization
*   **Race Condition in `init`**: If an `init` message arrives *after* the user has started typing (unlikely in this specific flow but possible if connection drops and reconnects), `setContent(text)` (Line 20) clobbers the user's local draft without merging.
*   **Drifting Cursor**: The logic `input.setSelectionRange(cursorStart, cursorEnd)` preserves the *index*, not the *identity* of the insertion point. If User B inserts text at index 0 while User A is typing at index 100, User A's cursor blindly stays at 100, effectively jumping *backward* relative to their text.

### Console Prophecy
*   **Memory Leak Warning**: If `client.onMessage` returns a cleanup function, you are calling it. Good. But if `client` connection changes, you re-subscribe.

---

## 3. The Remediation Plan (The "Fix")

### The Hotfix (Immediate)
1.  **Remove Dead Code**: Delete the duplicate return in `useEffect`.
2.  **Fix Initialization**: Actually call `setIsInitialized(true)` when `init` message is received or status becomes connected.
3.  **Debounce Writes**: Wrap `client.sendOperation` in a `useDebouncedCallback` (e.g., 300ms) or similar custom logic.

### The Systemic Fix (Architecture)
1.  **Cursor Preservation Hook**: Extract cursor logic to `useCursorPreservation(ref, value)`. This hook should use `getSnapshotBeforeUpdate` logic or `useLayoutEffect` (not `requestAnimationFrame`) to measure and restore cursor *synchronously* before the browser paints.
2.  **Visual Refactor**:
    *   Change `h-96` to `min-h-[12rem] h-auto` (auto-growing textarea) or `h-full` if in a flex container.
    *   Use a `ResizeObserver` to handle layout shifts.
    *   Switch to a semantic color variable (e.g., `border-border` instead of `border-gray-200`) to support dark mode automatically.

### Test Case Generation (Gherkin)

```gherkin
Feature: Collaborative Editing Robustness

  Scenario: Typing does not cause network flood
    Given I am connected to the EtherPly session
    When I type "Hello World" rapidly (100ms per char)
    Then the client should send at most 1 operation message per 300ms
    And the final document state should be "Hello World"

  Scenario: Cursor Stability on Re-render
    Given I have text "Hello |World" (cursor at pipe)
    When a remote update changes text to "Hello Beautiful |World"
    Then my cursor should be at "Hello Beautiful |World"
    And there should be no visual jump (layout shift < 0.01)

  Scenario: Mobile Viewport Safety
    Given I am on a mobile device (375px width)
    When the virtual keyboard opens
    Then the editor should remain visible above the fold
    And the "Is Connecting" overlay should not obscure the active text area if connected
```

---

## 4. Resolution Status (2025-12-13)

> [!TIP]
> The following issues from this audit have been **RESOLVED**:

| Issue | Status | Implementation |
|-------|--------|----------------|
| Dead Code (duplicate return) | ✅ Fixed | Removed in hook extraction |
| `isInitialized` never set | ✅ Fixed | Now set in `useCollaborativeEditor.ts` on `init` message |
| Network Flood (no debounce) | ✅ Fixed | `useDebounce` hook (300ms) in `hooks/useDebounce.ts` |
| Cursor Jitter | ✅ Fixed | `useLayoutEffect` for synchronous cursor restoration |
| Hardcoded `h-96` | ✅ Fixed | Changed to `min-h-[400px]` with flex layout |

**Remaining Architectural Recommendations:**
- ResizeObserver for dynamic layouts (deferred)
- Full OT/CRDT cursor transformation (requires backend changes)

