# React Quickstart

Get started with EtherPly in 5 minutes using our React hooks.

## 1. Install

```bash
npm install @etherply/sdk
```

## 2. Setup Provider

Wrap your application with `EtherPlyProvider` to initialize the client.

```tsx title="src/App.tsx"
import { EtherPlyProvider } from '@etherply/sdk/react';

const config = {
  workspaceId: 'my-first-app',
  token: 'YOUR_JWT_TOKEN' // In production, fetch this from your backend
};

export default function App() {
  return (
    <EtherPlyProvider config={config}>
      <CollaborativeEditor />
    </EtherPlyProvider>
  );
}
```

## 3. Sync State

Use `useDocument` to sync state across clients.

```tsx title="src/CollaborativeEditor.tsx"
import { useDocument } from '@etherply/sdk/react';

export function CollaborativeEditor() {
  // Syncs 'content' key in real-time
  const { value, setValue, isLoaded } = useDocument<string>({
    key: 'content',
    initialValue: ''
  });

  if (!isLoaded) return <div>Loading...</div>;

  return (
    <textarea
      value={value}
      onChange={(e) => setValue(e.target.value)}
      placeholder="Type here..."
    />
  );
}
```

## 4. Add Presence

Show who else is online with `usePresence`.

```tsx title="src/PresenceBar.tsx"
import { usePresence } from '@etherply/sdk/react';

export function PresenceBar() {
  const users = usePresence();

  return (
    <div className="flex gap-2">
      {users.map(user => (
        <div key={user.userId} className="avatar">
          {user.userId.slice(0, 2)}
        </div>
      ))}
      <div>{users.length} active</div>
    </div>
  );
}
```

## Next Steps

- [API Reference](../api/sdk/index.md)
- [Authentication Guide](./authentication.md)
