[**@etherply/sdk**](../../README.md)

***

# Function: EtherPlyProvider()

> **EtherPlyProvider**(`__namedParameters`): `Element`

Defined in: [src/react/context.tsx:55](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/react/context.tsx#L55)

Provider component that creates and manages an EtherPly client.

Wrap your app (or a portion of it) with this provider to share
a single client instance across multiple components.

## Parameters

| Parameter | Type |
| ------ | ------ |
| `__namedParameters` | [`EtherPlyProviderProps`](../interfaces/EtherPlyProviderProps.md) |

## Returns

`Element`

## Example

```tsx
import { EtherPlyProvider } from '@etherply/sdk/react';

function App() {
  return (
    <EtherPlyProvider
      config={{
        workspaceId: 'my-workspace',
        token: 'jwt-token'
      }}
    >
      <MyCollaborativeApp />
    </EtherPlyProvider>
  );
}
```
