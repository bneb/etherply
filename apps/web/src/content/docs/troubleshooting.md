# Troubleshooting Guide

Common issues and solutions when working with EtherPly.

## Connection Issues

### "Connection Refused"
**Symptoms**: The client fails to connect, and the console shows `WebSocket connection to 'ws://...' failed`.
**Cause**: The EtherPly sync server is not running or is not reachable at the configured host.
**Solution**:
- Ensure the backend is running (`go run cmd/server/main.go`).
- Check the `host` configuration in `EtherPlyProvider`.
- Verify firewall settings.

### "Authentication Failed"
**Symptoms**: The client connects but immediately receives an error or disconnects with code 4001.
**Cause**: The provided `token` is invalid or expired.
**Solution**:
- Verify the token generation logic on your backend.
- Ensure the token is signed with the correct secret.

## Synchronization Issues

### Updates not appearing on other clients
**Symptoms**: You make changes, but other connected clients don't see them.
**Cause**:
- Clients might be connected to different `workspaceId`s.
- `client.set()` calls might be failing (check console for errors).
- React components might not be subscribing correctly (ensure you use `useDocument` or `useEtherPly` hooks).

### "Duplicate Key" or State Conflicts
**Symptoms**: Data flickers or reverts to old values.
**Cause**: Two clients setting the same key simultaneously without conflict resolution logic.
**Solution**:
- For scalar values (strings, numbers), Last-Write-Wins (LWW) applies.
- For complex logic, consider using unique keys or atomic operations (coming soon).

## Python SDK

### "Module not found: etherply"
**Symptoms**: `ImportError: No module named 'etherply'`
**Solution**: Ensure you have installed the package: `pip install etherply` or `pip install -e packages/sdk-python` for local dev.

### "AsyncIO Event Loop Closed"
**Symptoms**: Runtime errors about closed event loops.
**Solution**: Do not mix sync and async code improperly. Use `asyncio.run()` for entry points.
