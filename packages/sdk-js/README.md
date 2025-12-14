# @etherply/sdk

The official JavaScript/TypeScript SDK for [EtherPly](https://etherply.com) — real-time sync infrastructure for collaborative apps.

[![npm version](https://img.shields.io/npm/v/@etherply/sdk)](https://www.npmjs.com/package/@etherply/sdk)
[![TypeScript](https://img.shields.io/badge/TypeScript-Ready-blue)](https://www.typescriptlang.org/)
[![License](https://img.shields.io/npm/l/@etherply/sdk)](./LICENSE)

## Features

- 🚀 **5-minute integration** — Add real-time sync to any app
- ⚛️ **React hooks** — First-class React support with `useEtherPly` and `useDocument`
- 🔄 **Automatic reconnection** — Exponential backoff, operation queueing
- 📦 **Tiny bundle** — Tree-shakeable, zero dependencies
- 🔒 **Type-safe** — Full TypeScript support with comprehensive types
- 🌐 **Works everywhere** — Browser, Node.js, React Native*

## Installation

```bash
npm install @etherply/sdk
# or
yarn add @etherply/sdk
# or
pnpm add @etherply/sdk
```

## Quick Start

### Vanilla JavaScript/TypeScript

```typescript
import { EtherPlyClient } from '@etherply/sdk';

// 1. Create a client
const client = new EtherPlyClient({
  workspaceId: 'my-workspace',
  token: 'your-jwt-token'  // Get one at dashboard.etherply.com
});

// 2. Connect
await client.connect();

// 3. Listen for updates
client.onMessage((msg) => {
  if (msg.type === 'init') {
    console.log('Initial state:', msg.data);
  }
  if (msg.type === 'op') {
    console.log('Update:', msg.payload.key, '=', msg.payload.value);
  }
});

// 4. Send updates
client.set('greeting', 'Hello, world!');
client.set('counter', 42);
client.set('user', { name: 'Alice', color: 'blue' });
```

### React

```tsx
import { useEtherPly } from '@etherply/sdk/react';

function CollaborativeEditor() {
  const { state, set, status } = useEtherPly({
    workspaceId: 'my-doc',
    token: 'your-jwt-token'
  });

  return (
    <div>
      <div>Status: {status === 'CONNECTED' ? '🟢' : '🔴'} {status}</div>
      <textarea
        value={(state.content as string) || ''}
        onChange={(e) => set('content', e.target.value)}
        placeholder="Start typing..."
      />
    </div>
  );
}
```

### React with Context (Recommended for larger apps)

```tsx
import { EtherPlyProvider, useDocument } from '@etherply/sdk/react';

// Wrap your app once
function App() {
  return (
    <EtherPlyProvider
      config={{
        workspaceId: 'my-app',
        token: 'your-jwt-token'
      }}
    >
      <Counter />
      <Title />
    </EtherPlyProvider>
  );
}

// Use in any component
function Counter() {
  const { value, setValue } = useDocument<number>({
    key: 'counter',
    initialValue: 0
  });

  return (
    <button onClick={() => setValue((value || 0) + 1)}>
      Count: {value}
    </button>
  );
}

function Title() {
  const { value, setValue } = useDocument<string>({ key: 'title' });

  return (
    <input
      value={value || ''}
      onChange={(e) => setValue(e.target.value)}
      placeholder="Document title..."
    />
  );
}
```

## API Reference

### `EtherPlyClient`

The core client class for connecting to EtherPly.

#### Constructor

```typescript
new EtherPlyClient(config: EtherPlyConfig)
```

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `workspaceId` | `string` | **required** | Workspace/room ID |
| `token` | `string` | **required** | JWT authentication token |
| `userId` | `string` | auto-generated | User identifier for presence |
| `serverUrl` | `string` | `wss://api.etherply.com` | Server URL |
| `autoReconnect` | `boolean` | `true` | Auto-reconnect on disconnect |
| `maxReconnectAttempts` | `number` | `10` | Max reconnection attempts |
| `reconnectBaseDelay` | `number` | `1000` | Base delay (ms) for backoff |
| `debug` | `boolean` | `false` | Enable console logging |

#### Methods

| Method | Description |
|--------|-------------|
| `connect(): Promise<void>` | Connect to the server |
| `disconnect(): void` | Disconnect from the server |
| `set(key, value): void` | Set a value (alias for `sendOperation`) |
| `get<T>(key): T \| undefined` | Get a value from local state |
| `getState(): Record<string, unknown>` | Get entire local state |
| `sendOperation(key, value): void` | Send an operation to the server |
| `onMessage(handler): () => void` | Subscribe to messages (returns unsubscribe fn) |
| `onStatusChange(handler): () => void` | Subscribe to status changes |
| `getStatus(): ConnectionStatus` | Get current connection status |

#### Connection Status

| Status | Description |
|--------|-------------|
| `IDLE` | Initial state before `connect()` |
| `CONNECTING` | Connection in progress |
| `CONNECTED` | Successfully connected |
| `DISCONNECTED` | Connection closed |
| `RECONNECTING` | Attempting to reconnect |
| `ERROR` | Fatal error, will not reconnect |

### React Hooks

#### `useEtherPly(options)`

All-in-one hook for connecting and syncing.

```typescript
const {
  state,        // Record<string, unknown> - current state
  set,          // (key, value) => void - update a value
  get,          // <T>(key) => T | undefined - get a value
  status,       // ConnectionStatus
  isConnected,  // boolean
  client,       // EtherPlyClient instance
  connect,      // () => Promise<void>
  disconnect,   // () => void
} = useEtherPly({
  workspaceId: 'my-workspace',
  token: 'jwt-token',
  onConnect: () => console.log('Connected!'),
  onDisconnect: () => console.log('Disconnected'),
  onError: (err) => console.error(err),
});
```

#### `useDocument<T>(options)`

Sync a single key with useState-like ergonomics. Must be used within `EtherPlyProvider`.

```typescript
const {
  value,     // T | undefined
  setValue,  // (newValue: T) => void
  isLoaded,  // boolean
} = useDocument<number>({
  key: 'counter',
  initialValue: 0,
});
```

#### `EtherPlyProvider`

Context provider for sharing a client across components.

```tsx
<EtherPlyProvider
  config={{ workspaceId: 'x', token: 'y' }}
  autoConnect={true}  // default: true
>
  {children}
</EtherPlyProvider>
```

#### `useEtherPlyContext()`

Access the client from context (for advanced usage).

```typescript
const client = useEtherPlyContext();
client.set('key', 'value');
```

## Common Patterns

### Optimistic Updates

Updates are applied locally immediately for instant feedback:

```typescript
const { state, set } = useEtherPly({ ... });

const handleChange = (newValue) => {
  set('text', newValue);  // Instant local update + server sync
};
```

### Offline Support

Operations are queued when disconnected and sent when reconnected:

```typescript
const client = new EtherPlyClient({
  workspaceId: 'my-workspace',
  token: 'jwt-token',
  autoReconnect: true,
});

// These will be queued and sent when connected
client.set('key1', 'value1');
client.set('key2', 'value2');

await client.connect();  // Queued ops are sent
```

### Multiple Workspaces

Create separate clients for each workspace:

```typescript
const doc1Client = new EtherPlyClient({
  workspaceId: 'document-1',
  token: 'jwt-token'
});

const doc2Client = new EtherPlyClient({
  workspaceId: 'document-2',
  token: 'jwt-token'
});
```

### Server-Side (Node.js)

```typescript
import { EtherPlyClient } from '@etherply/sdk';

const client = new EtherPlyClient({
  workspaceId: 'my-workspace',
  token: process.env.ETHERPLY_TOKEN!,
  serverUrl: 'wss://api.etherply.com'  // or your self-hosted URL
});

await client.connect();
client.set('serverTimestamp', Date.now());
```

## TypeScript

The SDK is written in TypeScript and provides full type definitions:

```typescript
import type {
  EtherPlyConfig,
  ConnectionStatus,
  EtherPlyMessage,
  Operation,
  InitMessage,
  OperationMessage,
} from '@etherply/sdk';
```

## Local Development

When developing locally with the EtherPly server:

```typescript
const client = new EtherPlyClient({
  workspaceId: 'my-workspace',
  token: 'dev-token',
  serverUrl: 'ws://localhost:8080',  // Local server
  debug: true,  // Enable console logging
});
```

## Troubleshooting

### "WebSocket connection failed"

1. Check that the server is running
2. Verify the `serverUrl` is correct
3. Check that your `token` is valid
4. Ensure your firewall allows WebSocket connections

### "Token is required"

The `token` parameter is required. Get a token from [dashboard.etherply.com](https://dashboard.etherply.com) or your authentication provider.

### Connection keeps dropping

- Check your network stability
- Increase `maxReconnectAttempts` if needed
- Enable `debug: true` to see connection logs

## License

MIT © EtherPly Team

## Links

- [Documentation](https://docs.etherply.com)
- [API Reference](https://docs.etherply.com/api)
- [Examples](https://github.com/bneb/etherply/tree/main/examples)
- [Discord Community](https://discord.gg/etherply)
