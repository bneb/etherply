"use client";

import { useEffect, useState, useRef } from 'react';
import { ConnectionStatus, EtherPlyClient } from '@/lib/etherply-client';

export default function CollaborativeEditor({ client }: { client: EtherPlyClient }) {
    const [content, setContent] = useState("");
    const [status, setStatus] = useState<ConnectionStatus>('IDLE');
    const [isInitialized, setIsInitialized] = useState(false);
    const textareaRef = useRef<HTMLTextAreaElement>(null);

    useEffect(() => {
        // Listen for incoming ops
        const unsubscribe = client.onMessage((msg) => {
            if (msg.type === 'init') {
                // Apply full state
                // In our LWW Engine, we get a map of keys.
                // Let's assume the text content is stored under key "document_text"
                const text = msg.data["document_text"]?.value || "";
                setContent(text);
            } else if (msg.type === 'op') {
                const { key, value } = msg.payload;
                if (key === 'document_text') {
                    // If this op is from us, we might ignore if we are optimistic, 
                    // but for simple MVP LWW, valid to just apply (or check timestamp).
                    // Here we just set it.
                    // In a real editor, we need cursor preservation.
                    // We'll do a simple check to avoid cursor jump if content matches.
                    // Cursor Preservation Logic
                    const input = textareaRef.current;
                    let cursorStart = 0;
                    let cursorEnd = 0;

                    if (input && document.activeElement === input) {
                        cursorStart = input.selectionStart;
                        cursorEnd = input.selectionEnd;
                    }

                    setContent(prev => {
                        if (prev !== value) {
                            // If content matches, no update needed (and no cursor jump risk)
                            // If content changes, React will re-render.
                            // We must restore cursor AFTER render.
                            if (input && document.activeElement === input) {
                                requestAnimationFrame(() => {
                                    // Naive cursor restoration: keep position.
                                    // (Real CRDT cursor requires transforming index, but MVP LWW just keeps index)
                                    // If text was inserted BEFORE cursor, this drifts. 
                                    // But strictly better than jumping to end.
                                    input.setSelectionRange(cursorStart, cursorEnd);
                                });
                            }
                            return value;
                        }
                        return prev;
                    });
                }
            }
        });

        return () => {
            if (unsubscribe) unsubscribe();
        };

        return () => {
            if (unsubscribe) unsubscribe();
        };
    }, [client]);

    useEffect(() => {
        // Subscribe to connection status
        const unsubStatus = client.onStatusChange((s) => {
            setStatus(s);
            if (s === 'CONNECTED' && !isInitialized) {
                // We wait for first message for init, but connected is good.
            }
        });
        client.connect();
        return () => unsubStatus();
    }, [client]);

    const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        const newVal = e.target.value;
        setContent(newVal);
        // Send to server
        client.sendOperation("document_text", newVal);
    };

    return (
        <div className="w-full h-96 border border-gray-200 rounded-lg shadow-sm focus-within:ring-2 focus-within:ring-brand-500 transition-all p-4 bg-white relative">
            {/* Loading Skeleton */}
            {(status === 'CONNECTING' || status === 'IDLE') && (
                <div className="absolute inset-0 bg-white/50 backdrop-blur-sm z-10 flex items-center justify-center rounded-lg">
                    <div className="w-full h-full p-4 space-y-4 animate-pulse">
                        <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                        <div className="h-4 bg-gray-200 rounded w-1/2"></div>
                        <div className="h-4 bg-gray-200 rounded w-5/6"></div>
                    </div>
                </div>
            )}

            {/* Error State */}
            {status === 'ERROR' && (
                <div className="absolute top-4 right-4 bg-red-100 border border-red-400 text-red-700 px-4 py-2 rounded shadow-md z-50 text-sm">
                    Connection Lost. Retrying...
                </div>
            )}

            <textarea
                ref={textareaRef}
                value={content}
                onChange={handleChange}
                placeholder={status === 'CONNECTED' ? "Start typing to initialize workspace..." : "Connecting..."}
                className="w-full h-full resize-none outline-none text-gray-800 font-mono disabled:opacity-50"
                disabled={status !== 'CONNECTED'}
            />
        </div>
    );
}
