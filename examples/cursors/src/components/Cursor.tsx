import React from 'react';

interface CursorProps {
    color: string;
    x: number;
    y: number;
    label?: string;
}

export function Cursor({ color, x, y, label }: CursorProps) {
    return (
        <div
            className="absolute top-0 left-0 pointer-events-none transition-transform duration-100 ease-linear z-50"
            style={{
                transform: `translate(${x}px, ${y}px)`,
            }}
        >
            <svg
                className="w-5 h-5 -mt-1 -ml-1"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
            >
                <path
                    d="M3 3L10.07 19.97L12.58 12.58L19.97 10.07L3 3Z"
                    fill={color}
                    stroke="white"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                />
            </svg>
            {label && (
                <div
                    className="absolute left-4 top-4 px-2 py-1 rounded-full text-xs font-semibold text-white whitespace-nowrap"
                    style={{ backgroundColor: color }}
                >
                    {label}
                </div>
            )}
        </div>
    );
}
