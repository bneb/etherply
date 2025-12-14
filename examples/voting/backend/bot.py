import asyncio
import logging
from typing import Dict, TypedDict
from etherply import EtherPlyClient, EtherPlyMessage

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger("VotingBot")

class VoteOption(TypedDict):
    id: str
    label: str
    votes: int

class PollState(TypedDict):
    question: str
    options: Dict[str, VoteOption]
    totalVotes: int
    winnerId: str

class VotingBot:
    def __init__(self):
        self.client = EtherPlyClient({
            "workspace_id": "voting-demo",
            "token": "bot-token",
            "user_id": "python-bot"
        })
        self.client.on_message(self.handle_message)
        self.state: PollState = {} # type: ignore

    def handle_message(self, msg: EtherPlyMessage):
        if msg["type"] == "init" and msg.get("data"):
            data = msg["data"]
            if "poll-data" in data:
                self.state = data["poll-data"]
                self.recalculate()
        
        elif msg["type"] == "op":
            payload = msg["payload"]
            if payload and payload["key"] == "poll-data":
                self.state = payload["value"]
                # If the update didn't come from us, recalculate
                if self.state.get("lastUpdatedBy") != "python-bot":
                    self.recalculate()

    def recalculate(self):
        """Calculate totals and winner."""
        if not self.state or "options" not in self.state:
            return

        options = self.state["options"]
        total_votes = sum(opt["votes"] for opt in options.values())
        
        # Determine winner
        winner_id = ""
        max_votes = -1
        
        for opt_id, opt in options.items():
            if opt["votes"] > max_votes:
                max_votes = opt["votes"]
                winner_id = opt_id
            elif opt["votes"] == max_votes:
                winner_id = "" # Tie (simplified)

        # Check if anything changed to avoid loops
        if (total_votes != self.state.get("totalVotes") or 
            winner_id != self.state.get("winnerId")):
            
            logger.info(f"Recalculating: Total={total_votes}, Winner={winner_id}")
            
            # Update state
            self.state["totalVotes"] = total_votes
            self.state["winnerId"] = winner_id
            self.state["lastUpdatedBy"] = "python-bot" # Prevent loops
            
            # Push update
            asyncio.create_task(self.client.set("poll-data", self.state))

    async def start(self):
        logger.info("Starting Voting Bot...")
        await self.client.connect()

if __name__ == "__main__":
    bot = VotingBot()
    try:
        asyncio.run(bot.start())
    except KeyboardInterrupt:
        pass
