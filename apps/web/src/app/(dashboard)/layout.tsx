
"use client";

import Link from 'next/link';
import Image from 'next/image';
import { usePathname } from 'next/navigation';
import { cn } from '@/lib/utils';
import { LayoutDashboard, CreditCard, Settings, LogOut } from 'lucide-react';

export default function DashboardLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    const pathname = usePathname();

    const navItems = [
        { href: '/dashboard', label: 'Projects', icon: LayoutDashboard },
        { href: '/billing', label: 'Billing', icon: CreditCard },
        { href: '/settings', label: 'Settings', icon: Settings },
    ];

    return (
        <div className="flex min-h-screen bg-black text-foreground font-sans selection:bg-blue-500/20">
            {/* Sidebar */}
            <aside className="w-64 flex flex-col fixed inset-y-0 left-0 border-r border-white/5 bg-black/50 backdrop-blur-xl z-50 hidden md:flex">
                <div className="p-8">
                    <div className="flex items-center gap-3 mb-10">
                        <div className="flex items-center justify-center">
                            <Image src="/logo.svg" alt="nMeshed Logo" width={24} height={24} className="h-6 w-6" />
                        </div>
                        <span className="font-bold text-lg tracking-tight text-white">nMeshed</span>
                    </div>

                    <nav className="space-y-1">
                        {navItems.map((item) => {
                            const Icon = item.icon;
                            const isActive = pathname === item.href;

                            return (
                                <Link
                                    key={item.href}
                                    href={item.href}
                                    className={cn(
                                        "flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-300",
                                        isActive
                                            ? "text-white bg-white/[0.03] border border-white/5 shadow-inner"
                                            : "text-zinc-500 hover:text-white hover:bg-white/[0.02]"
                                    )}
                                >
                                    <Icon className={cn("h-4 w-4 transition-colors", isActive ? "text-blue-500" : "text-zinc-600 group-hover:text-zinc-400")} />
                                    {item.label}
                                </Link>
                            )
                        })}
                    </nav>
                </div>

                <div className="mt-auto p-8 border-t border-white/5">
                    <Link href="/" className="flex items-center gap-3 text-sm font-medium text-zinc-600 hover:text-red-400 transition-colors">
                        <LogOut className="h-4 w-4" />
                        Sign Out
                    </Link>
                </div>
            </aside>

            {/* Main Content */}
            <main className="flex-1 md:ml-64 min-h-screen">
                <header className="h-16 border-b border-white/5 flex items-center justify-between px-6 bg-black/50 backdrop-blur-md sticky top-0 z-40 md:hidden">
                    <div className="flex items-center gap-2">
                        <div className="flex items-center justify-center">
                            <Image src="/logo.svg" alt="nMeshed Logo" width={24} height={24} className="h-6 w-6" />
                        </div>
                        <span className="font-bold text-white">nMeshed</span>
                    </div>
                </header>
                <div className="p-8 max-w-7xl mx-auto animate-in fade-in duration-500">
                    {children}
                </div>
            </main>
        </div>
    );
}
