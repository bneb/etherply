# The First Mile

> "Simplicity is the ultimate sophistication." — Leonardo da Vinci

You are here because you want to build something **alive**. 
The web we grew up with was a library—static pages, cached assets, solitary experiences. 
We are building the **Living Web**. Where cursors dance, text flows between screens, and state is a shared hallucination.

EtherPly is the engine for this new reality. It is not just a WebSocket wrapper; it is an opinionated, conflict-free, friction-less sync engine designed for the high-speed collision of ideas.

## The Ritual (Installation)

We believe in "Zero Friction." You should be multiplayer in the time it takes to brew a pour-over.

```bash
npm install @etherply/sdk @etherply/react-components
```

## The Invocation (Usage)

Do not manage WebSocket connections. Do not parse JSON. Do not handle reconnection logic.
Simply **declare** your intent.

### 1. The Provider
Wrap your application in the ether.

```tsx
import { EtherPlyProvider } from '@etherply/sdk/react';

<EtherPlyProvider config={{ workspaceId: 'room-1', token: 'dev' }}>
  <App />
</EtherPlyProvider>
```

### 2. The Magic
Drop in our pre-fabricated "Magic Components" to instantly enliven the space.

```tsx
import { LiveCursors } from '@etherply/react-components';

export default function Canvas() {
  return (
    <div>
        <LiveCursors /> {/* ✨ Presence, throttled and interpolated. */}
        <YourContent />
    </div>
  );
}
```

### 3. The State
Synchronize data as easily as `useState`.

```tsx
import { useDocument } from '@etherply/sdk/react';

const { value, setValue } = useDocument({ key: 'manifesto', initialValue: '' });
```

---

## The Library
We have prepared a suite of examples to guide your hand.

| Composition | Path | Use Case |
|---|---|---|
| **Text Editor** | [`examples/text-editor`](./examples/text-editor) | The "Hello World" of sync. |
| **Kanban Board** | [`examples/kanban`](./examples/kanban) | Complex, nested state. |
| **Cursors** | [`examples/cursors`](./examples/cursors) | High-frequency ephemeral data. |
| **IoT Dashboard** | [`examples/iot`](./examples/iot) | Machine-to-Machine telemetry. |

Begin.
