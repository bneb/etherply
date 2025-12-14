export interface Task {
    id: string;
    content: string;
}

export interface Column {
    id: string;
    title: string;
    taskIds: string[];
}

export interface BoardState {
    columns: Column[];
    tasks: Record<string, Task>;
}

export const INITIAL_BOARD: BoardState = {
    columns: [
        { id: 'todo', title: 'To Do', taskIds: ['task-1', 'task-2'] },
        { id: 'doing', title: 'In Progress', taskIds: [] },
        { id: 'done', title: 'Done', taskIds: [] },
    ],
    tasks: {
        'task-1': { id: 'task-1', content: 'Install EtherPly SDK' },
        'task-2': { id: 'task-2', content: 'Build cool stuff' },
    },
};
