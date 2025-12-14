'use client';

import { EtherPlyProvider } from '@etherply/sdk/react';
import { ReactNode, useState, useEffect } from 'react';

export function Providers({ children }: { children: ReactNode }) {
    const [mounted, setMounted] = useState(false);
    const [userId] = useState(() => 'monitor-' + Math.random().toString(36).slice(2, 7));

    useEffect(() => {
        setMounted(true);
    }, []);

    if (!mounted) return null;

    return (
        <EtherPlyProvider
            config={{
                workspaceId: 'iot-demo',
                token: 'monitor-token',
                userId: userId,
            }}
        >
            {/* Cast to any to avoid React types mismatch between linked packages */}
            {children as any}
        </EtherPlyProvider>
    );
}
