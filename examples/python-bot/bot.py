import asyncio
import logging
from etherply import EtherPlyClient, EtherPlyMessage

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)

async def main():
    # Initialize client
    client = EtherPlyClient({
        "workspace_id": "python-demo",
        "token": "dev-token-python",
        "host": "localhost:8080",
        "secure": False,
        "max_reconnect_attempts": 2
    })

    # Message handler
    def on_message(msg: EtherPlyMessage):
        if msg["type"] == "op":
            print(f"Update received: {msg['payload']['key']} = {msg['payload']['value']}")

    client.on_message(on_message)

    try:
        print("Connecting to EtherPly...")
        # Start connection in background task
        connect_task = asyncio.create_task(client.connect())
        
        # specific delay to allow connection
        await asyncio.sleep(1)
        
        print("Sending update...")
        await client.set("python-message", "Hello from Python SDK 🐍")
        
        # Keep alive for a bit to receive updates
        await asyncio.sleep(5)
        
        await client.disconnect()
        await connect_task
        
    except Exception as e:
        print(f"Demo error: {e}")

if __name__ == "__main__":
    asyncio.run(main())
