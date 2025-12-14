# Build stage
FROM golang:1.23-alpine AS builder

# Install git for fetching dependencies (some may use git)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
# CGO_ENABLED=0 for static linking (no glibc dependency)
# -ldflags="-s -w" strips debug info for smaller binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /etherply-sync-server .

# Runtime stage - using distroless for minimal attack surface
FROM gcr.io/distroless/static-debian12:nonroot

# Copy binary from builder
COPY --from=builder /etherply-sync-server /etherply-sync-server

# Copy CA certificates for HTTPS connections (webhooks, external services)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Create data directory for BadgerDB
# Note: distroless doesn't have mkdir, so we rely on volume mounts
VOLUME ["/data"]

# Environment variables with sensible defaults
ENV PORT=8080
ENV BADGER_PATH=/data/badger.db
ENV SHUTDOWN_TIMEOUT_SECONDS=30
# ETHERPLY_JWT_SECRET must be provided at runtime
# WEBHOOK_URL is optional

# Expose the server port
EXPOSE 8080

# Health check for container orchestrators
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD ["/etherply-sync-server", "-health-check"] || exit 1

# Run as non-root user (nonroot user from distroless)
USER nonroot:nonroot

# Start the server
ENTRYPOINT ["/etherply-sync-server"]
