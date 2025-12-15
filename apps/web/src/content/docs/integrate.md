# The First Mile: Integration

Integration should feel like magic. 
Whether you are building a React frontend, a Python bot, or a Go microservice, the connection handshake is identical.

## 1. Install the SDK

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<Tabs>
  <TabItem value="js" label="Node.js / React" default>
    ```bash
    npm install nmeshed
    ```
  </TabItem>
  <TabItem value="python" label="Python">
    ```bash
    pip install nmeshed
    ```
  </TabItem>
  <TabItem value="go" label="Go">
    ```bash
    go get github.com/bneb/etherply/pkg/go-sdk
    ```
  </TabItem>
</Tabs>

## 2. Initialize the Client

Configure your connection. The `token` is your passport; the `workspaceId` is your room.

<Tabs>
  <TabItem value="js" label="JavaScript / TS" default>
    ```typescript
    import { NMeshedClient } from 'nmeshed';

    const client = new NMeshedClient({
      workspaceId: 'my-room',
      token: 'sk_live_...', // Your secret token
    });
    
    await client.connect();
    ```
  </TabItem>
  <TabItem value="react" label="React Hook">
    ```tsx
    import { useNMeshed } from 'nmeshed/react';

    function App() {
      const { status } = useNMeshed({
        workspaceId: 'my-room',
        token: 'sk_live_...'
      });
      
      return <div>Status: {status}</div>;
    }
    ```
  </TabItem>
  <TabItem value="python" label="Python">
    ```python
    from nmeshed import NMeshedClient

    client = NMeshedClient(
        workspace_id="my-room",
        token="sk_live_..."
    )
    await client.connect()
    ```
  </TabItem>
  <TabItem value="go" label="Go">
    ```go
    c := client.New(client.Config{
        WorkspaceID: "my-room",
        Token:       "sk_live_...",
    })
    c.Connect(ctx)
    ```
  </TabItem>
</Tabs>

## 3. Synchronize State

Don't think about "messages." Think about "State."
When you set `key=value`, it propagates to every connected peer in milliseconds.

<Tabs>
  <TabItem value="js" label="JavaScript" default>
    ```typescript
    // Write
    client.set('status', 'online');

    // Read (Reactive)
    client.onMessage(msg => {
      if (msg.type === 'op') console.log(msg.payload);
    });
    ```
  </TabItem>
  <TabItem value="react" label="React (Magic)" default>
    ```tsx
    const { value, setValue } = useDocument({ 
      key: 'status', 
      initialValue: 'offline' 
    });

    // Write: automatically syncs + optimistic update
    setValue('online');
    ```
  </TabItem>
</Tabs>

## Next Steps

You are now connected to the nMeshed mesh.
- **[Build a Text Editor](/docs/examples/text-editor)**: The "Hello World" of multiplayer.
- **[Add Live Cursors](/docs/examples/cursors)**: Add presence in 30 seconds.
