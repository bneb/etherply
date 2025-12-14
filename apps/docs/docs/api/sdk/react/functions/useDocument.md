[**@etherply/sdk**](../../README.md)

***

# Function: useDocument()

> **useDocument**\<`T`\>(`options`): `UseDocumentReturn`\<`T`\>

Defined in: [src/react/useDocument.tsx:98](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/react/useDocument.tsx#L98)

Hook to sync a single key with EtherPly.

Provides a simple useState-like interface for a single synchronized value.
Must be used within an EtherPlyProvider.

## Type Parameters

| Type Parameter | Default type |
| ------ | ------ |
| `T` | `unknown` |

## Parameters

| Parameter | Type | Description |
| ------ | ------ | ------ |
| `options` | `UseDocumentOptions`\<`T`\> | Configuration options |

## Returns

`UseDocumentReturn`\<`T`\>

Object with value, setter, and loading state

## Examples

```tsx
function Counter() {
  const { value, setValue, isLoaded } = useDocument<number>({
    key: 'counter',
    initialValue: 0
  });
  
  if (!isLoaded) return <div>Loading...</div>;
  
  return (
    <div>
      <p>Count: {value}</p>
      <button onClick={() => setValue((value || 0) + 1)}>
        Increment
      </button>
    </div>
  );
}
```

```tsx
interface Todo {
  id: string;
  text: string;
  done: boolean;
}

function TodoItem({ id }: { id: string }) {
  const { value, setValue } = useDocument<Todo>({
    key: `todo:${id}`
  });
  
  if (!value) return null;
  
  return (
    <div>
      <input
        type="checkbox"
        checked={value.done}
        onChange={() => setValue({ ...value, done: !value.done })}
      />
      <span>{value.text}</span>
    </div>
  );
}
```
