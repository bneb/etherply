[**@etherply/sdk**](../../README.md)

***

# Function: EtherPlyProvider()

> **EtherPlyProvider**(`__namedParameters`): `Element`

Defined in: [src/react/context.tsx:55](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/react/context.tsx#L55)

Provider component that creates and manages an EtherPly client.

Wrap your app (or a portion of it) with this provider to share
a single client instance across multiple components.

## Parameters

| Parameter | Type |
| ------ | ------ |
| `__namedParameters` | `EtherPlyProviderProps` |

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
