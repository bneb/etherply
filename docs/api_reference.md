# EtherPly API Reference

## WebSocket API
**Endpoint:** `wss://api.etherply.com/v1/sync/{workspace_id}`

> **Local Development:** Use `ws://localhost:8080/v1/sync/{workspace_id}` when running the server locally.

### Query Parameters
- `token`: (Required) Bearer token for authentication.
- `userId`: (Optional) Stub for overriding user ID in dev.

### Messages

#### Client -> Server
**Operation**
```json
{
  "type": "op",
  "payload": {
    "key": "string",
    "value": "any"
  }
}
```

#### Server -> Client
**Initial State**
```json
{
  "type": "init",
  "data": {
    "key": "value"
  }
}
```
**Operation Broadcast**
```json
{
  "type": "op",
  "payload": {
    "key": "string",
    "value": "any"
  }
}
```

## REST API

### Get Presence
**GET** `/v1/presence/{workspace_id}`

**Response**
```json
[
  {
    "user_id": "string",
    "status": "online"
  }
]
```
