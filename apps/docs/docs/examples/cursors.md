# Multiplayer Cursors

Demonstrating high-performance presence tracking.

## Features
- **High Frequency**: Handles 30-60fps updates.
- **Throttling**: Uses `requestAnimationFrame` to prevent network flooding.
- **Presence**: Simulates ephemeral user state.

## Code Highlight

```typescript
// Throttle updates to ~30ms
if (now - lastUpdate < 33) return;

client.set(`presence:${myId}`, { x, y });
```

[View Source Code](https://github.com/etherply/etherply/tree/main/examples/cursors)
