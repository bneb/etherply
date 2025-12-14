'use client';

import { EtherPlyProvider } from '@etherply/sdk/react';
import { ReactNode, useState, useEffect } from 'react';

export function Providers({ children }: { children: ReactNode }) {
    const [mounted, setMounted] = useState(false);
    const [userId] = useState(() => 'user-' + Math.random().toString(36).slice(2, 7));

    useEffect(() => {
        setMounted(true);
    }, []);

    if (!mounted) return null;

    return (
        <EtherPlyProvider
            config={{
                workspaceId: 'kanban-demo',
                token: 'dev-token',
                userId: userId,
            }}
        >
            {children as any}
        </EtherPlyProvider>
    );
}
