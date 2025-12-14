# @etherply/sdk-python

The official Python SDK for [EtherPly](https://etherply.com). Designed for backend services, bots, and IoT devices.

![PyPI](https://img.shields.io/pypi/v/etherply)
![Python Version](https://img.shields.io/pypi/pyversions/etherply)
![License](https://img.shields.io/pypi/l/etherply)

## Features

- 🐍 **Asyncio-native** — Built for high-concurrency environments.
- 🤖 **Bot Ready** — Perfect for server-side logic and AI agents.
- 📡 **IoT Compatible** — Lightweight and resilient for edge devices.
- 🔒 **Type Strict** — Fully typed for robust development.

## Installation

```bash
pip install etherply
```

## Quick Start

### Basic Usage (Async)

```python
import asyncio
from etherply import EtherPlyClient

async def main():
    # 1. Initialize
    client = EtherPlyClient(
        workspace_id="my-workspace",
        token="dev-token",
        host="localhost:8080",
        secure=False  # True for production (WSS)
    )

    # 2. Connect
    await client.connect()
    print("✅ Connected!")

    # 3. Operations
    # Set a value
    await client.set("sensor:temp", 23.5)
    
    # Get a value (local state)
    temp = client.get("sensor:temp")
    print(f"Current Temp: {temp}")

    # 4. Subscribe
    async def on_update(update):
        print(f"Update received: {update.key} = {update.value}")

    client.subscribe(on_update)

    # Keep alive
    try:
        await asyncio.Future()
    except asyncio.CancelledError:
        await client.disconnect()

if __name__ == "__main__":
    asyncio.run(main())
```

## Architecture

The Python SDK uses `aiohttp` for WebSocket management. It maintains an internal generic `dict` representation of the document state, automatically merging CRDT operations received from the server.

### Key Concepts

1.  **State**: Accessible via `client.state` (read-only dict) or `client.get(key)`.
2.  **Operations**: `await client.set(key, value)` pushes an operation to the server.
3.  **Persistence**: The SDK automatically handles reconnection and state re-synchronization.

## Troubleshooting

### "Connection Refused"
- Ensure `host` does **not** include the protocol (e.g. `localhost:8080`, not `ws://localhost:8080`).
- Check `secure=True` for `wss://` (production) or `secure=False` for `ws://` (dev).

### "Task was destroyed but it is pending"
- Ensure you properly `await client.disconnect()` on shutdown.
- Use `try/finally` blocks in your main loop.

## License
MIT
