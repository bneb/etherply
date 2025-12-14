import { useState, useEffect, useRef } from 'react';
import { useEtherPlyContext } from './context';

export interface PresenceUser {
    userId: string;
    status: string;
    last_seen: string;
}

export interface UsePresenceOptions {
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
export function usePresence(options: UsePresenceOptions = {}): PresenceUser[] {
    const { interval = 10000 } = options;
    const client = useEtherPlyContext();
    const [users, setUsers] = useState<PresenceUser[]>([]);

    // Polling ref to avoid effect depending on interval changing frequently
    const savedCallback = useRef<() => void>();

    useEffect(() => {
        const fetchPresence = async () => {
            if (client.getStatus() !== 'CONNECTED') return;
            try {
                const data = await client.getPresence();
                setUsers(data);
            } catch (err) {
                console.warn('Failed to fetch presence:', err);
            }
        };
        savedCallback.current = fetchPresence;

        // Initial fetch
        fetchPresence();
    }, [client]);

    useEffect(() => {
        function tick() {
            if (savedCallback.current) {
                savedCallback.current();
            }
        }
        if (interval !== null && interval !== undefined) {
            const id = setInterval(tick, interval);
            return () => clearInterval(id);
        }
        return undefined;
    }, [interval]);

    return users;
}
