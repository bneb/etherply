# nMeshed Sync Server - Deployment Guide

This guide covers deploying the nMeshed Sync Server in various environments.

## Quick Start (Docker)

```bash
# Build the image
docker build -t nmeshed-sync-server .

# Run with required environment variables
docker run -d \
  -p 8080:8080 \
  -e NMESHED_JWT_SECRET="your-secure-secret-here" \
  -v nmeshed-data:/data \
  nmeshed-sync-server
```

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `NMESHED_JWT_SECRET` | **Yes** | - | JWT signing secret for authentication |
| `PORT` | No | `8080` | HTTP server port |
| `BADGER_PATH` | No | `./badger.db` | Path to BadgerDB data directory |
| `SHUTDOWN_TIMEOUT_SECONDS` | No | `30` | Graceful shutdown timeout |
| `WEBHOOK_URL` | No | - | URL for webhook event delivery |

## Docker Compose (Development)

```bash
# Start the server
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

## Kubernetes Deployment

### Deployment Manifest

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nmeshed-sync-server
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nmeshed-sync-server
  template:
    metadata:
      labels:
        app: nmeshed-sync-server
    spec:
      containers:
      - name: sync-server
        image: nmeshed-sync-server:latest
        ports:
        - containerPort: 8080
        env:
        - name: NMESHED_JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: nmeshed-secrets
              key: jwt-secret
        - name: BADGER_PATH
          value: "/data/badger.db"
        volumeMounts:
        - name: data
          mountPath: /data
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /readyz
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: nmeshed-data
```

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: nmeshed-sync-server
spec:
  selector:
    app: nmeshed-sync-server
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

### Secret

```bash
kubectl create secret generic nmeshed-secrets \
  --from-literal=jwt-secret='your-secure-secret-here'
```

## Health Endpoints

| Endpoint | Purpose | Use |
|----------|---------|-----|
| `/healthz` | Liveness probe | Restart if failing |
| `/readyz` | Readiness probe | Remove from LB if failing |

## Security Considerations

1. **JWT Secret**: Use a cryptographically random string (min 32 chars)
2. **Network**: Place behind a reverse proxy/ingress with TLS
3. **Volume Permissions**: Ensure data volume is only accessible by the container
4. **Image Updates**: Regularly rebuild to get security patches

## Backup & Restore

### Backup

```bash
# Stop the server gracefully
kubectl scale deployment nmeshed-sync-server --replicas=0

# Copy data (BadgerDB directory)
kubectl cp $(kubectl get pod -l app=nmeshed-sync-server -o name):/data ./backup/

# Restart
kubectl scale deployment nmeshed-sync-server --replicas=3
```

### Restore

```bash
# Scale down
kubectl scale deployment nmeshed-sync-server --replicas=0

# Restore data
kubectl cp ./backup/ $(kubectl get pod -l app=nmeshed-sync-server -o name):/data

# Scale up
kubectl scale deployment nmeshed-sync-server --replicas=3
```

## Troubleshooting

**Server won't start:**
- Check `NMESHED_JWT_SECRET` is set
- Verify `BADGER_PATH` is writable

**Health checks failing:**
- Check container logs: `docker logs <container-id>`
- Verify port is exposed correctly

**Data not persisting:**
- Ensure volume is mounted correctly
- Check volume permissions
