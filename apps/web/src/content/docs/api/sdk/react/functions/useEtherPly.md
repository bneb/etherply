[**@etherply/sdk**](../../README.md)

***

# Function: useEtherPly()

> **useEtherPly**(`options`): [`UseEtherPlyReturn`](../interfaces/UseEtherPlyReturn.md)

Defined in: [src/react/useEtherPly.tsx:110](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/react/useEtherPly.tsx#L110)

React hook for real-time synchronization with EtherPly.

This hook creates and manages an EtherPly client, provides reactive
state, and handles connection lifecycle automatically.

## Parameters

| Parameter | Type | Description |
| ------ | ------ | ------ |
| `options` | [`UseEtherPlyOptions`](../interfaces/UseEtherPlyOptions.md) | Configuration and callbacks |

## Returns

[`UseEtherPlyReturn`](../interfaces/UseEtherPlyReturn.md)

Object with state, setters, and connection status

## Examples

```tsx
function CollaborativeEditor() {
  const { state, set, status } = useEtherPly({
    workspaceId: 'my-doc',
    token: 'jwt-token'
  });
  
  return (
    <div>
      <p>Status: {status}</p>
      <textarea
        value={state.content as string || ''}
        onChange={(e) => set('content', e.target.value)}
      />
    </div>
  );
}
```

```tsx
const { state, set } = useEtherPly({
  workspaceId: 'my-doc',
  token: 'jwt-token',
  onConnect: () => console.log('Connected!'),
  onDisconnect: () => console.log('Disconnected'),
  onError: (err) => console.error('Error:', err)
});
```
