[**@etherply/sdk**](../../README.md)

***

# Class: EtherPlyClient

Defined in: [src/client.ts:88](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/client.ts#L88)

EtherPly client for real-time synchronization.

The client manages a WebSocket connection to an EtherPly server,
handles automatic reconnection, and provides methods for sending
and receiving state updates.

## Features
- Automatic reconnection with exponential backoff
- Connection timeout to prevent hangs
- Heartbeat pings to detect dead connections
- Operation queueing when offline
- Bounded queue to prevent memory issues
- Defensive message parsing

## Example

```typescript
const client = new EtherPlyClient({
  workspaceId: 'my-workspace',
  token: 'jwt-token'
});

await client.connect();

client.onMessage((msg) => {
  if (msg.type === 'op') {
    console.log('Update:', msg.payload.key, '=', msg.payload.value);
  }
});

client.set('greeting', 'Hello!');
```

## Constructors

### Constructor

> **new EtherPlyClient**(`config`): `EtherPlyClient`

Defined in: [src/client.ts:108](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/client.ts#L108)

Creates a new EtherPly client instance.

#### Parameters

| Parameter | Type | Description |
| ------ | ------ | ------ |
| `config` | [`EtherPlyConfig`](../interfaces/EtherPlyConfig.md) | Configuration options |

#### Returns

`EtherPlyClient`

#### Throws

If workspaceId or token is missing

## Methods

### close()

> **close**(): `void`

Defined in: [src/client.ts:637](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/client.ts#L637)

Alias for disconnect() for API consistency.

#### Returns

`void`

***

### connect()

> **connect**(): `Promise`\<`void`\>

Defined in: [src/client.ts:220](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/client.ts#L220)

Connects to the EtherPly server.

#### Returns

`Promise`\<`void`\>

A promise that resolves when connected, or rejects on error.

#### Throws

If connection fails or times out

***

### destroy()

> **destroy**(): `void`

Defined in: [src/client.ts:647](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/client.ts#L647)

Permanently destroys the client, releasing all resources.

After calling this, the client cannot be reconnected.
Use this for cleanup in React useEffect or similar.

#### Returns

`void`

***

### disconnect()

> **disconnect**(): `void`

Defined in: [src/client.ts:621](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/client.ts#L621)

Disconnects from the server.

After calling this, you can call `connect()` again to reconnect.

#### Returns

`void`

***

### get()

> **get**\<`T`\>(`key`): `T` \| `undefined`

Defined in: [src/client.ts:487](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/client.ts#L487)

Gets the current value of a key from local state.

Note: This returns the locally cached state, which may be
momentarily out of sync with the server.

#### Type Parameters

| Type Parameter | Default type |
| ------ | ------ |
| `T` | `unknown` |

#### Parameters

| Parameter | Type | Description |
| ------ | ------ | ------ |
| `key` | `string` | The key to get |

#### Returns

`T` \| `undefined`

The value, or undefined if not found

***

### getQueueSize()

> **getQueueSize**(): `number`

Defined in: [src/client.ts:612](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/client.ts#L612)

Gets the number of operations in the queue.

#### Returns

`number`

***

### getState()

> **getState**(): `Record`\<`string`, `unknown`\>

Defined in: [src/client.ts:496](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/client.ts#L496)

Gets the entire current state of the workspace.

#### Returns

`Record`\<`string`, `unknown`\>

A shallow copy of the current state

***

### getStatus()

> **getStatus**(): [`ConnectionStatus`](../type-aliases/ConnectionStatus.md)

Defined in: [src/client.ts:605](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/client.ts#L605)

Gets the current connection status.

#### Returns

[`ConnectionStatus`](../type-aliases/ConnectionStatus.md)

***

### onMessage()

> **onMessage**(`handler`): () => `void`

Defined in: [src/client.ts:566](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/client.ts#L566)

Subscribes to incoming messages.

#### Parameters

| Parameter | Type | Description |
| ------ | ------ | ------ |
| `handler` | [`MessageHandler`](../type-aliases/MessageHandler.md) | Function to call when a message is received |

#### Returns

A cleanup function to unsubscribe

> (): `void`

##### Returns

`void`

***

### onStatusChange()

> **onStatusChange**(`handler`): () => `void`

Defined in: [src/client.ts:585](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/client.ts#L585)

Subscribes to connection status changes.

The handler is called immediately with the current status.

#### Parameters

| Parameter | Type | Description |
| ------ | ------ | ------ |
| `handler` | [`StatusHandler`](../type-aliases/StatusHandler.md) | Function to call when status changes |

#### Returns

A cleanup function to unsubscribe

> (): `void`

##### Returns

`void`

***

### sendOperation()

> **sendOperation**(`key`, `value`): `void`

Defined in: [src/client.ts:508](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/client.ts#L508)

Sends an operation to update a key-value pair.

If not connected, the operation is queued and sent when reconnected.

#### Parameters

| Parameter | Type | Description |
| ------ | ------ | ------ |
| `key` | `string` | The key to update |
| `value` | `unknown` | The new value |

#### Returns

`void`

***

### set()

> **set**(`key`, `value`): `void`

Defined in: [src/client.ts:471](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/client.ts#L471)

Sets a key-value pair in the workspace.

#### Parameters

| Parameter | Type | Description |
| ------ | ------ | ------ |
| `key` | `string` | The key to set (must be non-empty string) |
| `value` | `unknown` | The value to set |

#### Returns

`void`

#### Throws

If key is invalid
