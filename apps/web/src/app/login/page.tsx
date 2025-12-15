
"use client";

import { useState } from 'react';
import Image from 'next/image';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { api } from '@/lib/api';
import { useRouter } from 'next/navigation';
import { ArrowRight } from 'lucide-react';
import Link from 'next/link';

export default function LoginPage() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [loading, setLoading] = useState(false);
    const router = useRouter();

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);

        try {
            await api.login(email, password);
            router.push('/dashboard');
        } catch (error) {
            console.error('Login failed:', error);
            // In a real app we would show a toast/alert here
            // For MVP velocity we just log it, as per plan
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-black p-4 selection:bg-blue-500/20">
            {/* Background Effects */}
            <div className="fixed inset-0 bg-[url('/grid.svg')] bg-center [mask-image:linear-gradient(180deg,white,rgba(255,255,255,0))]" />
            <div className="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[500px] h-[500px] bg-blue-600/10 blur-[100px] rounded-full -z-10" />

            <Card className="w-full max-w-md bg-black/40 backdrop-blur-xl border-white/5 shadow-2xl relative z-10">
                <CardHeader className="text-center pb-8">
                    <div className="mx-auto mb-6 bg-blue-600/10 p-3 rounded-xl ring-1 ring-blue-500/20 w-fit cursor-pointer hover:bg-blue-600/20 transition-colors">
                        <Link href="/">
                            <Image src="/logo.svg" alt="nMeshed Logo" width={24} height={24} className="h-6 w-6" />
                        </Link>
                    </div>
                    <CardTitle className="text-2xl font-bold tracking-tight text-white">Welcome Back</CardTitle>
                    <CardDescription className="text-zinc-500">Enter your email to sign in to the console.</CardDescription>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleLogin} className="space-y-4">
                        <div className="space-y-2">
                            <Input
                                type="email"
                                placeholder="developer@nmeshed.com"
                                required
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                className="h-11 bg-white/[0.03] border-white/10 text-white placeholder:text-zinc-600 focus-visible:ring-blue-500/50"
                            />
                        </div>
                        <div className="space-y-2">
                            <Input
                                type="password"
                                placeholder="••••••••"
                                required
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                className="h-11 bg-white/[0.03] border-white/10 text-white placeholder:text-zinc-600 focus-visible:ring-blue-500/50"
                            />
                        </div>
                        <Button type="submit" className="w-full h-11 bg-blue-600 hover:bg-blue-500 text-white shadow-[0_0_20px_-5px_rgba(37,99,235,0.4)]" disabled={loading}>
                            {loading ? 'Signing In...' : 'Continue with Email'}
                            {!loading && <ArrowRight className="ml-2 h-4 w-4" />}
                        </Button>
                    </form>
                </CardContent>
                <CardFooter className="flex justify-center pt-6 pb-6">
                    <p className="text-xs text-zinc-600 text-center max-w-xs">
                        By clicking continue, you agree to our <Link href="#" className="underline hover:text-zinc-400 transition-colors">Terms of Service</Link> and <Link href="#" className="underline hover:text-zinc-400 transition-colors">Privacy Policy</Link>.
                    </p>
                </CardFooter>
            </Card>
        </div>
    );
}
