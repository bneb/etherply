export type User = {
    userId: string;
    status: 'online' | 'idle';
}

export type EtherPlyConfig = {
    workspaceId: string;
    userId: string;
    token?: string;
}

export type ConnectionStatus = 'IDLE' | 'CONNECTING' | 'CONNECTED' | 'DISCONNECTED' | 'ERROR';

export class EtherPlyClient {
    private ws: WebSocket | null = null;
    private listeners: ((data: any) => void)[] = [];
    private statusListeners: ((status: ConnectionStatus) => void)[] = [];
    private config: EtherPlyConfig;
    private status: ConnectionStatus = 'IDLE';

    constructor(config: EtherPlyConfig) {
        this.config = config;
    }

    private setStatus(newStatus: ConnectionStatus) {
        this.status = newStatus;
        this.statusListeners.forEach(l => l(newStatus));
    }

    connect() {
        if (this.status === 'CONNECTED' || this.status === 'CONNECTING') return;

        this.setStatus('CONNECTING');
        // Connect to local server for demo
        const token = this.config.token || "demo-token-123";
        this.ws = new WebSocket(`ws://localhost:8080/v1/sync/${this.config.workspaceId}?userId=${this.config.userId}&token=${token}`);

        this.ws.onopen = () => {
            console.log("Connected to EtherPly");
            this.setStatus('CONNECTED');
            // Simulate Metric: demo_aha_moment check would be here or in UI
        };

        this.ws.onclose = () => {
            this.setStatus('DISCONNECTED');
            // Simple reconnect logic for MVP
            setTimeout(() => this.connect(), 3000);
        };

        this.ws.onerror = () => {
            this.setStatus('ERROR');
        };

        this.ws.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            this.listeners.forEach(l => l(msg));
        };
    }

    onMessage(callback: (data: any) => void) {
        this.listeners.push(callback);
        // Return cleanup function to unsubscribe
        return () => {
            this.listeners = this.listeners.filter(l => l !== callback);
        };
    }

    onStatusChange(callback: (status: ConnectionStatus) => void) {
        this.statusListeners.push(callback);
        // Immediate callback with current status
        callback(this.status);
        return () => {
            this.statusListeners = this.statusListeners.filter(l => l !== callback);
        };
    }

    /**
     * Sends a key-value operation to the sync server.
     * Includes a timestamp (Unix microseconds) for LWW conflict resolution.
     * @param key - The key to update (e.g., "document_text")
     * @param value - The new value for the key
     */
    sendOperation(key: string, value: any) {
        if (this.ws?.readyState === WebSocket.OPEN) {
            const start = performance.now();
            // CRITICAL: timestamp is required for LWW (Last-Write-Wins) conflict resolution.
            // Without it, operations may be dropped or lose to older operations.
            const timestamp = Date.now() * 1000; // Convert to Unix microseconds
            this.ws.send(JSON.stringify({
                type: 'op',
                payload: { key, value, timestamp }
            }));
            const latency = performance.now() - start;
            // Metric: client_sync_latency (fire to PostHog stub)
            console.log(`[METRIC] client_sync_latency | duration=${latency.toFixed(2)}ms`);
        }
    }
}
