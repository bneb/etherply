
import Link from 'next/link';
import Image from 'next/image';
import { Button } from '@/components/ui/button';
import { Github } from 'lucide-react';

export function Navbar() {
    return (
        <nav className="fixed top-0 left-0 right-0 z-50 border-b border-white/10 bg-black/50 backdrop-blur-lg">
            <div className="container mx-auto px-4 h-16 flex items-center justify-between">
                <Link href="/" className="flex items-center gap-2">
                    <div className="flex items-center justify-center">
                        <Image src="/logo.svg" alt="EtherPly Logo" width={24} height={24} className="h-6 w-6" />
                    </div>
                    <span className="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-white to-white/70">
                        EtherPly
                    </span>
                </Link>

                <div className="flex items-center gap-8">
                    <div className="hidden md:flex items-center gap-6 text-sm font-medium text-zinc-400">
                        <Link href="/docs" className="hover:text-white transition-colors">Documentation</Link>
                        <Link href="/pricing" className="hover:text-white transition-colors">Pricing</Link>
                    </div>

                    <div className="flex items-center gap-4">
                        <Link href="https://github.com/bneb/etherply" target="_blank">
                            <Button variant="ghost" size="icon" className="text-zinc-400 hover:text-white">
                                <Github className="h-5 w-5" />
                            </Button>
                        </Link>
                        <Link href="/login">
                            <Button variant="glass" className="font-semibold">
                                Console
                            </Button>
                        </Link>
                    </div>
                </div>
            </div>
        </nav>
    );
}
