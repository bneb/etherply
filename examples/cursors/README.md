# EtherPly Example: Live Cursors

> **Concept**: Demonstrates **Presence**, a high-frequency, ephemeral data layer separate from the persistent document store. Used for "Who is online" lists and live pointers.

## Architecture

Presence data is **ephemeral**. It is not stored in BadgerDB. It lives in memory (or Redis) and vanishes when the socket disconnects.

```mermaid
graph TD
    Mouse[MouseMove] -->|Throttle 30ms| Client[Client]
    Client -->|Send {x,y}| Server[Presence Hub]
    Server -->|Broadcast| Peers[Other Clients]
    Peers -->|Interpolate| Canvas[Canvas Layer]
```

### Data Structure
```ts
// Ephemeral Payload
{
  "user_id": "alice",
  "data": { "x": 100, "y": 200, "color": "#FF0000" },
  "last_seen": 1678900000
}
```

## Run It

### 1. Prerequisites
Ensure the **EtherPly Sync Server** is running (`:8080`).

### 2. Start Application
```bash
npm install
npm run dev
```

### 3. Verify
1. Open `http://localhost:3000` in multiple windows.
2. Move mouse in one.
3. See cursors glide in others.

## Troubleshooting

### Cursors imply lag
**Cause**: Lack of interpolation on the receiving end.
**Fix**: The `LiveCursors` component (used here) implements Linear Interpolation (Lerp). Ensure your frame rate is 60fps.

### "Too Many Requests" / Flooding
**Cause**: Sending every mousemove event without throttling.
**Fix**: Ensure `lodash.throttle` is applied to the event listener (default: 30ms).
