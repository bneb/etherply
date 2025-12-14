'use client';

import { useDocument, useEtherPlyContext } from '@etherply/sdk/react';
import { motion } from 'framer-motion';
import { useState, useEffect } from 'react';
import { INITIAL_POLL, PollState } from '../types/poll';
import { clsx } from 'clsx';

export function VotingApp() {
    const [hasVoted, setHasVoted] = useState(false);
    const client = useEtherPlyContext();
    const [isConnected, setIsConnected] = useState(false);

    useEffect(() => {
        return client.onStatusChange((status) => {
            setIsConnected(status === 'CONNECTED');
        });
    }, [client]);

    const { value: poll, setValue, isLoaded } = useDocument<PollState>({
        key: 'poll-data',
        initialValue: INITIAL_POLL,
    });

    const handleVote = (optionId: string) => {
        if (!poll) return;

        // Optimistic update
        const newPoll = { ...poll };
        newPoll.options[optionId].votes += 1;
        // Total votes will be recalculated by the Python bot!
        // But we can optimistically increment it too for better UX
        newPoll.totalVotes += 1;

        setValue(newPoll);
        setHasVoted(true);
    };

    const sortedOptions = Object.values(poll?.options || {}).sort((a, b) => b.votes - a.votes);
    const total = poll?.totalVotes || 0;

    if (!isLoaded) {
        return (
            <div className="flex items-center justify-center min-h-[400px]">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-500" />
            </div>
        );
    }

    return (
        <div className="w-full max-w-md mx-auto bg-white dark:bg-gray-800 rounded-2xl shadow-xl overflow-hidden border border-gray-100 dark:border-gray-700">
            <div className="p-8">
                <div className="flex justify-between items-start mb-6">
                    <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
                        {poll?.question}
                    </h2>
                    <div className={`w-3 h-3 rounded-full mt-2 ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} title={isConnected ? 'Live' : 'Disconnect'} />
                </div>

                <div className="space-y-4">
                    {sortedOptions.map((option) => {
                        const percentage = total > 0 ? Math.round((option.votes / total) * 100) : 0;
                        const isWinner = poll?.winnerId === option.id;

                        return (
                            <motion.button
                                layout
                                key={option.id}
                                onClick={() => handleVote(option.id)}
                                disabled={!isConnected}
                                className={clsx(
                                    "relative w-full text-left p-4 rounded-xl border-2 transition-colors",
                                    isWinner ? "border-amber-400 bg-amber-50 dark:bg-amber-900/20" : "border-gray-100 dark:border-gray-700 hover:border-indigo-200 dark:hover:border-indigo-800",
                                    "group overflow-hidden"
                                )}
                            >
                                {/* Progress Bar Background */}
                                <motion.div
                                    className="absolute inset-0 bg-indigo-50 dark:bg-indigo-900/30 origin-left z-0"
                                    initial={{ scaleX: 0 }}
                                    animate={{ scaleX: percentage / 100 }}
                                    transition={{ duration: 0.5, ease: "easeOut" }}
                                />

                                <div className="relative z-10 flex justify-between items-center">
                                    <span className="font-medium text-lg text-gray-900 dark:text-gray-100">
                                        {option.label}
                                        {isWinner && <span className="ml-2 text-amber-500">👑</span>}
                                    </span>
                                    <div className="flex flex-col items-end">
                                        <span className="font-bold text-gray-900 dark:text-white">{percentage}%</span>
                                        <span className="text-sm text-gray-500 dark:text-gray-400">{option.votes} votes</span>
                                    </div>
                                </div>
                            </motion.button>
                        );
                    })}
                </div>

                <div className="mt-6 text-center text-sm text-gray-500 dark:text-gray-400">
                    {total} votes cast • {isConnected ? 'Real-time' : 'Reconnecting...'}
                </div>
            </div>

            {/* Bot Status Banner */}
            <div className="bg-gray-50 dark:bg-gray-900/50 p-4 text-xs text-center border-t border-gray-100 dark:border-gray-800">
                <span className="font-mono text-indigo-500">Python Bot</span> is calculating totals & winners...
            </div>
        </div>
    );
}
