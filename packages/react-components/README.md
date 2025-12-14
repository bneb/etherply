# @etherply/react-components

Magical UI primitives for [EtherPly](https://etherply.com). Drop-in multiplayer features with zero config.

![npm version](https://img.shields.io/npm/v/@etherply/react-components)
![License](https://img.shields.io/npm/l/@etherply/react-components)

## Features

- ✨ **Zero Config** — Just import and render.
- 🚀 **Performance Optimized** — Built-in throttling (30fps) and interpolation.
- 🎨 **Beautiful Defaults** — Auto-generated colors and smooth animations.
- 🧩 **Customizable** — Override rendering while keeping the logic.

## Installation

```bash
npm install @etherply/react-components @etherply/sdk framer-motion
```

> **Note**: `framer-motion` and `@etherply/sdk` are peer dependencies.

## Components

### `<LiveCursors />`

Adds real-time multiplayer cursors to your application. Handles mouse/touch tracking, broadcasting, and rendering.

#### Usage

```tsx
import { EtherPlyProvider } from '@etherply/sdk/react';
import { LiveCursors } from '@etherply/react-components';

export default function App() {
  return (
    <EtherPlyProvider config={{ workspaceId: 'demo', token: 'dev' }}>
      
      {/* 1. Drop it in */}
      <LiveCursors />
      
      <YourContent />
    </EtherPlyProvider>
  );
}
```

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `throttleMs` | `number` | `33` | Throttle interval (ms) for network updates. |
| `timeoutMs` | `number` | `30000` | Time (ms) before removing inactive cursors. |
| `renderCursor` | `Function` | `undefined` | Custom render function. |

#### Custom Rendering

```tsx
<LiveCursors
  renderCursor={({ x, y, color, label }) => (
    <div style={{ transform: `translate(${x}px, ${y}px)` }}>
      Custom: {label}
    </div>
  )}
/>
```

## Architecture

This package subscribes to the `presence:*` key space in EtherPly. It effectively "side-loads" ephemeral data alongside your durable document state.

## License
MIT
