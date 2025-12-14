[**@etherply/sdk**](../../README.md)

***

# Function: truncate()

> **truncate**(`str`, `maxLength`): `string`

Defined in: [src/validation.ts:140](https://github.com/bneb/etherply/blob/cacd548d9b6df1697db3259d47218d0d3be5e820/packages/sdk-js/src/validation.ts#L140)

Safely truncates a string for logging/error messages.
Prevents log explosion from large payloads.

## Parameters

| Parameter | Type | Default value |
| ------ | ------ | ------ |
| `str` | `string` | `undefined` |
| `maxLength` | `number` | `200` |

## Returns

`string`
