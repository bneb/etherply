/**
 * Message validation utilities.
 * 
 * Defensive parsing to ensure messages from the server are well-formed
 * before processing. This prevents runtime errors from malformed data.
 */

import type {
    EtherPlyMessage,
    InitMessage,
    OperationMessage,
    PresenceMessage,
    PresenceUser,
} from './types';
import { MessageError } from './errors';

/**
 * Type guard for InitMessage.
 */
function isInitMessage(data: unknown): data is InitMessage {
    if (typeof data !== 'object' || data === null) return false;
    const obj = data as Record<string, unknown>;
    return (
        obj.type === 'init' &&
        typeof obj.data === 'object' &&
        obj.data !== null
    );
}

/**
 * Type guard for OperationMessage.
 */
function isOperationMessage(data: unknown): data is OperationMessage {
    if (typeof data !== 'object' || data === null) return false;
    const obj = data as Record<string, unknown>;
    if (obj.type !== 'op') return false;

    const payload = obj.payload;
    if (typeof payload !== 'object' || payload === null) return false;

    const p = payload as Record<string, unknown>;
    return (
        typeof p.key === 'string' &&
        p.key.length > 0 &&
        'value' in p &&
        typeof p.timestamp === 'number'
    );
}

/**
 * Type guard for PresenceUser.
 */
function isPresenceUser(data: unknown): data is PresenceUser {
    if (typeof data !== 'object' || data === null) return false;
    const obj = data as Record<string, unknown>;
    return (
        typeof obj.userId === 'string' &&
        (obj.status === 'online' || obj.status === 'idle')
    );
}

/**
 * Type guard for PresenceMessage.
 */
function isPresenceMessage(data: unknown): data is PresenceMessage {
    if (typeof data !== 'object' || data === null) return false;
    const obj = data as Record<string, unknown>;
    if (obj.type !== 'presence') return false;

    const users = obj.users;
    if (!Array.isArray(users)) return false;

    return users.every(isPresenceUser);
}

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
export function parseMessage(raw: string): EtherPlyMessage {
    // Step 1: Parse JSON
    let data: unknown;
    try {
        data = JSON.parse(raw);
    } catch (error) {
        throw new MessageError(
            `Failed to parse message as JSON: ${error instanceof Error ? error.message : 'Unknown error'}`,
            raw
        );
    }

    // Step 2: Validate structure
    if (typeof data !== 'object' || data === null) {
        throw new MessageError('Message must be an object', raw);
    }

    const obj = data as Record<string, unknown>;
    if (typeof obj.type !== 'string') {
        throw new MessageError('Message must have a "type" field', raw);
    }

    // Step 3: Validate specific message types
    if (isInitMessage(data)) {
        return data;
    }

    if (isOperationMessage(data)) {
        return data;
    }

    if (isPresenceMessage(data)) {
        return data;
    }

    throw new MessageError(
        `Unknown or malformed message type: "${obj.type}"`,
        raw
    );
}

/**
 * Safely truncates a string for logging/error messages.
 * Prevents log explosion from large payloads.
 */
export function truncate(str: string, maxLength: number = 200): string {
    if (str.length <= maxLength) return str;
    return str.substring(0, maxLength) + `... (${str.length - maxLength} more chars)`;
}
