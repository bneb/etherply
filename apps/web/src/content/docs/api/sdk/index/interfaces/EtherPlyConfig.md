[**@etherply/sdk**](../../README.md)

***

# Interface: EtherPlyConfig

Defined in: [src/types.ts:4](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/types.ts#L4)

Configuration options for the EtherPly client.

## Extended by

- [`UseEtherPlyOptions`](../../react/interfaces/UseEtherPlyOptions.md)

## Properties

### autoReconnect?

> `optional` **autoReconnect**: `boolean`

Defined in: [src/types.ts:36](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/types.ts#L36)

Enable automatic reconnection on disconnect.

#### Default

```ts
true
```

***

### connectionTimeout?

> `optional` **connectionTimeout**: `number`

Defined in: [src/types.ts:63](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/types.ts#L63)

Timeout (in ms) for initial connection.
If connection isn't established within this time, it fails.

#### Default

```ts
10000
```

***

### debug?

> `optional` **debug**: `boolean`

Defined in: [src/types.ts:84](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/types.ts#L84)

Enable debug logging to console.

#### Default

```ts
false
```

***

### heartbeatInterval?

> `optional` **heartbeatInterval**: `number`

Defined in: [src/types.ts:70](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/types.ts#L70)

Interval (in ms) to send heartbeat pings.
Set to 0 to disable heartbeats.

#### Default

```ts
30000
```

***

### maxQueueSize?

> `optional` **maxQueueSize**: `number`

Defined in: [src/types.ts:78](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/types.ts#L78)

Maximum number of operations to queue while disconnected.
When exceeded, oldest operations are dropped (FIFO).
Set to 0 for unlimited (not recommended).

#### Default

```ts
1000
```

***

### maxReconnectAttempts?

> `optional` **maxReconnectAttempts**: `number`

Defined in: [src/types.ts:42](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/types.ts#L42)

Maximum number of reconnection attempts.

#### Default

```ts
10
```

***

### maxReconnectDelay?

> `optional` **maxReconnectDelay**: `number`

Defined in: [src/types.ts:56](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/types.ts#L56)

Maximum delay (in ms) between reconnection attempts.
Caps the exponential backoff to prevent excessive waits.

#### Default

```ts
30000
```

***

### reconnectBaseDelay?

> `optional` **reconnectBaseDelay**: `number`

Defined in: [src/types.ts:49](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/types.ts#L49)

Base delay (in ms) between reconnection attempts.
Uses exponential backoff: delay * 2^attempt

#### Default

```ts
1000
```

***

### serverUrl?

> `optional` **serverUrl**: `string`

Defined in: [src/types.ts:30](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/types.ts#L30)

WebSocket server URL.
Defaults to 'wss://api.etherply.com' in production.

#### Default

```ts
'wss://api.etherply.com'
```

***

### token

> **token**: `string`

Defined in: [src/types.ts:17](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/types.ts#L17)

JWT authentication token.
Required for production use.
Get one from dashboard.etherply.com

***

### userId?

> `optional` **userId**: `string`

Defined in: [src/types.ts:23](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/types.ts#L23)

Optional user identifier for presence tracking.
If not provided, a random ID will be generated.

***

### workspaceId

> **workspaceId**: `string`

Defined in: [src/types.ts:10](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/types.ts#L10)

The workspace ID to connect to.
A workspace is a collaborative room or document.

#### Example

```ts
'my-project-123'
```
