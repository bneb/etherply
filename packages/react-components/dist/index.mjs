var __defProp = Object.defineProperty;
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, { get: all[name], enumerable: true });
};

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

// src/components/Button/Button.tsx
import React from "react";
import { jsx as jsx3, jsxs as jsxs2 } from "react/jsx-runtime";
var Button = React.forwardRef(
  ({ className = "", variant = "primary", size = "md", isLoading, children, ...props }, ref) => {
    const baseStyles = "inline-flex items-center justify-center font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none";
    const variants = {
      primary: "bg-[#0094c6] text-white hover:bg-[#005e7c]",
      secondary: "bg-[#001242] text-white hover:bg-[#000022]",
      ghost: "hover:bg-gray-100 hover:text-gray-900 dark:hover:bg-gray-800 dark:hover:text-gray-50",
      destructive: "bg-red-500 text-white hover:bg-red-600"
    };
    const sizes = {
      sm: "h-8 px-3 text-xs",
      md: "h-10 px-4 py-2 text-sm",
      lg: "h-12 px-8 text-md"
    };
    const rounded = "rounded-md";
    const combinedClassName = `
      ${baseStyles} 
      ${variants[variant]} 
      ${sizes[size]} 
      ${rounded} 
      ${className}
    `.trim().replace(/\s+/g, " ");
    return /* @__PURE__ */ jsxs2(
      "button",
      {
        ref,
        className: combinedClassName,
        disabled: isLoading || props.disabled,
        ...props,
        children: [
          isLoading && /* @__PURE__ */ jsxs2("svg", { className: "animate-spin -ml-1 mr-2 h-4 w-4 text-current", xmlns: "http://www.w3.org/2000/svg", fill: "none", viewBox: "0 0 24 24", children: [
            /* @__PURE__ */ jsx3("circle", { className: "opacity-25", cx: "12", cy: "12", r: "10", stroke: "currentColor", strokeWidth: "4" }),
            /* @__PURE__ */ jsx3("path", { className: "opacity-75", fill: "currentColor", d: "M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" })
          ] }),
          children
        ]
      }
    );
  }
);
Button.displayName = "Button";

// src/tokens/colors.ts
var colors_exports = {};
__export(colors_exports, {
  colors: () => colors
});
var colors = {
  primary: {
    DEFAULT: "#0094c6",
    // Ocean Blue
    hover: "#005e7c",
    // Baltic Blue
    foreground: "#ffffff"
  },
  secondary: {
    DEFAULT: "#001242",
    // Deep Navy
    hover: "#000022",
    // Prussian Blue
    foreground: "#ffffff"
  },
  destructive: {
    DEFAULT: "#ef4444",
    hover: "#dc2626",
    foreground: "#ffffff"
  },
  surface: {
    DEFAULT: "#ffffff",
    subtle: "#f3f4f6",
    // gray-100
    dark: "#1f2937"
    // gray-800
  },
  border: {
    DEFAULT: "#e5e7eb",
    // gray-200
    dark: "#374151"
    // gray-700
  }
};
export {
  Button,
  DefaultCursor,
  LiveCursors,
  colors_exports as tokens
};
