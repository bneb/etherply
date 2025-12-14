# Real-time Voting

A cross-language example bridging TypeScript and Python.

## Features
- **Frontend (Next.js)**: Optimistic voting UI with real-time charts.
- **Backend (Python)**: A bot (`bot.py`) that observes votes and calculates the winner.

## Architecture

1. **Voter** (JS): `client.set('votes:python', count + 1)`
2. **Bot** (Python): Listens for changes, sums totals, and updates `state.winner`.
3. **Observer** (JS): Reacts to `state.winner` updates instantly.

## Code Highlight (Python)

```python
async for message in ws:
    if message['type'] == 'op':
        recalculate_winner()
        await client.set('winner', new_winner)
```

[View Source Code](https://github.com/etherply/etherply/tree/main/examples/voting)
