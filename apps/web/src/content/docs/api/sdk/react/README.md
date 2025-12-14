[**@etherply/sdk**](../README.md)

***

# react

React hooks for EtherPly

## Example

```tsx
import { useEtherPly, usePresence } from '@etherply/sdk/react';

function App() {
  const { state, set, status } = useEtherPly({
    workspaceId: 'my-workspace',
    token: 'jwt-token'
  });
  
  const users = usePresence();
  
  return (
    <div>
      <p>Status: {status}</p>
      <p>Title: {state.title}</p>
      <input onChange={(e) => set('title', e.target.value)} />
      <p>{users.length} users online</p>
    </div>
  );
}
```

## Interfaces

- [EtherPlyProviderProps](interfaces/EtherPlyProviderProps.md)
- [PresenceUser](interfaces/PresenceUser.md)
- [UseDocumentOptions](interfaces/UseDocumentOptions.md)
- [UseDocumentReturn](interfaces/UseDocumentReturn.md)
- [UseEtherPlyOptions](interfaces/UseEtherPlyOptions.md)
- [UseEtherPlyReturn](interfaces/UseEtherPlyReturn.md)
- [UsePresenceOptions](interfaces/UsePresenceOptions.md)

## Functions

- [EtherPlyProvider](functions/EtherPlyProvider.md)
- [useDocument](functions/useDocument.md)
- [useEtherPly](functions/useEtherPly.md)
- [useEtherPlyContext](functions/useEtherPlyContext.md)
- [usePresence](functions/usePresence.md)
