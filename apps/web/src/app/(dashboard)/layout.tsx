
"use client";

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { cn } from '@/lib/utils';
import { LayoutDashboard, CreditCard, Settings, LogOut, Zap } from 'lucide-react';

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
        <div className="flex min-h-screen bg-black text-white">
            {/* Sidebar */}
            <aside className="w-64 border-r border-white/10 bg-zinc-900/50 hidden md:flex flex-col">
                <div className="p-6 border-b border-white/10 flex items-center gap-2">
                    <div className="p-1 bg-blue-600 rounded">
                        <Zap className="h-4 w-4 text-white fill-white" />
                    </div>
                    <span className="font-bold text-lg">EtherPly</span>
                </div>

                <div className="flex-1 py-6 px-4 space-y-1">
                    {navItems.map((item) => {
                        const Icon = item.icon;
                        const isActive = pathname === item.href;

                        return (
                            <Link
                                key={item.href}
                                href={item.href}
                                className={cn(
                                    "flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors",
                                    isActive
                                        ? "bg-blue-600/10 text-blue-400"
                                        : "text-zinc-400 hover:text-white hover:bg-white/5"
                                )}
                            >
                                <Icon className="h-4 w-4" />
                                {item.label}
                            </Link>
                        )
                    })}
                </div>

                <div className="p-4 border-t border-white/10">
                    <Link href="/" className="flex items-center gap-3 px-3 py-2 text-sm font-medium text-zinc-400 hover:text-red-400 transition-colors">
                        <LogOut className="h-4 w-4" />
                        Sign Out
                    </Link>
                </div>
            </aside>

            {/* Main Content */}
            <main className="flex-1 overflow-auto">
                <header className="h-16 border-b border-white/10 flex items-center justify-between px-8 bg-zinc-900/50 md:hidden">
                    <span className="font-bold">EtherPly Console</span>
                    {/* Mobile Menu Toggle would go here */}
                </header>
                <div className="p-8 max-w-6xl mx-auto">
                    {children}
                </div>
            </main>
        </div>
    );
}
