export interface DeviceState {
    id: string;
    name: string;
    type: string;
    temperature: number;
    pressure: number;
    rpm: number;
    status: 'IDLE' | 'RUNNING' | 'WARNING' | 'CRITICAL';
    last_update: number;
}

export type TelemetryMap = Record<string, DeviceState>;
