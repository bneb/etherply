
import Link from 'next/link';
import Image from 'next/image';

export default function DocsIndex() {
    return (
        <div className="flex flex-col items-center justify-center min-h-[80vh] text-center max-w-5xl mx-auto">
            <h1 className="text-6xl font-extrabold tracking-tight mb-6">
                <span className="bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-blue-600">
                    Global Sync
                </span> in &lt;50ms.
            </h1>
            <p className="text-xl text-zinc-400 max-w-2xl mb-12">
                nMeshed is the open-source sync engine for modern applications.
                Stop building WebSockets from scratch.
            </p>

            <div className="flex gap-4 mb-20">
                <Link
                    href="/docs/intro"
                    className="px-8 py-3 bg-blue-600 hover:bg-blue-500 text-white font-semibold rounded-lg transition-all"
                >
                    Get Started →
                </Link>
                <Link
                    href="https://github.com/bneb/etherply"
                    target="_blank"
                    className="px-8 py-3 bg-white/10 hover:bg-white/20 text-white font-semibold rounded-lg transition-all"
                >
                    View on GitHub
                </Link>
            </div>

            {/* Feature Grid Visualization */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8 w-full max-w-6xl px-4">
                {/* Card 1: Collaborative Editor */}
                <div className="rounded-xl overflow-hidden border border-white/10 bg-white/[0.02]">
                    <div className="p-3 border-b border-white/5 text-sm font-semibold text-zinc-400">
                        Multiplayer Primitives
                    </div>
                    <div className="relative aspect-[4/3] w-full">
                        <Image src="/docs/img/view-editor.png" alt="Collaborative Editor" fill className="object-cover" />
                    </div>
                </div>

                {/* Card 2: Analytics */}
                <div className="rounded-xl overflow-hidden border border-white/10 bg-white/[0.02]">
                    <div className="p-3 border-b border-white/5 text-sm font-semibold text-zinc-400">
                        Real-time Insights
                    </div>
                    <div className="relative aspect-[4/3] w-full">
                        <Image src="/docs/img/view-analytics.png" alt="Real-time Analytics" fill className="object-cover" />
                    </div>
                </div>

                {/* Card 3: Topology */}
                <div className="rounded-xl overflow-hidden border border-white/10 bg-white/[0.02]">
                    <div className="p-3 border-b border-white/5 text-sm font-semibold text-zinc-400">
                        Global Mesh
                    </div>
                    <div className="relative aspect-[4/3] w-full">
                        <Image src="/docs/img/view-topology.png" alt="Global Infrastructure" fill className="object-cover" />
                    </div>
                </div>
            </div>

            {/* Quick Install */}
            <div className="mt-24 w-full max-w-xl">
                <div className="bg-zinc-900 border border-zinc-800 rounded-lg overflow-hidden text-left">
                    <div className="px-4 py-2 bg-zinc-800 text-xs font-bold text-zinc-400 uppercase tracking-wider">
                        Installation
                    </div>
                    <div className="p-6 font-mono text-zinc-200 flex items-center justify-between">
                        <span>npm install nmeshed</span>
                        <button className="text-zinc-500 hover:text-white transition-colors" title="Copy">
                            Click to copy
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
