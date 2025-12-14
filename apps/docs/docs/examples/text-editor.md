# Collaborative Text Editor

A "Hello World" example of real-time text synchronization.

## Features
- **Real-time Sync**: Uses `client.set()` and `useDocument()` to sync text field state.
- **Connection Status**: Visual indicator of WebSocket connection state.

## Code Highlight

```typescript
const { value, setValue } = useDocument<{ text: string }>({
  key: 'shared-text',
  initialValue: { text: '' }
});

// Sync on change
<textarea 
  value={value?.text}
  onChange={(e) => setValue({ text: e.target.value })}
/>
```

[View Source Code](https://github.com/etherply/etherply/tree/main/examples/text-editor)
