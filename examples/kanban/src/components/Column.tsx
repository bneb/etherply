'use client';

import { useDroppable } from '@dnd-kit/core';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { Column as ColumnType, Task } from '../types/kanban';
import { Card } from './Card';

interface ColumnProps {
    column: ColumnType;
    tasks: Task[];
}

export function Column({ column, tasks }: ColumnProps) {
    const { setNodeRef } = useDroppable({
        id: column.id,
    });

    return (
        <div className="flex flex-col w-[85vw] md:w-80 bg-gray-50 dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 h-full max-h-full snap-center shrink-0">
            <div className="p-4 border-b border-gray-200 dark:border-gray-800 flex justify-between items-center">
                <h3 className="font-semibold text-gray-700 dark:text-gray-200">
                    {column.title}
                </h3>
                <span className="bg-gray-200 dark:bg-gray-800 text-gray-600 dark:text-gray-400 px-2 py-1 rounded text-xs font-medium">
                    {tasks.length}
                </span>
            </div>

            <div ref={setNodeRef} className="flex-1 p-2 overflow-y-auto">
                <div className="flex flex-col gap-2">
                    <SortableContext items={tasks.map((t) => t.id)} strategy={verticalListSortingStrategy}>
                        {tasks.map((task) => (
                            <Card key={task.id} task={task} />
                        ))}
                    </SortableContext>
                    {tasks.length === 0 && (
                        <div className="h-full flex items-center justify-center p-8 text-gray-400 border-2 border-dashed border-gray-200 dark:border-gray-800 rounded-lg">
                            Empty
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
