# Integrate EtherPly in 5 Minutes

> [!TIP]
> **SDK Status:** The `@etherply/sdk` package is ready for use!
> 
> ```bash
> npm install @etherply/sdk
> ```

## Prerequisites
- A valid EtherPly API Key (get one at dashboard.etherply.com)
- A React / Next.js / Vue / Svelte project

## 1. Install the SDK

```bash
npm install @etherply/sdk
```

## 2. Initialize the Client

```typescript
import { EtherPlyClient } from '@etherply/sdk';

const client = new EtherPlyClient({
    workspaceId: "my-first-app",
    userId: "user-123",  // Replace with your auth logic
    token: "your-jwt-token"  // Required in production
});

client.connect();
```

## 3. Listen for Data (React)

Use the built-in hooks for easiest integration:

```tsx
import { useEtherPly, useDocument } from '@etherply/sdk/react';

function MyComponent() {
  const { isConnected } = useEtherPly({
    workspaceId: 'my-workspace',
    token: 'jwt'
  });
  
  const { value, setValue } = useDocument({ 
    key: 'shared-text',
    initialValue: '' 
  });
  
  if (!isConnected) return <div>Connecting...</div>;
  
  return (
    <input 
      value={value} 
      onChange={(e) => setValue(e.target.value)} 
    />
  );
}
```

---

## Next Steps

| Resource | Description |
|----------|-------------|
| [API Reference](./api/http-api.md) | Full WebSocket and REST API docs |
| [Architecture](./concepts/architecture.md) | How EtherPly works under the hood |
| [Demo App](https://github.com/bneb/etherply/tree/main/examples/demo) | Full Next.js collaborative editor |
| [Go SDK](https://github.com/bneb/etherply/tree/main/pkg/go-sdk) | For backend-to-backend sync |
