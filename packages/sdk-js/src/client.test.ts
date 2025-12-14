import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { EtherPlyClient } from './client';
import { ConfigurationError, ConnectionError } from './errors';

// Mock WebSocket with all static constants
class MockWebSocket {
    static instances: MockWebSocket[] = [];

    // WebSocket ready state constants (must match browser WebSocket)
    static readonly CONNECTING = 0;
    static readonly OPEN = 1;
    static readonly CLOSING = 2;
    static readonly CLOSED = 3;

    readyState = MockWebSocket.CONNECTING;
    onopen: (() => void) | null = null;
    onclose: ((event: { code: number; reason: string }) => void) | null = null;
    onerror: ((event: unknown) => void) | null = null;
    onmessage: ((event: { data: string }) => void) | null = null;

    constructor(public url: string) {
        MockWebSocket.instances.push(this);
    }

    send = vi.fn();
    close = vi.fn();

    simulateOpen() {
        this.readyState = MockWebSocket.OPEN;
        this.onopen?.();
    }

    simulateMessage(data: unknown) {
        this.onmessage?.({ data: JSON.stringify(data) });
    }

    simulateClose(code = 1000, reason = '') {
        this.readyState = MockWebSocket.CLOSED;
        this.onclose?.({ code, reason });
    }

    simulateError() {
        this.onerror?.({});
    }
}

const originalWebSocket = global.WebSocket;

beforeEach(() => {
    MockWebSocket.instances = [];
    global.WebSocket = MockWebSocket as unknown as typeof WebSocket;
    vi.useFakeTimers();
});

afterEach(() => {
    global.WebSocket = originalWebSocket;
    vi.useRealTimers();
});

describe('EtherPlyClient', () => {
    const defaultConfig = {
        workspaceId: 'test-workspace',
        token: 'test-token',
    };

    describe('constructor', () => {
        it('throws ConfigurationError if workspaceId is missing', () => {
            expect(() => new EtherPlyClient({ workspaceId: '', token: 'token' }))
                .toThrow(ConfigurationError);
        });

        it('throws ConfigurationError if token is missing', () => {
            expect(() => new EtherPlyClient({ workspaceId: 'workspace', token: '' }))
                .toThrow(ConfigurationError);
        });

        it('throws ConfigurationError for invalid config', () => {
            expect(() => new EtherPlyClient({
                workspaceId: 'workspace',
                token: 'token',
                maxReconnectAttempts: -1,
            })).toThrow(ConfigurationError);
        });

        it('creates client with valid config', () => {
            const client = new EtherPlyClient(defaultConfig);
            expect(client.getStatus()).toBe('IDLE');
        });
    });

    describe('connect', () => {
        it('establishes WebSocket connection', async () => {
            const client = new EtherPlyClient(defaultConfig);
            const connectPromise = client.connect();
            MockWebSocket.instances[0].simulateOpen();
            await connectPromise;
            expect(client.getStatus()).toBe('CONNECTED');
        });

        it('resolves immediately if already connected', async () => {
            const client = new EtherPlyClient(defaultConfig);
            const promise1 = client.connect();
            MockWebSocket.instances[0].simulateOpen();
            await promise1;
            await client.connect();
            expect(MockWebSocket.instances.length).toBe(1);
        });

        it('rejects with ConnectionError on error', async () => {
            const client = new EtherPlyClient(defaultConfig);
            const connectPromise = client.connect();
            MockWebSocket.instances[0].simulateError();
            await expect(connectPromise).rejects.toThrow(ConnectionError);
        });

        it('times out if connection takes too long', async () => {
            const client = new EtherPlyClient({
                ...defaultConfig,
                connectionTimeout: 1000,
            });
            const connectPromise = client.connect();
            vi.advanceTimersByTime(1001);
            await expect(connectPromise).rejects.toThrow('timed out');
        });
    });

    describe('messaging', () => {
        it('handles init message and updates state', async () => {
            const client = new EtherPlyClient(defaultConfig);
            const connectPromise = client.connect();
            const ws = MockWebSocket.instances[0];
            ws.simulateOpen();
            await connectPromise;
            ws.simulateMessage({
                type: 'init',
                data: { greeting: 'Hello', count: 42 },
            });
            expect(client.get('greeting')).toBe('Hello');
            expect(client.get('count')).toBe(42);
        });

        it('handles op message and updates state', async () => {
            const client = new EtherPlyClient(defaultConfig);
            const connectPromise = client.connect();
            MockWebSocket.instances[0].simulateOpen();
            await connectPromise;
            MockWebSocket.instances[0].simulateMessage({
                type: 'op',
                payload: { key: 'title', value: 'New Title', timestamp: 123 },
            });
            expect(client.get('title')).toBe('New Title');
        });

        it('notifies message listeners', async () => {
            const client = new EtherPlyClient(defaultConfig);
            const listener = vi.fn();
            client.onMessage(listener);
            const connectPromise = client.connect();
            MockWebSocket.instances[0].simulateOpen();
            await connectPromise;
            MockWebSocket.instances[0].simulateMessage({
                type: 'op',
                payload: { key: 'test', value: 'value', timestamp: 123 },
            });
            expect(listener).toHaveBeenCalled();
        });

        it('allows unsubscribing from messages', async () => {
            const client = new EtherPlyClient(defaultConfig);
            const listener = vi.fn();
            const unsubscribe = client.onMessage(listener);
            const connectPromise = client.connect();
            MockWebSocket.instances[0].simulateOpen();
            await connectPromise;
            unsubscribe();
            MockWebSocket.instances[0].simulateMessage({
                type: 'op',
                payload: { key: 'test', value: 'value', timestamp: 123 },
            });
            expect(listener).not.toHaveBeenCalled();
        });
    });

    describe('sendOperation', () => {
        it('sends operation when connected', async () => {
            const client = new EtherPlyClient(defaultConfig);
            const connectPromise = client.connect();
            const ws = MockWebSocket.instances[0];
            ws.simulateOpen();
            await connectPromise;
            client.set('key', 'value');
            expect(ws.send).toHaveBeenCalled();
        });

        it('queues operations when disconnected', () => {
            const client = new EtherPlyClient(defaultConfig);
            client.set('queued', 'value');
            expect(client.get('queued')).toBe('value');
            expect(client.getQueueSize()).toBe(1);
        });

        it('flushes queue on connect', async () => {
            const client = new EtherPlyClient(defaultConfig);
            client.set('key1', 'value1');
            client.set('key2', 'value2');
            const connectPromise = client.connect();
            const ws = MockWebSocket.instances[0];
            ws.simulateOpen();
            await connectPromise;
            expect(ws.send).toHaveBeenCalledTimes(2);
            expect(client.getQueueSize()).toBe(0);
        });
    });

    describe('status changes', () => {
        it('notifies status listeners', async () => {
            const client = new EtherPlyClient(defaultConfig);
            const listener = vi.fn();
            client.onStatusChange(listener);
            expect(listener).toHaveBeenCalledWith('IDLE');
            const connectPromise = client.connect();
            expect(listener).toHaveBeenCalledWith('CONNECTING');
            MockWebSocket.instances[0].simulateOpen();
            await connectPromise;
            expect(listener).toHaveBeenCalledWith('CONNECTED');
        });
    });

    describe('disconnect', () => {
        it('closes WebSocket connection', async () => {
            const client = new EtherPlyClient(defaultConfig);
            const connectPromise = client.connect();
            const ws = MockWebSocket.instances[0];
            ws.simulateOpen();
            await connectPromise;
            client.disconnect();
            expect(ws.close).toHaveBeenCalled();
            expect(client.getStatus()).toBe('DISCONNECTED');
        });
    });

    describe('destroy', () => {
        it('cleans up all resources', async () => {
            const client = new EtherPlyClient(defaultConfig);
            const connectPromise = client.connect();
            MockWebSocket.instances[0].simulateOpen();
            await connectPromise;
            client.set('key', 'value');
            client.destroy();
            expect(client.getStatus()).toBe('DISCONNECTED');
            expect(client.getState()).toEqual({});
            expect(client.getQueueSize()).toBe(0);
        });
    });

    describe('getState', () => {
        it('returns a copy of the state', async () => {
            const client = new EtherPlyClient(defaultConfig);
            const connectPromise = client.connect();
            const ws = MockWebSocket.instances[0];
            ws.simulateOpen();
            await connectPromise;
            ws.simulateMessage({
                type: 'init',
                data: { a: 1, b: 2 },
            });
            const state = client.getState();
            expect(state).toEqual({ a: 1, b: 2 });
            state.c = 3;
            expect(client.get('c')).toBeUndefined();
        });
    });
});
