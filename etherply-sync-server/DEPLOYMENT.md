# Deployment Guide

EtherPly is designed to be "Enterprise-Ready" out of the box. We provide official Docker images and Compose files for self-hosting.

## Prerequisites
- Docker Engine 20.10+
- Docker Compose v2.0+

## Quick Start (On-Premise)

1. **Clone and Navigate**:
   ```bash
   git clone https://github.com/etherply/etherply.git
   cd etherply/etherply-sync-server
   ```

2. **Start Services**:
   ```bash
   docker-compose up -d
   ```

3. **Verify**:
   The server will be available at `http://localhost:8080`.
   - Health Check: `curl http://localhost:8080/healthz`
   - Metrics: `curl http://localhost:8080/metrics`

## Configuration

1. **Create a `.env` file**:
   ```bash
   cp .env.example .env
   # Or create manually:
   # PORT=8080
   # ETHERPLY_JWT_SECRET=production-secret-123
   ```

2. **Run with Config**:
   ```bash
   docker-compose --env-file .env up -d
   ```

### Reference Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP Port | `8080` |
| `ETHERPLY_JWT_SECRET` | **Required**. Secret for verifying tokens. | - |
| `SYNC_STRATEGY` | Algorithm (`automerge`, `lww`) | `automerge` |
| `BADGER_PATH` | Path to DB file inside container | `/data/badger.db` |

## Production Notes

- **Persistence**: Data is stored in the `etherply_data` Docker volume. Ensure this volume is backed up.
- **Security**: 
  - Change `etherply_JWT_SECRET` to a strong random string.
  - Run behind a reverse proxy (Nginx/Traefik) for TLS termination.
