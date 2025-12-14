

export interface CursorProps {
    color: string;
    x: number;
    y: number;
    label?: string;
}

export function DefaultCursor({ color, x, y, label }: CursorProps) {
    return (
        <div
            style={{
                position: 'absolute',
                left: 0,
                top: 0,
                transform: `translate(${x}px, ${y}px)`,
                pointerEvents: 'none',
                transition: 'transform 120ms linear',
                zIndex: 9999,
            }}
        >
            <svg
                width="24"
                height="36"
                viewBox="0 0 24 36"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
            >
                <path
                    d="M5.65376 12.3673H5.46026L5.31717 12.4976L0.500002 16.8829L0.500002 1.19138L11.7841 12.3673H5.65376Z"
                    fill={color}
                    stroke="white"
                />
            </svg>
            {label && (
                <div
                    style={{
                        position: 'absolute',
                        left: 16,
                        top: 16,
                        backgroundColor: color,
                        color: 'white',
                        borderRadius: '12px',
                        padding: '4px 8px',
                        fontSize: '12px',
                        fontWeight: 600,
                        whiteSpace: 'nowrap',
                    }}
                >
                    {label}
                </div>
            )}
        </div>
    );
}
