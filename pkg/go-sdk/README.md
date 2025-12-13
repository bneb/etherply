# EtherPly Go Client SDK

## Overview
A lightweight, idiomatic Go library for connecting to the EtherPly Sync Engine via WebSockets.
Designed for backend-to-backend synchronization (e.g., seeding data) or building custom Go-based clients (e.g., CLI tools, bots).

## Installation
```bash
go get github.com/bneb/etherply/pkg/go-sdk
```

## Quick Start
```go
package main

import (
    "log"
    "github.com/bneb/etherply/pkg/go-sdk"
)

func main() {
    // 1. Initialize
    // Tip: In production, fetch tokens from your auth provider, do not hardcode.
    client := etherply.NewClient("ws://localhost:8080", "valid-jwt-token")

    // 2. Connect
    // This performs the WebSocket handshake and authenticates.
    if err := client.Connect("my-workspace"); err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer client.Close()

    // 3. Listen (Non-blocking)
    client.Listen(func(msg map[string]interface{}) {
        log.Printf("Received: %v", msg)
    })

    // 4. Send Operation
    // Operations are LWW (Last-Write-Wins). The SDK automatically attaches a microsecond timestamp.
    if err := client.SendOperation("greeting", "Hello EtherPly"); err != nil {
         log.Printf("Send failed: %v", err)
    }
    
    // Keep main thread alive
    select {}
}
```

## Architecture & Thread Safety

### Concurrency Model
- **`SendOperation`**: Thread-safe. You can call this from multiple goroutines.
- **`Listen`**: Runs on its own goroutine. The callback is executed sequentially for each incoming message.
- **`Close`**: Thread-safe.

### Error Handling
The SDK follows Go idioms. Errors are returned explicitly.
- If `Connect()` fails, the client is unusable.
- If `SendOperation()` fails, it typically means the connection dropped. The current SDK version does **not** automatically reconnect (see [Roadmap](../docs/tech_debt.md)).

## Internals
- **Protocol**: Uses `gorilla/websocket`.
- **Serialization**: JSON.
- **Timestamps**: Uses `time.Now().UnixMicro()` for CRDT conflict resolution.
