[**@etherply/sdk**](../README.md)

***

# index

@etherply/sdk - Official JavaScript/TypeScript SDK for EtherPly

Real-time sync infrastructure in 5 minutes.

## Example

```typescript
import { EtherPlyClient } from '@etherply/sdk';

const client = new EtherPlyClient({
  workspaceId: 'my-workspace',
  token: 'your-jwt-token'
});

await client.connect();
client.set('greeting', 'Hello, world!');
```

## Classes

- [AuthenticationError](classes/AuthenticationError.md)
- [ConfigurationError](classes/ConfigurationError.md)
- [ConnectionError](classes/ConnectionError.md)
- [EtherPlyClient](classes/EtherPlyClient.md)
- [EtherPlyError](classes/EtherPlyError.md)
- [MessageError](classes/MessageError.md)
- [QueueOverflowError](classes/QueueOverflowError.md)

## Interfaces

- [EtherPlyConfig](interfaces/EtherPlyConfig.md)
- [InitMessage](interfaces/InitMessage.md)
- [Operation](interfaces/Operation.md)
- [OperationMessage](interfaces/OperationMessage.md)

## Type Aliases

- [ConnectionStatus](type-aliases/ConnectionStatus.md)
- [EtherPlyMessage](type-aliases/EtherPlyMessage.md)
- [MessageHandler](type-aliases/MessageHandler.md)
- [StatusHandler](type-aliases/StatusHandler.md)

## Functions

- [parseMessage](functions/parseMessage.md)
- [truncate](functions/truncate.md)
