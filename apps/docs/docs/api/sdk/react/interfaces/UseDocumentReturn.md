[**@etherply/sdk**](../../README.md)

***

# Interface: UseDocumentReturn\<T\>

Defined in: [src/react/useDocument.tsx:23](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/react/useDocument.tsx#L23)

Return value of the useDocument hook.

## Type Parameters

| Type Parameter |
| ------ |
| `T` |

## Properties

### isLoaded

> **isLoaded**: `boolean`

Defined in: [src/react/useDocument.tsx:37](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/react/useDocument.tsx#L37)

Whether the value has been loaded from the server.

***

### setValue()

> **setValue**: (`newValue`) => `void`

Defined in: [src/react/useDocument.tsx:32](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/react/useDocument.tsx#L32)

Update the value.

#### Parameters

| Parameter | Type |
| ------ | ------ |
| `newValue` | `T` |

#### Returns

`void`

***

### value

> **value**: `T` \| `undefined`

Defined in: [src/react/useDocument.tsx:27](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/react/useDocument.tsx#L27)

Current value of the document.
