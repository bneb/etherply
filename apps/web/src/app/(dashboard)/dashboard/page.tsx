
"use client";

import { useEffect, useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card';
import { ProjectService, type Project } from '@/lib/mocks';
import { Plus, Server, Activity } from 'lucide-react';
import { Input } from '@/components/ui/input';

export default function DashboardPage() {
    const [projects, setProjects] = useState<Project[]>([]);
    const [loading, setLoading] = useState(true);
    const [newProjectName, setNewProjectName] = useState('');
    const [creating, setCreating] = useState(false);

    useEffect(() => {
        loadProjects();
    }, []);

    const loadProjects = async () => {
        const data = await ProjectService.listProjects('org_1');
        setProjects(data);
        setLoading(false);
    };

    const handleCreate = async () => {
        if (!newProjectName) return;
        setCreating(true);
        await ProjectService.createProject(newProjectName, 'us-east-1');
        setNewProjectName('');
        await loadProjects();
        setCreating(false);
    };

    return (
        <div className="space-y-8">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Projects</h1>
                    <p className="text-zinc-400">Manage your synced applications and regions.</p>
                </div>
                <div className="flex gap-2">
                    <Input
                        placeholder="New Project Name"
                        className="w-48 bg-transparent border-white/10"
                        value={newProjectName}
                        onChange={(e) => setNewProjectName(e.target.value)}
                    />
                    <Button onClick={handleCreate} disabled={creating || !newProjectName} className="bg-blue-600 hover:bg-blue-500">
                        <Plus className="h-4 w-4 mr-2" />
                        Create
                    </Button>
                </div>
            </div>

            {loading ? (
                <div className="text-zinc-500">Loading projects...</div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {projects.map((project) => (
                        <Card key={project.id} className="bg-zinc-900/50 border-white/10 hover:border-blue-500/50 transition-all">
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-lg font-medium">
                                    {project.name}
                                </CardTitle>
                                <Server className={cn("h-4 w-4", project.status === 'healthy' ? "text-green-500" : "text-yellow-500")} />
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold flex items-center gap-2">
                                    {project.activeConnections.toLocaleString()}
                                    <span className="text-xs font-normal text-zinc-500">conns</span>
                                </div>
                                <p className="text-xs text-zinc-500 mt-1">Region: {project.region}</p>
                            </CardContent>
                            <CardFooter>
                                <div className="w-full flex items-center gap-2 text-xs text-zinc-500">
                                    <Activity className="h-3 w-3" />
                                    Realtime Sync Active
                                </div>
                            </CardFooter>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    );
}

function cn(...inputs: (string | undefined)[]) {
    return inputs.filter(Boolean).join(' ');
}
