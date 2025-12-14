[**@etherply/sdk**](../../README.md)

***

# Function: usePresence()

> **usePresence**(`options`): [`PresenceUser`](../interfaces/PresenceUser.md)[]

Defined in: src/react/usePresence.tsx:26

Hook to get the current presence list for the workspace.

Note: This currently uses polling every 10 seconds.

## Parameters

| Parameter | Type | Description |
| ------ | ------ | ------ |
| `options` | [`UsePresenceOptions`](../interfaces/UsePresenceOptions.md) | Configuration options |

## Returns

[`PresenceUser`](../interfaces/PresenceUser.md)[]

Array of active users
