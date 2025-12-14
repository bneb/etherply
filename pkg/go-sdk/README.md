# EtherPly Go SDK

The official Go client for [EtherPly](https://etherply.com).
Build high-performance, real-time backend services, synchronizers, and bots.

![Go Reference](https://pkg.go.dev/badge/github.com/bneb/etherply/pkg/go-sdk.svg)
![Build Status](https://img.shields.io/github/actions/workflow/status/bneb/etherply/ci.yml)
![License](https://img.shields.io/github/license/bneb/etherply)

## Features

- **Type-Safe**: Native Go structs and interfaces.
- **Resilient**: Automatic reconnection with exponential backoff.
- **Concurrent**: Thread-safe design powered by channels.
- **Efficient**: Binary message format support (future).

## Installation

```bash
go get github.com/bneb/etherply/pkg/go-sdk
```

## Quick Start

### 1. Connect to the EtherPly Cloud

```go
package main

import (
    "context"
    "log"

    "github.com/bneb/etherply/pkg/go-sdk/client"
)

func main() {
    // Initialize the client
    c := client.New(client.Config{
        WorkspaceID: "my-workspace",
        Token:       "your-jwt-token",
        Host:        "api.etherply.com", // or localhost:8080
        Secure:      true,               // false for localhost
    })

    // Connect (blocking or async)
    if err := c.Connect(context.Background()); err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer c.Close()

    log.Println("✅ Connected to EtherPly!")
}
```

### 2. Subscribe to Updates

Use Go channels to handle real-time updates without blocking your main loop.

```go
// Create a channel for updates
updates := make(chan client.Update, 100)

// Subscribe
c.Subscribe(updates)

// Process in a goroutine
go func() {
    for update := range updates {
        log.Printf("Received update: Key=%s Value=%v", update.Key, update.Value)
    }
}()
```

### 3. Mutate State

Sending operations is thread-safe and asynchronous.

```go
// Set a simple value
if err := c.Set("status", "active"); err != nil {
    log.Println("Error sending op:", err)
}

// Set a complex object (automatically marshaled)
data := map[string]interface{}{
    "users": 42,
    "load":  0.85,
}
c.Set("metrics", data)
```

## Configuration

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `WorkspaceID` | `string` | **Required** | The room/workspace identifier. |
| `Token` | `string` | **Required** | Access token. |
| `Host` | `string` | `api.etherply.com` | Server hostname. |
| `Secure` | `bool` | `true` | Use WSS (TLS) vs WS. |
| `LogLevel` | `int` | `0` (Info) | Logging verbosity. |

## Error Handling

The SDK uses standard Go errors. You should check for `client.ErrConnectionClosed` if you need to detect disconnects manually, although the client handles reconnection automatically.

```go
if err := c.Connect(ctx); err != nil {
    // Handle initial connection failure
}
```

## License
MIT
