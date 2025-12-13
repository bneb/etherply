# EtherPly Go Client SDK

## Overview
A lightweight, idiomatic Go library for connecting to the EtherPly Sync Engine.
Designed for backend-to-backend synchronization or building custom Go-based clients (e.g., CLI tools, bots).

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
    client := etherply.NewClient("ws://localhost:8080", "my-token")

    // 2. Connect
    if err := client.Connect("my-workspace"); err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }

    // 3. Listen
    client.Listen(func(msg map[string]interface{}) {
        log.Printf("Received: %v", msg)
    })

    // 4. Send
    client.SendOperation("greeting", "Hello EtherPly")
}
```

## Internals
- Use `NewClient` to configure.
- `Connect` performs the WebSocket handshake.
- `Listen` spawns a goroutine to read messages.
