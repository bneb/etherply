[**@etherply/sdk**](../../README.md)

***

# Function: truncate()

> **truncate**(`str`, `maxLength`): `string`

Defined in: [src/validation.ts:140](https://github.com/bneb/etherply/blob/303e876d3c20fd29bbe2577e2c1219ab266cb9d9/packages/sdk-js/src/validation.ts#L140)

Safely truncates a string for logging/error messages.
Prevents log explosion from large payloads.

## Parameters

| Parameter | Type | Default value |
| ------ | ------ | ------ |
| `str` | `string` | `undefined` |
| `maxLength` | `number` | `200` |

## Returns

`string`
