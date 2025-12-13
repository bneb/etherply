# Integrate EtherPly in 5 Minutes

## Prerequisites
- A valid EtherPly API Key (get one at dashboard.etherply.com)
- A React / Next.js project

## 1. Install the SDK
```bash
npm install @etherply/sdk
```
*(For now, copy `lib/etherply-client.ts` to your project)*

## 2. Initialize the Client
```typescript
import { EtherPlyClient } from './lib/etherply-client';

const client = new EtherPlyClient({
    workspaceId: "my-first-app",
    userId: "user-123" // Replace with your auth logic
});

client.connect();
```

## 3. Listen for Data
```typescript
client.onMessage((msg) => {
    if (msg.type === 'op' && msg.payload.key === 'shared-text') {
        console.log("New value:", msg.payload.value);
    }
});
```

## 4. Send Updates
```typescript
// Call this when your user types
client.sendOperation("shared-text", "Hello World");
```

## 5. Add Presence
Use our pre-built hook (coming soon) or fetch from API:
```typescript
const users = await fetch('http://api.etherply.com/v1/presence/my-first-app').then(r => r.json());
```
