[**@etherply/sdk**](../../README.md)

***

# Interface: UseEtherPlyOptions

Defined in: [src/react/useEtherPly.tsx:8](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/react/useEtherPly.tsx#L8)

Options for the useEtherPly hook.

## Extends

- [`EtherPlyConfig`](../../index/interfaces/EtherPlyConfig.md)

## Properties

### autoReconnect?

> `optional` **autoReconnect**: `boolean`

Defined in: [src/types.ts:36](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/types.ts#L36)

Enable automatic reconnection on disconnect.

#### Default

```ts
true
```

#### Inherited from

[`EtherPlyConfig`](../../index/interfaces/EtherPlyConfig.md).[`autoReconnect`](../../index/interfaces/EtherPlyConfig.md#autoreconnect)

***

### connectionTimeout?

> `optional` **connectionTimeout**: `number`

Defined in: [src/types.ts:63](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/types.ts#L63)

Timeout (in ms) for initial connection.
If connection isn't established within this time, it fails.

#### Default

```ts
10000
```

#### Inherited from

[`EtherPlyConfig`](../../index/interfaces/EtherPlyConfig.md).[`connectionTimeout`](../../index/interfaces/EtherPlyConfig.md#connectiontimeout)

***

### debug?

> `optional` **debug**: `boolean`

Defined in: [src/types.ts:84](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/types.ts#L84)

Enable debug logging to console.

#### Default

```ts
false
```

#### Inherited from

[`EtherPlyConfig`](../../index/interfaces/EtherPlyConfig.md).[`debug`](../../index/interfaces/EtherPlyConfig.md#debug)

***

### heartbeatInterval?

> `optional` **heartbeatInterval**: `number`

Defined in: [src/types.ts:70](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/types.ts#L70)

Interval (in ms) to send heartbeat pings.
Set to 0 to disable heartbeats.

#### Default

```ts
30000
```

#### Inherited from

[`EtherPlyConfig`](../../index/interfaces/EtherPlyConfig.md).[`heartbeatInterval`](../../index/interfaces/EtherPlyConfig.md#heartbeatinterval)

***

### maxQueueSize?

> `optional` **maxQueueSize**: `number`

Defined in: [src/types.ts:78](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/types.ts#L78)

Maximum number of operations to queue while disconnected.
When exceeded, oldest operations are dropped (FIFO).
Set to 0 for unlimited (not recommended).

#### Default

```ts
1000
```

#### Inherited from

[`EtherPlyConfig`](../../index/interfaces/EtherPlyConfig.md).[`maxQueueSize`](../../index/interfaces/EtherPlyConfig.md#maxqueuesize)

***

### maxReconnectAttempts?

> `optional` **maxReconnectAttempts**: `number`

Defined in: [src/types.ts:42](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/types.ts#L42)

Maximum number of reconnection attempts.

#### Default

```ts
10
```

#### Inherited from

[`EtherPlyConfig`](../../index/interfaces/EtherPlyConfig.md).[`maxReconnectAttempts`](../../index/interfaces/EtherPlyConfig.md#maxreconnectattempts)

***

### maxReconnectDelay?

> `optional` **maxReconnectDelay**: `number`

Defined in: [src/types.ts:56](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/types.ts#L56)

Maximum delay (in ms) between reconnection attempts.
Caps the exponential backoff to prevent excessive waits.

#### Default

```ts
30000
```

#### Inherited from

[`EtherPlyConfig`](../../index/interfaces/EtherPlyConfig.md).[`maxReconnectDelay`](../../index/interfaces/EtherPlyConfig.md#maxreconnectdelay)

***

### onConnect()?

> `optional` **onConnect**: () => `void`

Defined in: [src/react/useEtherPly.tsx:12](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/react/useEtherPly.tsx#L12)

Callback when connected.

#### Returns

`void`

***

### onDisconnect()?

> `optional` **onDisconnect**: () => `void`

Defined in: [src/react/useEtherPly.tsx:17](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/react/useEtherPly.tsx#L17)

Callback when disconnected.

#### Returns

`void`

***

### onError()?

> `optional` **onError**: (`error`) => `void`

Defined in: [src/react/useEtherPly.tsx:22](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/react/useEtherPly.tsx#L22)

Callback when an error occurs.

#### Parameters

| Parameter | Type |
| ------ | ------ |
| `error` | `Error` |

#### Returns

`void`

***

### reconnectBaseDelay?

> `optional` **reconnectBaseDelay**: `number`

Defined in: [src/types.ts:49](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/types.ts#L49)

Base delay (in ms) between reconnection attempts.
Uses exponential backoff: delay * 2^attempt

#### Default

```ts
1000
```

#### Inherited from

[`EtherPlyConfig`](../../index/interfaces/EtherPlyConfig.md).[`reconnectBaseDelay`](../../index/interfaces/EtherPlyConfig.md#reconnectbasedelay)

***

### serverUrl?

> `optional` **serverUrl**: `string`

Defined in: [src/types.ts:30](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/types.ts#L30)

WebSocket server URL.
Defaults to 'wss://api.etherply.com' in production.

#### Default

```ts
'wss://api.etherply.com'
```

#### Inherited from

[`EtherPlyConfig`](../../index/interfaces/EtherPlyConfig.md).[`serverUrl`](../../index/interfaces/EtherPlyConfig.md#serverurl)

***

### token

> **token**: `string`

Defined in: [src/types.ts:17](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/types.ts#L17)

JWT authentication token.
Required for production use.
Get one from dashboard.etherply.com

#### Inherited from

[`EtherPlyConfig`](../../index/interfaces/EtherPlyConfig.md).[`token`](../../index/interfaces/EtherPlyConfig.md#token)

***

### userId?

> `optional` **userId**: `string`

Defined in: [src/types.ts:23](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/types.ts#L23)

Optional user identifier for presence tracking.
If not provided, a random ID will be generated.

#### Inherited from

[`EtherPlyConfig`](../../index/interfaces/EtherPlyConfig.md).[`userId`](../../index/interfaces/EtherPlyConfig.md#userid)

***

### workspaceId

> **workspaceId**: `string`

Defined in: [src/types.ts:10](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/types.ts#L10)

The workspace ID to connect to.
A workspace is a collaborative room or document.

#### Example

```ts
'my-project-123'
```

#### Inherited from

[`EtherPlyConfig`](../../index/interfaces/EtherPlyConfig.md).[`workspaceId`](../../index/interfaces/EtherPlyConfig.md#workspaceid)
