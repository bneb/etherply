# IoT Dashboard

Simulating high-frequency telemetry from embedded devices.

## Features
- **Python Simulator**: Generates physics-based sensor data (Temp, RPM).
- **Batching**: Streams 10Hz updates efficiently.
- **Alerting**: Dashboard highlights critical states (e.g. Overheating).

## Code Highlight

```python
# Streaming telemetry map
await client.set("telemetry", {
    "device-1": { temp: 88.5, status: "CRITICAL" },
    "device-2": { temp: 42.1, status: "OK" }
})
```

[View Source Code](https://github.com/etherply/etherply/tree/main/examples/iot)
