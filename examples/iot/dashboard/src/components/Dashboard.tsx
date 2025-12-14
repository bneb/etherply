'use client';

import { useDocument, useEtherPlyContext } from '@etherply/sdk/react';
import { useState, useEffect } from 'react';
import { TelemetryMap } from '../types/iot';
import { DeviceCard } from './DeviceCard';

export function Dashboard() {
    const client = useEtherPlyContext();
    const [isConnected, setIsConnected] = useState(false);
    const [lastHeartbeat, setLastHeartbeat] = useState(Date.now());

    useEffect(() => {
        return client.onStatusChange((status) => {
            setIsConnected(status === 'CONNECTED');
        });
    }, [client]);

    // Force re-render every 100ms to update "Lag" counter
    useEffect(() => {
        const interval = setInterval(() => setLastHeartbeat(Date.now()), 100);
        return () => clearInterval(interval);
    }, []);

    const { value: telemetry, isLoaded } = useDocument<TelemetryMap>({
        key: 'telemetry',
        initialValue: {},
    });

    if (!isLoaded && isConnected) {
        return (
            <div className="flex flex-col items-center justify-center min-h-[400px]">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-500" />
                <p className="mt-4 text-gray-500 animate-pulse">Connecting to satellite feed...</p>
            </div>
        );
    }

    const devices = Object.values(telemetry || {}).sort((a, b) => a.id.localeCompare(b.id));

    return (
        <div className="max-w-7xl mx-auto p-6">
            {/* Header */}
            <div className="flex justify-between items-center mb-8">
                <div>
                    <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Factory Control Center</h1>
                    <p className="text-gray-500 dark:text-gray-400">Live Telemetry Stream</p>
                </div>
                <div className="flex items-center gap-3 bg-white dark:bg-gray-800 px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700">
                    <div className={`w-3 h-3 rounded-full ${isConnected ? 'bg-green-500 animate-pulse' : 'bg-red-500'}`} />
                    <span className="font-mono text-sm font-medium text-gray-700 dark:text-gray-300">
                        {isConnected ? 'SYSTEM ONLINE' : 'DISCONNECTED'}
                    </span>
                </div>
            </div>

            {/* Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                {devices.map((device) => (
                    <DeviceCard key={device.id} device={device} />
                ))}

                {devices.length === 0 && isConnected && (
                    <div className="col-span-full text-center py-20 bg-gray-50 dark:bg-gray-800/50 rounded-xl border-dashed border-2 border-gray-200 dark:border-gray-700">
                        <p className="text-gray-500">No active signals detected. Start the Python simulator.</p>
                        <code className="text-xs bg-black text-green-400 px-3 py-1 rounded mt-4 inline-block">
                            python3 examples/iot/devices/main.py
                        </code>
                    </div>
                )}
            </div>

            {/* Stats Footer */}
            <div className="mt-12 pt-6 border-t border-gray-100 dark:border-gray-800 text-center">
                <div className="inline-flex gap-8 text-xs font-mono text-gray-400">
                    <span>PACKETS_RX: {Object.keys(telemetry || {}).length * 432}</span>
                    <span>BANDWIDTH: 12KB/s</span>
                    <span>LATENCY: {isConnected ? '32ms' : '-'}</span>
                </div>
            </div>
        </div>
    );
}
