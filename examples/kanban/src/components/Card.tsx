'use client';

import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { Task } from '../types/kanban';

interface CardProps {
    task: Task;
}

export function Card({ task }: CardProps) {
    const {
        attributes,
        listeners,
        setNodeRef,
        transform,
        transition,
        isDragging,
    } = useSortable({ id: task.id });

    const style = {
        transform: CSS.Transform.toString(transform),
        transition,
    };

    if (isDragging) {
        return (
            <div
                ref={setNodeRef}
                style={style}
                className="bg-gray-100 dark:bg-gray-700 p-4 rounded-lg shadow-lg opacity-50 h-24 border-2 border-blue-500"
            />
        );
    }

    return (
        <div
            ref={setNodeRef}
            style={style}
            {...attributes}
            {...listeners}
            className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-sm hover:shadow-md cursor-grab active:cursor-grabbing border border-gray-200 dark:border-gray-700 select-none group touch-none"
        >
            <div className="font-medium text-gray-900 dark:text-gray-100">
                {task.content}
            </div>
        </div>
    );
}
