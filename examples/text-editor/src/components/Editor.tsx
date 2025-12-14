'use client';

import { useDocument, useEtherPlyContext } from '@etherply/sdk/react';
import { useEffect, useState } from 'react';

export function Editor() {
    const client = useEtherPlyContext();
    const [isConnected, setIsConnected] = useState(false);

    useEffect(() => {
        return client.onStatusChange((status) => {
            setIsConnected(status === 'CONNECTED');
        });
    }, [client]);

    const { value, setValue, isLoaded } = useDocument<string>({
        key: 'document-text',
        initialValue: '',
    });

    return (
        <div className="w-full max-w-4xl mx-auto p-4">
            <div className="mb-4 flex items-center justify-between">
                <h1 className="text-2xl font-bold">Collaborative Editor</h1>
                <div className="flex items-center gap-2">
                    <div
                        className={`w-3 h-3 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'
                            }`}
                    />
                    <span className="text-sm text-gray-500">
                        {isConnected ? 'Connected' : 'Disconnected'}
                    </span>
                </div>
            </div>

            <div className="relative">
                {!isLoaded && isConnected && (
                    <div className="absolute inset-0 bg-white/50 flex items-center justify-center">
                        <span className="text-sm font-medium">Loading...</span>
                    </div>
                )}
                <textarea
                    value={value}
                    onChange={(e) => setValue(e.target.value)}
                    className="w-full h-[600px] p-4 text-lg font-mono border rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none resize-none"
                    placeholder="Start typing..."
                    disabled={!isConnected}
                />
            </div>

            <div className="mt-4 text-sm text-gray-400 text-right">
                {(value || '').length} characters
            </div>
        </div>
    );
}
