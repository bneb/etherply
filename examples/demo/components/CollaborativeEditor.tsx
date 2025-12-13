"use client";

import { useEffect, useState, useRef, useLayoutEffect } from 'react';
import { ConnectionStatus, EtherPlyClient } from '@/lib/etherply-client';
import { useDebounce } from './hooks/useDebounce';

export default function CollaborativeEditor({ client }: { client: EtherPlyClient }) {
    const [content, setContent] = useState("");
    const [status, setStatus] = useState<ConnectionStatus>('IDLE');
    const [isInitialized, setIsInitialized] = useState(false);
    const textareaRef = useRef<HTMLTextAreaElement>(null);
    const cursorRef = useRef<{ start: number; end: number } | null>(null);

    // Debounced sender to prevent network flooding (300ms)
    const sendOperationDebounced = useDebounce((key: string, value: string) => {
        client.sendOperation(key, value);
    }, 300);

    // Cursor Preservation: Restore cursor position synchronously after render
    useLayoutEffect(() => {
        const input = textareaRef.current;
        if (input && cursorRef.current && document.activeElement === input) {
            input.setSelectionRange(cursorRef.current.start, cursorRef.current.end);
            cursorRef.current = null; // Clear after restore
        }
    }, [content]);

    useEffect(() => {
        // Listen for incoming ops
        const unsubscribe = client.onMessage((msg) => {
            if (msg.type === 'init') {
                const text = msg.data["document_text"]?.value || "";
                setContent(text);
                setIsInitialized(true);
            } else if (msg.type === 'op') {
                const { key, value } = msg.payload;
                if (key === 'document_text') {
                    setContent((prev) => {
                        if (prev !== value) {
                            // Capture cursor before update if focused
                            const input = textareaRef.current;
                            if (input && document.activeElement === input) {
                                cursorRef.current = {
                                    start: input.selectionStart,
                                    end: input.selectionEnd,
                                };
                            }
                            return value;
                        }
                        return prev;
                    });
                }
            }
        });

        return () => {
            unsubscribe();
        };
    }, [client]);

    useEffect(() => {
        // Subscribe to connection status
        const unsubStatus = client.onStatusChange((s) => {
            setStatus(s);
            if (s === 'CONNECTED' && !isInitialized) {
                // Connected, waiting for init message
            }
        });
        client.connect();
        return () => unsubStatus();
    }, [client, isInitialized]);

    const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        const newVal = e.target.value;
        setContent(newVal);
        // Send to server (debounced)
        sendOperationDebounced("document_text", newVal);
    };

    return (
        <div className="w-full relative min-h-[400px] border border-gray-200 rounded-lg shadow-sm focus-within:ring-2 focus-within:ring-brand-500 transition-all bg-white flex flex-col overflow-hidden">
            {/* Loading / Connection Overlay */}
            {(status === 'CONNECTING' || status === 'IDLE') && (
                <div className="absolute inset-0 bg-white/80 backdrop-blur-sm z-20 flex items-center justify-center">
                    <div className="flex flex-col items-center space-y-3">
                        <div className="w-6 h-6 border-2 border-brand-500 border-t-transparent rounded-full animate-spin"></div>
                        <span className="text-sm text-gray-500 font-medium">Connecting to EtherPly...</span>
                    </div>
                </div>
            )}

            {/* Error State */}
            {status === 'ERROR' && (
                <div className="absolute top-4 right-4 bg-red-50 border border-red-200 text-red-600 px-3 py-1.5 rounded-md shadow-sm z-50 text-xs font-medium flex items-center gap-2">
                    <span className="w-2 h-2 rounded-full bg-red-500 animate-pulse"></span>
                    Connection Lost. Retrying...
                </div>
            )}

            <textarea
                ref={textareaRef}
                value={content}
                onChange={handleChange}
                placeholder="Start typing..."
                className="w-full h-full p-6 resize-none outline-none text-gray-800 font-mono text-sm leading-relaxed bg-transparent"
                disabled={status !== 'CONNECTED'}
                spellCheck={false}
            />
        </div>
    );
}
