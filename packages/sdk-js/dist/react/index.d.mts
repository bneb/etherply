import { b as EtherPlyConfig, C as ConnectionStatus, a as EtherPlyClient } from '../client-CtBTxa2k.mjs';
import * as react_jsx_runtime from 'react/jsx-runtime';
import { ReactNode } from 'react';

/**
 * Options for the useEtherPly hook.
 */
interface UseEtherPlyOptions extends EtherPlyConfig {
    /**
     * Callback when connected.
     */
    onConnect?: () => void;
    /**
     * Callback when disconnected.
     */
    onDisconnect?: () => void;
    /**
     * Callback when an error occurs.
     */
    onError?: (error: Error) => void;
}
/**
 * Return value of the useEtherPly hook.
 */
interface UseEtherPlyReturn {
    /**
     * Current state of the workspace as a reactive object.
     */
    state: Record<string, unknown>;
    /**
     * Set a value in the workspace.
     */
    set: (key: string, value: unknown) => void;
    /**
     * Get a value from the workspace.
     */
    get: <T = unknown>(key: string) => T | undefined;
    /**
     * Current connection status.
     */
    status: ConnectionStatus;
    /**
     * Whether the client is connected.
     */
    isConnected: boolean;
    /**
     * The underlying EtherPly client instance.
     */
    client: EtherPlyClient;
    /**
     * Manually connect to the server.
     */
    connect: () => Promise<void>;
    /**
     * Manually disconnect from the server.
     */
    disconnect: () => void;
}
/**
 * React hook for real-time synchronization with EtherPly.
 *
 * This hook creates and manages an EtherPly client, provides reactive
 * state, and handles connection lifecycle automatically.
 *
 * @param options - Configuration and callbacks
 * @returns Object with state, setters, and connection status
 *
 * @example Basic Usage
 * ```tsx
 * function CollaborativeEditor() {
 *   const { state, set, status } = useEtherPly({
 *     workspaceId: 'my-doc',
 *     token: 'jwt-token'
 *   });
 *
 *   return (
 *     <div>
 *       <p>Status: {status}</p>
 *       <textarea
 *         value={state.content as string || ''}
 *         onChange={(e) => set('content', e.target.value)}
 *       />
 *     </div>
 *   );
 * }
 * ```
 *
 * @example With Callbacks
 * ```tsx
 * const { state, set } = useEtherPly({
 *   workspaceId: 'my-doc',
 *   token: 'jwt-token',
 *   onConnect: () => console.log('Connected!'),
 *   onDisconnect: () => console.log('Disconnected'),
 *   onError: (err) => console.error('Error:', err)
 * });
 * ```
 */
declare function useEtherPly(options: UseEtherPlyOptions): UseEtherPlyReturn;

/**
 * Options for the useDocument hook.
 */
interface UseDocumentOptions<T> {
    /**
     * The key to sync.
     */
    key: string;
    /**
     * Initial value before server state is received.
     */
    initialValue?: T;
}
/**
 * Return value of the useDocument hook.
 */
interface UseDocumentReturn<T> {
    /**
     * Current value of the document.
     */
    value: T | undefined;
    /**
     * Update the value.
     */
    setValue: (newValue: T) => void;
    /**
     * Whether the value has been loaded from the server.
     */
    isLoaded: boolean;
}
/**
 * Hook to sync a single key with EtherPly.
 *
 * Provides a simple useState-like interface for a single synchronized value.
 * Must be used within an EtherPlyProvider.
 *
 * @param options - Configuration options
 * @returns Object with value, setter, and loading state
 *
 * @example
 * ```tsx
 * function Counter() {
 *   const { value, setValue, isLoaded } = useDocument<number>({
 *     key: 'counter',
 *     initialValue: 0
 *   });
 *
 *   if (!isLoaded) return <div>Loading...</div>;
 *
 *   return (
 *     <div>
 *       <p>Count: {value}</p>
 *       <button onClick={() => setValue((value || 0) + 1)}>
 *         Increment
 *       </button>
 *     </div>
 *   );
 * }
 * ```
 *
 * @example With Complex Objects
 * ```tsx
 * interface Todo {
 *   id: string;
 *   text: string;
 *   done: boolean;
 * }
 *
 * function TodoItem({ id }: { id: string }) {
 *   const { value, setValue } = useDocument<Todo>({
 *     key: `todo:${id}`
 *   });
 *
 *   if (!value) return null;
 *
 *   return (
 *     <div>
 *       <input
 *         type="checkbox"
 *         checked={value.done}
 *         onChange={() => setValue({ ...value, done: !value.done })}
 *       />
 *       <span>{value.text}</span>
 *     </div>
 *   );
 * }
 * ```
 */
declare function useDocument<T = unknown>(options: UseDocumentOptions<T>): UseDocumentReturn<T>;

interface PresenceUser {
    userId: string;
    status: string;
    last_seen: string;
}
interface UsePresenceOptions {
    /**
     * Polling interval in milliseconds.
     * Default: 10000 (10 seconds)
     */
    interval?: number;
}
/**
 * Hook to get the current presence list for the workspace.
 *
 * Note: This currently uses polling every 10 seconds.
 *
 * @param options - Configuration options
 * @returns Array of active users
 */
declare function usePresence(options?: UsePresenceOptions): PresenceUser[];

/**
 * Props for EtherPlyProvider.
 */
interface EtherPlyProviderProps {
    /**
     * Configuration for the EtherPly client.
     */
    config: EtherPlyConfig;
    /**
     * Child components that will have access to the client.
     */
    children: ReactNode;
    /**
     * Whether to automatically connect on mount.
     * @default true
     */
    autoConnect?: boolean;
}
/**
 * Provider component that creates and manages an EtherPly client.
 *
 * Wrap your app (or a portion of it) with this provider to share
 * a single client instance across multiple components.
 *
 * @example
 * ```tsx
 * import { EtherPlyProvider } from '@etherply/sdk/react';
 *
 * function App() {
 *   return (
 *     <EtherPlyProvider
 *       config={{
 *         workspaceId: 'my-workspace',
 *         token: 'jwt-token'
 *       }}
 *     >
 *       <MyCollaborativeApp />
 *     </EtherPlyProvider>
 *   );
 * }
 * ```
 */
declare function EtherPlyProvider({ config, children, autoConnect, }: EtherPlyProviderProps): react_jsx_runtime.JSX.Element;
/**
 * Hook to access the EtherPly client from context.
 *
 * Must be used within an EtherPlyProvider.
 *
 * @returns The EtherPly client instance
 * @throws {Error} If used outside of EtherPlyProvider
 *
 * @example
 * ```tsx
 * function MyComponent() {
 *   const client = useEtherPlyContext();
 *
 *   const handleClick = () => {
 *     client.set('clicked', true);
 *   };
 *
 *   return <button onClick={handleClick}>Click me</button>;
 * }
 * ```
 */
declare function useEtherPlyContext(): EtherPlyClient;

export { EtherPlyProvider, type UseEtherPlyOptions, type UseEtherPlyReturn, useDocument, useEtherPly, useEtherPlyContext, usePresence };
