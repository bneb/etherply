[**@etherply/sdk**](../../README.md)

***

# Function: parseMessage()

> **parseMessage**(`raw`): [`EtherPlyMessage`](../type-aliases/EtherPlyMessage.md)

Defined in: [src/validation.ts:95](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/validation.ts#L95)

Parses and validates a raw message string from the server.

## Parameters

| Parameter | Type | Description |
| ------ | ------ | ------ |
| `raw` | `string` | Raw JSON string from WebSocket |

## Returns

[`EtherPlyMessage`](../type-aliases/EtherPlyMessage.md)

Validated EtherPlyMessage

## Throws

If message is invalid

## Example

```typescript
try {
  const message = parseMessage(event.data);
  // message is now type-safe
} catch (error) {
  if (error instanceof MessageError) {
    console.error('Invalid message:', error.rawMessage);
  }
}
```
