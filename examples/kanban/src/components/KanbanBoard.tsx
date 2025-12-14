'use client';

import {
    DndContext,
    DragOverlay,
    useSensor,
    useSensors,
    PointerSensor,
    closestCorners,
    DragStartEvent,
    DragOverEvent,
    DragEndEvent,
} from '@dnd-kit/core';
import { useDocument, useEtherPlyContext } from '@etherply/sdk/react';
import { useState, useEffect, useMemo } from 'react';
import { BoardState, INITIAL_BOARD } from '../types/kanban';
import { Column } from './Column';
import { Card } from './Card';
import { arrayMove } from '@dnd-kit/sortable';

export function KanbanBoard() {
    const client = useEtherPlyContext();
    const [isConnected, setIsConnected] = useState(false);
    const [activeId, setActiveId] = useState<string | null>(null);

    useEffect(() => {
        return client.onStatusChange((status) => {
            setIsConnected(status === 'CONNECTED');
        });
    }, [client]);

    const { value: board, setValue } = useDocument<BoardState>({
        key: 'kanban-board',
        initialValue: INITIAL_BOARD,
    });

    const sensors = useSensors(
        useSensor(PointerSensor, {
            activationConstraint: {
                distance: 5,
            },
        })
    );

    const activeTask = useMemo(() =>
        activeId && board?.tasks[activeId] ? board.tasks[activeId] : null
        , [activeId, board?.tasks]);

    // ... inside component ...

    const tasksToColumns = useMemo(() => {
        const map = new Map<string, string>();
        if (!board) return map;
        for (const column of board.columns) {
            for (const taskId of column.taskIds) {
                map.set(taskId, column.id);
            }
        }
        return map;
    }, [board]);

    function findColumn(id: string) {
        if (!board) return null;
        if (tasksToColumns.has(id)) {
            return tasksToColumns.get(id);
        }
        // Fallback: Check if the ID itself is a column ID
        return board.columns.find((col) => col.id === id)?.id || null;
    }

    function handleDragStart(event: DragStartEvent) {
        setActiveId(event.active.id as string);
    }

    function handleDragOver(event: DragOverEvent) {
        const { active, over } = event;
        const overId = over?.id;

        if (!overId || active.id === overId || !board) return;

        const activeColumnId = findColumn(active.id as string);
        const overColumnId = findColumn(overId as string);

        if (!activeColumnId || !overColumnId || activeColumnId === overColumnId) {
            return;
        }

        // Moving between columns
        const activeColumn = board.columns.find((c) => c.id === activeColumnId)!;
        const overColumn = board.columns.find((c) => c.id === overColumnId)!;

        const activeItems = activeColumn.taskIds;
        const overItems = overColumn.taskIds;

        const activeIndex = activeItems.indexOf(active.id as string);
        const overIndex = overItems.indexOf(overId as string);

        let newIndex;
        if (overItems.includes(overId as string)) {
            newIndex = overItems.indexOf(overId as string);
        } else {
            newIndex = overItems.length + 1;
        }

        // Clone state
        const newBoard = {
            ...board,
            columns: board.columns.map((col) => {
                if (col.id === activeColumnId) {
                    return {
                        ...col,
                        taskIds: col.taskIds.filter((id) => id !== active.id),
                    };
                }
                if (col.id === overColumnId) {
                    return {
                        ...col,
                        taskIds: [
                            ...col.taskIds.slice(0, newIndex),
                            active.id as string,
                            ...col.taskIds.slice(newIndex, col.taskIds.length),
                        ],
                    };
                }
                return col;
            }),
        };

        setValue(newBoard);
    }

    function handleDragEnd(event: DragEndEvent) {
        const { active, over } = event;
        const activeId = active.id as string;
        const overId = over?.id as string;

        if (!overId || !board) {
            setActiveId(null);
            return;
        }

        const activeColumnId = findColumn(activeId);
        const overColumnId = findColumn(overId);

        if (activeColumnId && overColumnId && activeColumnId === overColumnId) {
            const columnIndex = board.columns.findIndex((col) => col.id === activeColumnId);
            const column = board.columns[columnIndex];
            const oldIndex = column.taskIds.indexOf(activeId);
            const newIndex = column.taskIds.indexOf(overId);

            if (oldIndex !== newIndex) {
                const newBoard = {
                    ...board,
                    columns: board.columns.map((col, index) => {
                        if (index === columnIndex) {
                            return {
                                ...col,
                                taskIds: arrayMove(col.taskIds, oldIndex, newIndex),
                            };
                        }
                        return col;
                    })
                };
                setValue(newBoard);
            }
        }

        setActiveId(null);
    }

    if (!board) return null;

    return (
        <DndContext
            sensors={sensors}
            collisionDetection={closestCorners}
            onDragStart={handleDragStart}
            onDragOver={handleDragOver}
            onDragEnd={handleDragEnd}
        >
            <div className="flex justify-start md:justify-center p-2 md:p-8 overflow-x-auto snap-x snap-mandatory">
                <div className="flex gap-4 md:gap-6 h-[80vh] min-w-full px-4 md:px-0">
                    {board.columns.map((col) => (
                        <Column
                            key={col.id}
                            column={col}
                            tasks={col.taskIds.map((id) => board.tasks[id])}
                        />
                    ))}
                </div>
            </div>

            <div className="fixed bottom-4 right-4 flex items-center gap-2 bg-white dark:bg-gray-800 p-2 rounded-full shadow-lg border border-gray-200 dark:border-gray-700">
                <div
                    className={`w-3 h-3 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'
                        }`}
                />
                <span className="text-sm font-medium pr-2">
                    {isConnected ? 'Connected' : 'Reconnecting...'}
                </span>
            </div>

            <DragOverlay>
                {activeTask ? <Card task={activeTask} /> : null}
            </DragOverlay>
        </DndContext>
    );
}
