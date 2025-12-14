# Integrate EtherPly in 5 Minutes

> [!IMPORTANT]
> **SDK Status:** The `@etherply/sdk` npm package is coming in Q4 2026. For now, copy the TypeScript client from the demo app. See [Operation Stripe Roadmap](./ROADMAP.md#q4-2026-operation-stripe-the-developer-experience-push).

## Prerequisites
- A valid EtherPly API Key (get one at dashboard.etherply.com)
- A React / Next.js / Vue / Svelte project

## 1. Install the SDK

**Coming Soon:**
```bash
npm install @etherply/sdk
```

**Current (copy from demo):**
```bash
# Copy the client to your project
cp examples/demo/lib/etherply-client.ts ./lib/
```

## 2. Initialize the Client
```typescript
import { EtherPlyClient } from './lib/etherply-client';

const client = new EtherPlyClient({
    workspaceId: "my-first-app",
    userId: "user-123",  // Replace with your auth logic
    token: "your-jwt-token"  // Required in production
});

client.connect();
```

## 3. Listen for Data
```typescript
// In React, use with useEffect
useEffect(() => {
    const cleanup = client.onMessage((msg) => {
        if (msg.type === 'op' && msg.payload.key === 'shared-text') {
            setSharedText(msg.payload.value);
        }
    });
    return cleanup;
}, []);
```

## 4. Send Updates
```typescript
// Call this when your user types
function handleInput(value: string) {
    client.sendOperation("shared-text", value);
}
```

## 5. Add Presence

Show who's currently connected:
```typescript
// REST API
const response = await fetch('http://api.etherply.com/v1/presence/my-first-app', {
    headers: { 'Authorization': `Bearer ${token}` }
});
const users = await response.json();
// users: [{ user_id: "user-123", status: "online" }, ...]
```

## 6. Handle Connection Status
```typescript
client.onStatusChange((status) => {
    // status: 'IDLE' | 'CONNECTING' | 'CONNECTED' | 'DISCONNECTED' | 'ERROR'
    console.log('Connection status:', status);
});
```

---

## Full Example (React)

```tsx
import { useState, useEffect } from 'react';
import { EtherPlyClient } from './lib/etherply-client';

const client = new EtherPlyClient({
    workspaceId: "demo-workspace",
    userId: "user-" + Math.random().toString(36).slice(2),
    token: "your-jwt-token"
});

export function CollaborativeEditor() {
    const [text, setText] = useState("");
    const [connected, setConnected] = useState(false);

    useEffect(() => {
        client.connect();
        
        const msgCleanup = client.onMessage((msg) => {
            if (msg.type === 'init') {
                setText(msg.data?.['document'] || "");
            }
            if (msg.type === 'op' && msg.payload.key === 'document') {
                setText(msg.payload.value);
            }
        });

        const statusCleanup = client.onStatusChange((status) => {
            setConnected(status === 'CONNECTED');
        });

        return () => {
            msgCleanup();
            statusCleanup();
        };
    }, []);

    const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        const newValue = e.target.value;
        setText(newValue);
        client.sendOperation("document", newValue);
    };

    return (
        <div>
            <div>Status: {connected ? '🟢 Connected' : '🔴 Disconnected'}</div>
            <textarea 
                value={text} 
                onChange={handleChange} 
                placeholder="Start typing..." 
            />
        </div>
    );
}
```

---

## Next Steps

| Resource | Description |
|----------|-------------|
| [API Reference](./api_reference.md) | Full WebSocket and REST API docs |
| [Architecture](./architecture.md) | How EtherPly works under the hood |
| [Demo App](../examples/demo) | Full Next.js collaborative editor |
| [Go SDK](../pkg/go-sdk/README.md) | For backend-to-backend sync |
