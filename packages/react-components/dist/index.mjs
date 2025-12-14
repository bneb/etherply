// src/components/LiveCursors.tsx
import { useEffect, useRef, useState } from "react";
import { useEtherPlyContext } from "@etherply/sdk/react";

// src/components/DefaultCursor.tsx
import { jsx, jsxs } from "react/jsx-runtime";
function DefaultCursor({ color, x, y, label }) {
  return /* @__PURE__ */ jsxs(
    "div",
    {
      style: {
        position: "absolute",
        left: 0,
        top: 0,
        transform: `translate(${x}px, ${y}px)`,
        pointerEvents: "none",
        transition: "transform 120ms linear",
        zIndex: 9999
      },
      children: [
        /* @__PURE__ */ jsx(
          "svg",
          {
            width: "24",
            height: "36",
            viewBox: "0 0 24 36",
            fill: "none",
            xmlns: "http://www.w3.org/2000/svg",
            children: /* @__PURE__ */ jsx(
              "path",
              {
                d: "M5.65376 12.3673H5.46026L5.31717 12.4976L0.500002 16.8829L0.500002 1.19138L11.7841 12.3673H5.65376Z",
                fill: color,
                stroke: "white"
              }
            )
          }
        ),
        label && /* @__PURE__ */ jsx(
          "div",
          {
            style: {
              position: "absolute",
              left: 16,
              top: 16,
              backgroundColor: color,
              color: "white",
              borderRadius: "12px",
              padding: "4px 8px",
              fontSize: "12px",
              fontWeight: 600,
              whiteSpace: "nowrap"
            },
            children: label
          }
        )
      ]
    }
  );
}

// src/components/LiveCursors.tsx
import { jsx as jsx2 } from "react/jsx-runtime";
var COLORS = ["#f87171", "#fb923c", "#fbbf24", "#a3e635", "#34d399", "#22d3ee", "#818cf8", "#e879f9"];
function getRandomColor(id) {
  let hash = 0;
  for (let i = 0; i < id.length; i++) {
    hash = id.charCodeAt(i) + ((hash << 5) - hash);
  }
  return COLORS[Math.abs(hash) % COLORS.length];
}
function LiveCursors({ renderCursor, throttleMs = 33, timeoutMs = 3e4 }) {
  const client = useEtherPlyContext();
  const [myId] = useState(() => "user-" + Math.random().toString(36).slice(2, 7));
  const [cursors, setCursors] = useState({});
  const lastUpdateRef = useRef(0);
  useEffect(() => {
    const unsub = client.onMessage((msg) => {
      if (msg.type === "op" && msg.payload.key.startsWith("presence:")) {
        setCursors((prev) => ({
          ...prev,
          [msg.payload.key]: msg.payload.value
        }));
      } else if (msg.type === "init") {
        const initial = {};
        for (const [k, v] of Object.entries(msg.data)) {
          if (k.startsWith("presence:")) {
            initial[k] = v;
          }
        }
        setCursors((prev) => ({ ...prev, ...initial }));
      }
    });
    return unsub;
  }, [client]);
  useEffect(() => {
    const handleMouseMove = (e) => {
      const now = Date.now();
      if (now - lastUpdateRef.current < throttleMs) return;
      lastUpdateRef.current = now;
      const position = {
        x: e.clientX,
        y: e.clientY,
        userId: myId,
        lastUpdate: now
      };
      client.set(`presence:${myId}`, position);
    };
    window.addEventListener("mousemove", handleMouseMove);
    return () => window.removeEventListener("mousemove", handleMouseMove);
  }, [client, myId, throttleMs]);
  const activeCursors = Object.entries(cursors).filter(([key, value]) => {
    return key !== `presence:${myId}` && value?.lastUpdate > Date.now() - timeoutMs;
  }).map(([key, value]) => {
    const data = value;
    return {
      id: data.userId || key,
      x: data.x,
      y: data.y,
      color: getRandomColor(data.userId || key)
    };
  });
  return /* @__PURE__ */ jsx2("div", { style: { position: "fixed", inset: 0, pointerEvents: "none", overflow: "hidden", zIndex: 9999 }, children: activeCursors.map((cursor) => renderCursor ? renderCursor({ ...cursor, label: cursor.id }) : /* @__PURE__ */ jsx2(
    DefaultCursor,
    {
      color: cursor.color,
      x: cursor.x,
      y: cursor.y,
      label: cursor.id
    },
    cursor.id
  )) });
}
export {
  DefaultCursor,
  LiveCursors
};
