import { useState, useRef, useEffect, useLayoutEffect, RefObject } from 'react';
import { ConnectionStatus, EtherPlyClient } from '@/lib/etherply-client';
import { useDebounce } from './useDebounce';

type UseCollaborativeEditorReturn = {
    content: string;
    status: ConnectionStatus;
    handleChange: (e: React.ChangeEvent<HTMLTextAreaElement>) => void;
};

export function useCollaborativeEditor(
    client: EtherPlyClient,
    textareaRef: RefObject<HTMLTextAreaElement>
): UseCollaborativeEditorReturn {
    const [content, setContent] = useState("");
    const [status, setStatus] = useState<ConnectionStatus>('IDLE');
    const [isInitialized, setIsInitialized] = useState(false);
    const cursorRef = useRef<{ start: number; end: number } | null>(null);

    // Debounced sender to prevent network flooding (300ms)
    // NOTE: In a production app, we would use a proper CRDT JSON patch.
    // Here we send full text (LWW limit).
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
    }, [content, textareaRef]);

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
    }, [client, textareaRef]);

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

    return { content, status, handleChange };
}
