'use client';

import { useDocument, useEtherPlyContext, usePresence } from '@etherply/sdk/react';
import { useEffect, useState } from 'react';

export function Editor() {
    const client = useEtherPlyContext();
    const [isConnected, setIsConnected] = useState(false);

    // Fetch presence (polls every 5s)
    const users = usePresence({ interval: 5000 });

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
        <div className="w-full max-w-6xl mx-auto p-4 flex gap-4">
            <div className="flex-1">
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

            {/* Presence Sidebar */}
            <div className="w-64 bg-gray-50 p-4 rounded-lg border h-[600px] mt-[52px]">
                <h3 className="font-semibold text-gray-700 mb-4 flex items-center justify-between">
                    Active Users
                    <span className="bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded-full">
                        {users.length}
                    </span>
                </h3>
                <div className="space-y-3">
                    {users.map((user) => (
                        <div key={user.userId} className="flex items-center gap-3">
                            <div className="w-8 h-8 rounded-full bg-gradient-to-br from-blue-400 to-indigo-500 flex items-center justify-center text-white text-xs font-bold ring-2 ring-white shadow-sm">
                                {user.userId.slice(5, 7).toUpperCase()}
                            </div>
                            <div className="flex flex-col">
                                <span className="text-sm font-medium text-gray-700">
                                    {user.userId}
                                </span>
                                <span className="text-xs text-green-600 flex items-center gap-1">
                                    <div className="w-1.5 h-1.5 rounded-full bg-green-500" />
                                    Online
                                </span>
                            </div>
                        </div>
                    ))}
                    {users.length === 0 && (
                        <p className="text-gray-400 text-sm text-center italic mt-10">
                            No active users
                        </p>
                    )}
                </div>
            </div>
        </div>
    );
}
