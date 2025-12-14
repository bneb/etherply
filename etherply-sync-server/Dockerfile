# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy dependencies first to leverage cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
# CGO_ENABLED=0 for static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o etherply-server cmd/server/main.go

# Run Stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies (ca-certificates for HTTPS)
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/etherply-server .

# Expose port
EXPOSE 8080

# Environment variables with defaults
ENV PORT=8080
ENV LOG_FORMAT=json
ENV LOG_LEVEL=info
ENV SYNC_STRATEGY=automerge
ENV BADGER_PATH=/data/badger.db

# Create data directory for persistence
RUN mkdir -p /data

# Volume for persistence
VOLUME ["/data"]

# Run
CMD ["./etherply-server"]
