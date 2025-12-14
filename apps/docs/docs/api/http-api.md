# EtherPly API Reference

## WebSocket API
**Endpoint:** `wss://api.etherply.com/v1/sync/{workspace_id}`

> [!TIP]
> **Local Development:** Use `ws://localhost:8080/v1/sync/{workspace_id}` when running the server locally.

### Connection Handshake
To establish a connection, you must perform a standard WebSocket handshake with authentication.

**Query Parameters:**
- `token` (Required): A valid JWT Bearer token.
- `userId` (Optional): A string identifier for the user (useful for debugging/dev).

**Example (JS Client):**
```javascript
const socket = new WebSocket('ws://localhost:8080/v1/sync/my-workspace?token=eyJh...');
```

### Protocol Messages

#### 1. Client -> Server: Operation (Op)
Send this message to mutate state.

```json
{
  "type": "op",
  "payload": {
    "key": "documentTitle",
    "value": "My New Thesis",
    "timestamp": 1702483200000000
  }
}
```

> [!IMPORTANT]
> **Timestamp Format**: The `timestamp` field MUST be in **Unix Microseconds**.
> - JavaScript: `Date.now() * 1000`
> - Go: `time.Now().UnixMicro()`
> Failure to provide this will result in the server rejecting the operation or unpredictable LWW behavior.

#### 2. Server -> Client: Initial State (Init)
Sent immediately after connection is accepted. Contains the full current state of the workspace.

```json
{
  "type": "init",
  "data": {
    "documentTitle": "Old Thesis Title",
    "cursorPosition": 12
  }
}
```

#### 3. Server -> Client: Broadcast (Op)
Sent to *other* clients when a change is made. The sender does NOT receive their own echo.

```json
{
  "type": "op",
  "payload": {
    "key": "documentTitle",
    "value": "My New Thesis"
  }
}
```

---

## REST API

### Get Presence
Retrieve the list of currently active users in a workspace.

**Endpoint:** `GET /v1/presence/{workspace_id}`

#### Request Example (cURL)
```bash
curl -v -H "Authorization: Bearer <YOUR_TOKEN>" \
http://localhost:8080/v1/presence/workspace-123
```

#### Success Response (200 OK)
```json
[
  {
    "user_id": "user_abc123",
    "status": "online"
  },
  {
    "user_id": "user_xyz789",
    "status": "online"
  }
]
```

#### Error Response (401 Unauthorized)
```json
{
  "error": "Unauthorized",
  "message": "Invalid authentication credentials"
}
```

### Health Check (Implicit)
EtherPly's sync server doesn't currently expose a dedicated `/health` endpoint, but you can verify connectivity by hitting the presence endpoint with an invalid ID.

**Command:**
```bash
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/v1/presence/health-check-probe
```
**Expected Output:** `200` (returns empty list `[]`)
**What If It Fails?**
- If you get `000` or `Connection refused`: The server process is not running. Check logs with `tail -f std.log` or verify `PORT`.
- If you get `404`: The routing configuration is incorrect.
