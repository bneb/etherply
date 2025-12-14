import * as react_jsx_runtime from 'react/jsx-runtime';
import React$1 from 'react';

interface LiveCursorsProps {
    /**
     * Optional custom cursor component.
     */
    renderCursor?: (props: {
        x: number;
        y: number;
        color: string;
        label: string;
        id: string;
    }) => React.ReactNode;
    /**
     * Throttling interval in milliseconds (default: 33ms ~30fps).
     */
    throttleMs?: number;
    /**
     * Timeout in milliseconds to remove stale cursors (default: 30000ms).
     */
    timeoutMs?: number;
}
declare function LiveCursors({ renderCursor, throttleMs, timeoutMs }: LiveCursorsProps): react_jsx_runtime.JSX.Element;

interface CursorProps {
    color: string;
    x: number;
    y: number;
    label?: string;
}
declare function DefaultCursor({ color, x, y, label }: CursorProps): react_jsx_runtime.JSX.Element;

interface ButtonProps extends React$1.ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: 'primary' | 'secondary' | 'ghost' | 'destructive';
    size?: 'sm' | 'md' | 'lg';
    isLoading?: boolean;
}
declare const Button: React$1.ForwardRefExoticComponent<ButtonProps & React$1.RefAttributes<HTMLButtonElement>>;

declare const colors: {
    readonly primary: {
        readonly DEFAULT: "#0094c6";
        readonly hover: "#005e7c";
        readonly foreground: "#ffffff";
    };
    readonly secondary: {
        readonly DEFAULT: "#001242";
        readonly hover: "#000022";
        readonly foreground: "#ffffff";
    };
    readonly destructive: {
        readonly DEFAULT: "#ef4444";
        readonly hover: "#dc2626";
        readonly foreground: "#ffffff";
    };
    readonly surface: {
        readonly DEFAULT: "#ffffff";
        readonly subtle: "#f3f4f6";
        readonly dark: "#1f2937";
    };
    readonly border: {
        readonly DEFAULT: "#e5e7eb";
        readonly dark: "#374151";
    };
};

declare const colors$1_colors: typeof colors;
declare namespace colors$1 {
  export { colors$1_colors as colors };
}

export { Button, type CursorProps, DefaultCursor, LiveCursors, type LiveCursorsProps, colors$1 as tokens };
