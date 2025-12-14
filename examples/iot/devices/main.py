import asyncio
import logging
import random
import time
from typing import Dict, TypedDict, List
from etherply import EtherPlyClient

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger("IoTSimulator")

class DeviceState(TypedDict):
    id: str
    name: str
    type: str
    temperature: float
    pressure: float
    rpm: int
    status: str
    last_update: int

class DeviceSimulator:
    def __init__(self, device_id: str, name: str, device_type: str):
        self.id = device_id
        self.name = name
        self.type = device_type
        self.temp = 40.0 + random.random() * 20
        self.pressure = 1000.0
        self.rpm = 0
        self.status = "IDLE"
        self.running = True

    def update(self) -> DeviceState:
        # Simulate physics
        if self.status == "RUNNING":
            target_rpm = 2000
            self.rpm += (target_rpm - self.rpm) * 0.1
            self.temp += 0.5 if self.temp < 85 else -0.1
            self.pressure = 1000 + (self.rpm / 100) + (random.random() * 5)
        else:
            self.rpm *= 0.95
            self.temp *= 0.99
            self.pressure = 1000

        # Random status changes
        if random.random() < 0.01:
            self.status = "RUNNING" if self.status == "IDLE" else "IDLE"
        
        # Alerts
        if self.temp > 80:
            self.status = "WARNING"
        if self.temp > 90:
            self.status = "CRITICAL"

        # Add noise
        self.temp += (random.random() - 0.5)
        
        return {
            "id": self.id,
            "name": self.name,
            "type": self.type,
            "temperature": round(self.temp, 1),
            "pressure": round(self.pressure, 1),
            "rpm": int(self.rpm),
            "status": self.status,
            "last_update": int(time.time() * 1000)
        }

class IoTGateway:
    def __init__(self):
        self.client = EtherPlyClient({
            "workspace_id": "iot-demo",
            "token": "gateway-token",
            "user_id": "python-gateway",
            "host": "localhost:8080", 
            "secure": False
        })
        self.devices: List[DeviceSimulator] = [
            DeviceSimulator("dev-001", "Factory Cell A", "cnc-lathe"),
            DeviceSimulator("dev-002", "Factory Cell B", "robotic-arm"),
            DeviceSimulator("dev-003", "Warehouse HVAC", "hvac-unit"),
            DeviceSimulator("dev-004", "Generator 1", "diesel-gen"),
        ]

    async def run(self):
        logger.info("Starting IoT Gateway...")
        await self.client.connect()
        
        logger.info("Connected. Streaming telemetry...")
        
        while True:
            # Collect telemetry from all devices
            telemetry: Dict[str, DeviceState] = {}
            for device in self.devices:
                state = device.update()
                telemetry[device.id] = state
            
            # Push batch update
            # We use a single 'telemetry' key containing all devices map
            # In a real app we might shard this, but map sync is efficient here
            await self.client.set("telemetry", telemetry)
            
            # High frequency: 10Hz updates
            await asyncio.sleep(0.1)

if __name__ == "__main__":
    gateway = IoTGateway()
    try:
        asyncio.run(gateway.run())
    except KeyboardInterrupt:
        pass
