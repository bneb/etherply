export interface VoteOption {
    id: string;
    label: string;
    votes: number;
}

export interface PollState {
    question: string;
    options: Record<string, VoteOption>;
    totalVotes: number;
    winnerId?: string;
}

export const INITIAL_POLL: PollState = {
    question: "What's the best programming language?",
    options: {
        'python': { id: 'python', label: 'Python 🐍', votes: 0 },
        'typescript': { id: 'typescript', label: 'TypeScript 🔷', votes: 0 },
        'rust': { id: 'rust', label: 'Rust 🦀', votes: 0 },
        'go': { id: 'go', label: 'Go 🐹', votes: 0 },
    },
    totalVotes: 0,
};
