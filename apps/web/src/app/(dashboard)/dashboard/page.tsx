"use client";

import { useEffect, useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { api, type Project } from '@/lib/api';
import { Plus, Server, Activity, Globe, Wifi } from 'lucide-react';
import { Input } from '@/components/ui/input';

export default function DashboardPage() {
    const [projects, setProjects] = useState<Project[]>([]);
    const [loading, setLoading] = useState(true);
    const [newProjectName, setNewProjectName] = useState('');
    const [creating, setCreating] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        loadProjects();
    }, []);

    const loadProjects = async () => {
        try {
            const data = await api.listProjects();
            setProjects(data);
        } catch (e) {
            console.error(e);
            setError("Failed to load projects. Ensure backend is running.");
        } finally {
            setLoading(false);
        }
    };

    const handleCreate = async () => {
        if (!newProjectName) return;
        setCreating(true);
        try {
            await api.createProject(newProjectName, 'us-east-1');
            setNewProjectName('');
            await loadProjects();
        } catch (e) {
            setError("Failed to create project.");
        } finally {
            setCreating(false);
        }
    };

    return (
        <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight text-white mb-1">Projects</h1>
                    <p className="text-zinc-500">Manage your synced applications and regions.</p>
                </div>
                <div className="flex gap-3">
                    <Input
                        placeholder="New Project Name"
                        className="w-full md:w-64 bg-black/20 border-white/10 text-white placeholder:text-zinc-600 focus-visible:ring-blue-500/30"
                        value={newProjectName}
                        onChange={(e) => setNewProjectName(e.target.value)}
                    />
                    <Button onClick={handleCreate} disabled={creating || !newProjectName} className="bg-blue-600 hover:bg-blue-500 text-white shadow-[0_0_20px_-5px_rgba(37,99,235,0.4)]">
                        <Plus className="h-4 w-4 mr-2" />
                        {creating ? 'Creating...' : 'Create'}
                    </Button>
                </div>
            </div>

            {error && (
                <div className="p-4 bg-red-500/10 border border-red-500/20 text-red-500 rounded-md">
                    {error}
                </div>
            )}

            {loading ? (
                <div className="text-zinc-500 flex items-center gap-2">
                    <div className="h-4 w-4 animate-spin rounded-full border-2 border-white/10 border-t-blue-500"></div>
                    Loading projects...
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {projects.map((project) => (
                        <Card key={project.id} className="group bg-white/[0.02] border-white/5 hover:border-blue-500/30 hover:bg-white/[0.04] transition-all duration-300">
                            <CardHeader className="flex flex-row items-start justify-between space-y-0 pb-2">
                                <div>
                                    <CardTitle className="text-lg font-medium text-white group-hover:text-blue-100 transition-colors">
                                        {project.name}
                                    </CardTitle>
                                    <div className="flex items-center gap-2 mt-2">
                                        <Badge variant="secondary" className="text-[10px] px-1.5 py-0 h-5 font-normal tracking-wide uppercase bg-white/5 text-zinc-500 border border-white/5">
                                            {project.region}
                                        </Badge>
                                        <Badge variant="default" className="text-[10px] px-1.5 py-0 h-5 font-normal tracking-wide uppercase border border-white/5 bg-emerald-500/10 text-emerald-500 hover:bg-emerald-500/20">
                                            Healthy
                                        </Badge>
                                    </div>
                                </div>
                                <div className="p-2 bg-white/5 rounded-lg group-hover:bg-blue-500/10 transition-colors">
                                    <Server className="h-4 w-4 text-zinc-500 group-hover:text-blue-400 transition-colors" />
                                </div>
                            </CardHeader>
                            <CardContent className="mt-4">
                                <div className="flex items-end gap-2">
                                    <span className="text-3xl font-bold text-white tracking-tight">0</span>
                                    <span className="text-sm text-zinc-500 mb-1.5">active conns</span>
                                </div>
                            </CardContent>
                            <CardFooter className="pt-0">
                                <div className="w-full flex items-center justify-between text-xs text-zinc-500 border-t border-white/5 pt-4 mt-2">
                                    <div className="flex items-center gap-1.5">
                                        <Activity className="h-3 w-3 text-emerald-500" />
                                        <span className="text-zinc-400">Syncing (12ms)</span>
                                    </div>
                                    <div className="flex items-center gap-1.5">
                                        <Wifi className="h-3 w-3" />
                                        ws://api.etherply.io
                                    </div>
                                </div>
                            </CardFooter>
                        </Card>
                    ))}

                    {/* Add New Project Placeholder Card */}
                    <div
                        onClick={() => document.querySelector<HTMLInputElement>('input[placeholder="New Project Name"]')?.focus()}
                        className="cursor-pointer rounded-xl border border-dashed border-white/10 bg-white/[0.01] hover:bg-white/[0.03] hover:border-white/20 transition-all duration-300 flex flex-col items-center justify-center p-6 gap-3 text-zinc-600 hover:text-zinc-400 min-h-[200px]"
                    >
                        <div className="p-3 bg-white/5 rounded-full">
                            <Plus className="h-6 w-6" />
                        </div>
                        <span className="text-sm font-medium">Create New Project</span>
                    </div>
                </div>
            )}
        </div>
    );
}
