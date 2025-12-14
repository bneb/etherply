
"use client";

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { AuthService } from '@/lib/mocks';
import { useRouter } from 'next/navigation';
import { Zap } from 'lucide-react';
import Link from 'next/link';

export default function LoginPage() {
    const [email, setEmail] = useState('');
    const [loading, setLoading] = useState(false);
    const router = useRouter();

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);

        try {
            const user = await AuthService.signIn(email);
            if (typeof window !== 'undefined') {
                localStorage.setItem('etherply_user', JSON.stringify(user));
            }
            router.push('/dashboard');
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-black/90 p-4">
            <Card className="w-full max-w-md bg-white/5 border-white/10">
                <CardHeader className="text-center">
                    <div className="mx-auto mb-4 bg-blue-600 p-2 rounded-lg w-fit">
                        <Zap className="h-6 w-6 text-white text-fill-white" />
                    </div>
                    <CardTitle className="text-2xl">Welcome Back</CardTitle>
                    <CardDescription>Enter your email to sign in to the console.</CardDescription>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleLogin} className="space-y-4">
                        <div className="space-y-2">
                            <Input
                                type="email"
                                placeholder="developer@company.com"
                                required
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                className="bg-black/20 border-white/10 text-white placeholder:text-zinc-500"
                            />
                        </div>
                        <Button type="submit" className="w-full" disabled={loading}>
                            {loading ? 'Signing In...' : 'Sign In with Email'}
                        </Button>
                    </form>
                </CardContent>
                <CardFooter className="flex justify-center border-t border-white/5 pt-6">
                    <p className="text-xs text-zinc-500">
                        By clicking continue, you agree to our <Link href="#" className="underline hover:text-white">Terms of Service</Link>
                    </p>
                </CardFooter>
            </Card>
        </div>
    );
}
