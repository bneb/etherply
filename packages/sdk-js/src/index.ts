/**
 * @etherply/sdk - Official JavaScript/TypeScript SDK for EtherPly
 * 
 * Real-time sync infrastructure in 5 minutes.
 * 
 * @example
 * ```typescript
 * import { EtherPlyClient } from '@etherply/sdk';
 * 
 * const client = new EtherPlyClient({
 *   workspaceId: 'my-workspace',
 *   token: 'your-jwt-token'
 * });
 * 
 * await client.connect();
 * client.set('greeting', 'Hello, world!');
 * ```
 * 
 * @packageDocumentation
 */

export { EtherPlyClient } from './client';

// Types
export type {
    EtherPlyConfig,
    ConnectionStatus,
    MessageHandler,
    StatusHandler,
    Operation,
    InitMessage,
    OperationMessage,
    EtherPlyMessage,
} from './types';

// Errors
export {
    EtherPlyError,
    ConfigurationError,
    ConnectionError,
    AuthenticationError,
    MessageError,
    QueueOverflowError,
} from './errors';

// Utilities (for advanced users)
export { parseMessage, truncate } from './validation';
