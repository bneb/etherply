[**@etherply/sdk**](../../README.md)

***

# Function: useEtherPlyContext()

> **useEtherPlyContext**(): [`EtherPlyClient`](../../index/classes/EtherPlyClient.md)

Defined in: [src/react/context.tsx:110](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/react/context.tsx#L110)

Hook to access the EtherPly client from context.

Must be used within an EtherPlyProvider.

## Returns

[`EtherPlyClient`](../../index/classes/EtherPlyClient.md)

The EtherPly client instance

## Throws

If used outside of EtherPlyProvider

## Example

```tsx
function MyComponent() {
  const client = useEtherPlyContext();
  
  const handleClick = () => {
    client.set('clicked', true);
  };
  
  return <button onClick={handleClick}>Click me</button>;
}
```
