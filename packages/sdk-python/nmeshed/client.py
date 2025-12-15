import asyncio
import json
import logging
import time
import websockets
from typing import Optional, Callable, Dict, Set, Any
from .types import EtherPlyConfig, EtherPlyMessage
from .errors import ConfigurationError, ConnectionError

DEFAULT_CONFIG: EtherPlyConfig = {
    "host": "localhost:8080",
    "secure": False,
    "auto_reconnect": True,
    "max_reconnect_attempts": 5,
    "reconnect_base_delay": 1000,
}

class EtherPlyClient:
    def __init__(self, config: EtherPlyConfig):
        if not config.get("workspace_id"):
            raise ConfigurationError("workspace_id is required")
        if not config.get("token"):
            raise ConfigurationError("token is required")
            
        self.config = {**DEFAULT_CONFIG, **config}
        self.ws: Optional[websockets.WebSocketClientProtocol] = None
        self.state: Dict[str, Any] = {}
        self.message_handlers: Set[Callable[[EtherPlyMessage], None]] = set()
        self.reconnect_attempts = 0
        self._running = False
        self._logger = logging.getLogger("EtherPly")

    @property
    def url(self) -> str:
        protocol = "wss" if self.config["secure"] else "ws"
        host = self.config["host"]
        ws_id = self.config["workspace_id"]
        return f"{protocol}://{host}/ws/{ws_id}"

    async def connect(self):
        """Connect to the EtherPly server."""
        self._running = True
        while self._running:
            try:
                headers = {"Authorization": f"Bearer {self.config['token']}"}
                self._logger.info(f"Connecting to {self.url}")
                
                async with websockets.connect(self.url, additional_headers=headers) as ws:
                    self.ws = ws
                    self.reconnect_attempts = 0
                    self._logger.info("Connected")
                    
                    # Listen loop
                    async for message in ws:
                        await self._handle_message(message)
                        
            except Exception as e:
                self._logger.error(f"Connection error: {e}")
                if not self.config["auto_reconnect"]:
                    raise ConnectionError(f"Failed to connect: {e}")
                
                if self.reconnect_attempts >= self.config["max_reconnect_attempts"]:
                    self._logger.error("Max reconnect attempts reached")
                    self._running = False
                    raise ConnectionError("Max reconnect attempts reached")
                
                delay = (self.config["reconnect_base_delay"] / 1000) * (2 ** self.reconnect_attempts)
                self._logger.info(f"Reconnecting in {delay}s...")
                await asyncio.sleep(delay)
                self.reconnect_attempts += 1

    async def disconnect(self):
        """Disconnect from the server."""
        self._running = False
        if self.ws:
            await self.ws.close()

    async def set(self, key: str, value: Any):
        """Set a value in the workspace."""
        if not self.ws:
            raise ConnectionError("Not connected")
            
        # Optimistic update
        self.state[key] = value
        
        message = {
            "type": "op",
            "payload": {
                "key": key,
                "value": value,
                "timestamp": int(time.time() * 1000000) # Microseconds
            }
        }
        await self.ws.send(json.dumps(message))

    def on_message(self, handler: Callable[[EtherPlyMessage], None]):
        """Register a message handler."""
        self.message_handlers.add(handler)

    async def _handle_message(self, raw_message: str):
        try:
            data: EtherPlyMessage = json.loads(raw_message)
            
            # Update local state
            if data["type"] == "init" and data.get("data"):
                self.state.update(data["data"])
            elif data["type"] == "op" and data.get("payload"):
                payload = data["payload"]
                self.state[payload["key"]] = payload["value"]
            
            # Notify handlers
            for handler in self.message_handlers:
                try:
                    handler(data)
                except Exception as e:
                    self._logger.error(f"Error in message handler: {e}")
                    
        except json.JSONDecodeError:
            self._logger.error("Received invalid JSON")
