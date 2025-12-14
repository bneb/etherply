import { createContext, useContext, useRef, useEffect, type ReactNode } from 'react';
import { EtherPlyClient } from '../client';
import type { EtherPlyConfig } from '../types';

/**
 * Context for sharing an EtherPly client across components.
 */
const EtherPlyContext = createContext<EtherPlyClient | null>(null);

/**
 * Props for EtherPlyProvider.
 */
export interface EtherPlyProviderProps {
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
export function EtherPlyProvider({
    config,
    children,
    autoConnect = true,
}: EtherPlyProviderProps) {
    const clientRef = useRef<EtherPlyClient | null>(null);

    // Create client on first render
    if (!clientRef.current) {
        clientRef.current = new EtherPlyClient(config);
    }

    useEffect(() => {
        const client = clientRef.current;
        if (!client) return;

        if (autoConnect) {
            client.connect().catch((error) => {
                console.error('[EtherPly] Auto-connect failed:', error);
            });
        }

        return () => {
            client.disconnect();
        };
    }, [autoConnect]);

    return (
        <EtherPlyContext.Provider value={clientRef.current}>
            {children}
        </EtherPlyContext.Provider>
    );
}

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
export function useEtherPlyContext(): EtherPlyClient {
    const client = useContext(EtherPlyContext);

    if (!client) {
        throw new Error(
            'useEtherPlyContext must be used within an EtherPlyProvider. ' +
            'Wrap your component tree with <EtherPlyProvider>.'
        );
    }

    return client;
}
