[**@etherply/sdk**](../../README.md)

***

# Interface: UseEtherPlyReturn

Defined in: [src/react/useEtherPly.tsx:28](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/react/useEtherPly.tsx#L28)

Return value of the useEtherPly hook.

## Properties

### client

> **client**: [`EtherPlyClient`](../../index/classes/EtherPlyClient.md)

Defined in: [src/react/useEtherPly.tsx:57](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/react/useEtherPly.tsx#L57)

The underlying EtherPly client instance.

***

### connect()

> **connect**: () => `Promise`\<`void`\>

Defined in: [src/react/useEtherPly.tsx:62](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/react/useEtherPly.tsx#L62)

Manually connect to the server.

#### Returns

`Promise`\<`void`\>

***

### disconnect()

> **disconnect**: () => `void`

Defined in: [src/react/useEtherPly.tsx:67](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/react/useEtherPly.tsx#L67)

Manually disconnect from the server.

#### Returns

`void`

***

### get()

> **get**: \<`T`\>(`key`) => `T` \| `undefined`

Defined in: [src/react/useEtherPly.tsx:42](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/react/useEtherPly.tsx#L42)

Get a value from the workspace.

#### Type Parameters

| Type Parameter | Default type |
| ------ | ------ |
| `T` | `unknown` |

#### Parameters

| Parameter | Type |
| ------ | ------ |
| `key` | `string` |

#### Returns

`T` \| `undefined`

***

### isConnected

> **isConnected**: `boolean`

Defined in: [src/react/useEtherPly.tsx:52](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/react/useEtherPly.tsx#L52)

Whether the client is connected.

***

### set()

> **set**: (`key`, `value`) => `void`

Defined in: [src/react/useEtherPly.tsx:37](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/react/useEtherPly.tsx#L37)

Set a value in the workspace.

#### Parameters

| Parameter | Type |
| ------ | ------ |
| `key` | `string` |
| `value` | `unknown` |

#### Returns

`void`

***

### state

> **state**: `Record`\<`string`, `unknown`\>

Defined in: [src/react/useEtherPly.tsx:32](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/react/useEtherPly.tsx#L32)

Current state of the workspace as a reactive object.

***

### status

> **status**: [`ConnectionStatus`](../../index/type-aliases/ConnectionStatus.md)

Defined in: [src/react/useEtherPly.tsx:47](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/react/useEtherPly.tsx#L47)

Current connection status.
