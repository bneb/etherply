"use client";

import { useEffect, useState } from 'react';
import { EtherPlyClient, User } from '@/lib/etherply-client';

export default function PresenceWidget({ client }: { client: EtherPlyClient }) {
    const [users, setUsers] = useState<User[]>([]);

    useEffect(() => {
        // In a real app, we would listen for specific presence events.
        // For MVP with basic Go server, we might poll or get pushed updates.
        // To keep it simple, let's assume the WebSocket broadcasts presence or we fetch it.
        // We'll define a simple fetch poll for this MVP demo since the WB interface for presence was simple REST in the plan.

        const fetchPresence = async () => {
            // Assuming client.config has workspaceId (we need to expose it or pass it)
            // We will just hardcode fetching from the known endpoint for the demo sake
            // In reality, this should be fully handled by the SDK.
            try {
                // We need to know where the server is. 
                // client library has hardcoded localhost:8080 or config.
                const res = await fetch('http://localhost:8080/v1/presence/demo-workspace');
                if (res.ok) {
                    const data = await res.json();
                    setUsers(data || []);
                }
            } catch (e) {
                console.error("Presence fetch error", e);
            }
        };

        const interval = setInterval(fetchPresence, 2000); // Poll every 2s for presence updates
        fetchPresence();

        return () => clearInterval(interval);
    }, []);

    return (
        <div className="flex -space-x-2 overflow-hidden items-center">
            {users.map((user, i) => (
                <div key={user.userId + i} title={user.userId} className="relative inline-block h-8 w-8 rounded-full ring-2 ring-white bg-gradient-to-br from-brand-500 to-purple-500 flex items-center justify-center text-white text-xs font-bold uppercase">
                    {user.userId.slice(0, 2)}
                    <span className="absolute bottom-0 right-0 block h-2.5 w-2.5 rounded-full bg-green-400 ring-2 ring-white" />
                </div>
            ))}
            <div className="ml-4 text-sm text-gray-500">
                {users.length} Active
            </div>
        </div>
    );
}
