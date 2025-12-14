# Real-time Kanban Board

A complex state management example using Drag and Drop.

## Features
- **Complex Objects**: Syncs nested arrays and objects (`columns`, `cards`).
- **Drag and Drop**: Integration with `@dnd-kit`.
- **Optimistic UI**: Instant feedback while syncing.

## Code Highlight

```typescript
function handleDragEnd(event) {
    const { active, over } = event;
    // ... calculate new state locally ...
    
    // One-line sync
    setValue(newState);
}
```

[View Source Code](https://github.com/etherply/etherply/tree/main/examples/kanban)
