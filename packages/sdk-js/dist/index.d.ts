import { E as EtherPlyMessage } from './client-Dg1Ob3EQ.js';
export { C as ConnectionStatus, a as EtherPlyClient, b as EtherPlyConfig, I as InitMessage, M as MessageHandler, O as Operation, c as OperationMessage, S as StatusHandler } from './client-Dg1Ob3EQ.js';

/**
 * Error types for EtherPly SDK.
 *
 * Using typed errors allows consumers to handle specific failure modes.
 */
/**
 * Base class for all EtherPly errors.
 */
declare class EtherPlyError extends Error {
    readonly code: string;
    constructor(message: string, code: string);
}
/**
 * Thrown when configuration is invalid.
 */
declare class ConfigurationError extends EtherPlyError {
    constructor(message: string);
}
/**
 * Thrown when connection fails or times out.
 */
declare class ConnectionError extends EtherPlyError {
    readonly cause?: Error | undefined;
    readonly isRetryable: boolean;
    constructor(message: string, cause?: Error | undefined, isRetryable?: boolean);
}
/**
 * Thrown when authentication fails.
 */
declare class AuthenticationError extends EtherPlyError {
    constructor(message?: string);
}
/**
 * Thrown when a message fails to parse or validate.
 */
declare class MessageError extends EtherPlyError {
    readonly rawMessage?: string | undefined;
    constructor(message: string, rawMessage?: string | undefined);
}
/**
 * Thrown when the operation queue exceeds capacity.
 */
declare class QueueOverflowError extends EtherPlyError {
    constructor(maxSize: number);
}

/**
 * Message validation utilities.
 *
 * Defensive parsing to ensure messages from the server are well-formed
 * before processing. This prevents runtime errors from malformed data.
 */

/**
 * Parses and validates a raw message string from the server.
 *
 * @param raw - Raw JSON string from WebSocket
 * @returns Validated EtherPlyMessage
 * @throws {MessageError} If message is invalid
 *
 * @example
 * ```typescript
 * try {
 *   const message = parseMessage(event.data);
 *   // message is now type-safe
 * } catch (error) {
 *   if (error instanceof MessageError) {
 *     console.error('Invalid message:', error.rawMessage);
 *   }
 * }
 * ```
 */
declare function parseMessage(raw: string): EtherPlyMessage;
/**
 * Safely truncates a string for logging/error messages.
 * Prevents log explosion from large payloads.
 */
declare function truncate(str: string, maxLength?: number): string;

export { AuthenticationError, ConfigurationError, ConnectionError, EtherPlyError, EtherPlyMessage, MessageError, QueueOverflowError, parseMessage, truncate };
