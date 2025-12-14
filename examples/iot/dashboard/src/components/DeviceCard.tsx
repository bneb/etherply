'use client';

import { Activity, Gauge, Thermometer, Fan, AlertTriangle } from 'lucide-react';
import { clsx } from 'clsx';
import { DeviceState } from '../types/iot';

interface DeviceCardProps {
    device: DeviceState;
}

export function DeviceCard({ device }: DeviceCardProps) {
    const isStale = Date.now() - device.last_update > 2000;

    const statusColor = {
        'IDLE': 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400',
        'RUNNING': 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
        'WARNING': 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
        'CRITICAL': 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
    }[device.status];

    return (
        <div className={clsx(
            "relative overflow-hidden rounded-xl border p-6 transition-all",
            isStale ? "opacity-60 grayscale border-gray-200" : "bg-white dark:bg-gray-800 border-gray-100 dark:border-gray-700 shadow-lg",
            device.status === 'CRITICAL' && "animate-pulse ring-2 ring-red-500"
        )}>
            {isStale && (
                <div className="absolute inset-0 flex items-center justify-center bg-gray-50/80 dark:bg-gray-900/80 z-10 backdrop-blur-[1px]">
                    <span className="bg-gray-800 text-white px-3 py-1 rounded-full text-xs font-mono">OFFLINE</span>
                </div>
            )}

            <div className="flex justify-between items-start mb-6">
                <div>
                    <h3 className="text-lg font-bold text-gray-900 dark:text-white">{device.name}</h3>
                    <p className="text-xs text-gray-500 font-mono mt-1">{device.id}</p>
                </div>
                <div className={clsx("px-2.5 py-1 rounded-full text-xs font-bold", statusColor)}>
                    {device.status}
                </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
                {/* Temperature */}
                <div className="space-y-1">
                    <div className="flex items-center gap-1.5 text-gray-500 dark:text-gray-400 text-xs uppercase tracking-wider font-semibold">
                        <Thermometer size={14} /> Temp
                    </div>
                    <div className={clsx(
                        "text-2xl font-mono font-medium",
                        device.temperature > 85 ? "text-red-500" : "text-gray-900 dark:text-gray-100"
                    )}>
                        {device.temperature.toFixed(1)}°C
                    </div>
                </div>

                {/* RPM */}
                <div className="space-y-1">
                    <div className="flex items-center gap-1.5 text-gray-500 dark:text-gray-400 text-xs uppercase tracking-wider font-semibold">
                        <Fan size={14} /> RPM
                    </div>
                    <div className="text-2xl font-mono font-medium text-gray-900 dark:text-gray-100">
                        {device.rpm}
                    </div>
                </div>

                {/* Pressure */}
                <div className="space-y-1">
                    <div className="flex items-center gap-1.5 text-gray-500 dark:text-gray-400 text-xs uppercase tracking-wider font-semibold">
                        <Gauge size={14} /> Press
                    </div>
                    <div className="text-2xl font-mono font-medium text-gray-900 dark:text-gray-100">
                        {device.pressure.toFixed(0)} hPa
                    </div>
                </div>

                {/* Last Update */}
                <div className="space-y-1">
                    <div className="flex items-center gap-1.5 text-gray-500 dark:text-gray-400 text-xs uppercase tracking-wider font-semibold">
                        <Activity size={14} /> Lag
                    </div>
                    <div className="text-2xl font-mono font-medium text-gray-900 dark:text-gray-100">
                        {Math.max(0, Date.now() - device.last_update)}ms
                    </div>
                </div>
            </div>
        </div>
    );
}
