'use client';

import { useEtherPly } from '@etherply/sdk/react';
import { useEffect, useRef, useState } from 'react';
import { Cursor } from './Cursor';

const COLORS = ['#f87171', '#fb923c', '#fbbf24', '#a3e635', '#34d399', '#22d3ee', '#818cf8', '#e879f9'];

function getRandomColor(id: string) {
    let hash = 0;
    for (let i = 0; i < id.length; i++) {
        hash = id.charCodeAt(i) + ((hash << 5) - hash);
    }
    return COLORS[Math.abs(hash) % COLORS.length];
}

export function CursorOverlay() {
    const { client, state, isConnected } = useEtherPly({ workspaceId: 'cursors-demo', token: 'dev' });
    const [myId] = useState(() => 'user-' + Math.random().toString(36).slice(2, 7));
    const rafRef = useRef<number | null>(null);
    const lastUpdateRef = useRef<number>(0);

    useEffect(() => {
        if (!isConnected) return;

        function handleMouseMove(e: MouseEvent) {
            const now = Date.now();
            // Throttle to ~30fps to avoid flooding websocket
            if (now - lastUpdateRef.current < 33) return;

            lastUpdateRef.current = now;

            const position = {
                x: e.clientX,
                y: e.clientY,
                userId: myId,
                lastUpdate: now
            };

            // In a real app we'd use a dedicated presence channel.
            // Here we simulate it with document keys.
            client.set(`presence:${myId}`, position);
        }

        window.addEventListener('mousemove', handleMouseMove);
        return () => window.removeEventListener('mousemove', handleMouseMove);
    }, [client, isConnected, myId]);

    // Clean up stale cursors (optional, for demo polish)
    const cursors = Object.entries(state)
        .filter(([key, value]) => {
            // Filter for presence keys, exclude my own, and check for recent updates
            return key.startsWith('presence:') &&
                key !== `presence:${myId}` &&
                (value as any)?.lastUpdate > Date.now() - 30000; // 30s timeout
        })
        .map(([key, value]) => {
            const data = value as any;
            return {
                id: data.userId,
                x: data.x,
                y: data.y,
                color: getRandomColor(data.userId),
            };
        });

    return (
        <div className="fixed inset-0 pointer-events-none overflow-hidden">
            {cursors.map((cursor) => (
                <Cursor
                    key={cursor.id}
                    color={cursor.color}
                    x={cursor.x}
                    y={cursor.y}
                    label={cursor.id}
                />
            ))}

            {/* Connection Status Badge */}
            <div className="fixed bottom-4 right-4 pointer-events-auto flex items-center gap-2 bg-white dark:bg-gray-800 p-2 rounded-full shadow-lg border border-gray-200 dark:border-gray-700">
                <div
                    className={`w-3 h-3 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'
                        }`}
                />
                <span className="text-sm font-medium pr-2 text-gray-900 dark:text-gray-100">
                    {isConnected ? 'Live' : 'Connecting'}
                </span>
            </div>
        </div>
    );
}
