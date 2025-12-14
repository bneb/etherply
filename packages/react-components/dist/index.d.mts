import * as react_jsx_runtime from 'react/jsx-runtime';

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

export { type CursorProps, DefaultCursor, LiveCursors, type LiveCursorsProps };
