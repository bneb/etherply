
import Link from 'next/link';
import { getDocsTree, DocNode } from '@/lib/docs';
import { Navbar } from '@/components/navbar';

// Recursive sidebar item component
function SidebarItem({ node }: { node: DocNode }) {
    if (node.isDirectory) {
        return (
            <div className="mb-4">
                <h4 className="font-semibold text-zinc-100 mb-2 capitalize px-2">{node.name.replace('-', ' ')}</h4>
                <div className="pl-0 space-y-1">
                    {node.children?.map((child) => (
                        <SidebarItem key={child.path} node={child} />
                    ))}
                </div>
            </div>
        );
    }

    return (
        <Link
            href={`/docs/${node.path}`}
            className="block px-2 py-1.5 text-sm text-zinc-400 hover:text-white hover:bg-white/5 rounded-md transition-colors"
        >
            {node.metadata?.sidebar_label || node.metadata?.title || node.name}
        </Link>
    );
}

export default function DocsLayout({ children }: { children: React.ReactNode }) {
    const docsTree = getDocsTree();

    return (
        <div className="min-h-screen bg-black">
            <Navbar />
            <div className="flex container mx-auto pt-16">
                {/* Sidebar */}
                <aside className="w-64 fixed h-[calc(100vh-4rem)] overflow-y-auto border-r border-white/10 py-8 hidden lg:block">
                    <div className="space-y-1 pr-4">
                        <Link href="/docs" className="block px-2 py-1.5 text-sm font-semibold text-white hover:bg-white/5 rounded-md mb-6">
                            Documentation Home
                        </Link>
                        {docsTree.map((node) => (
                            <SidebarItem key={node.path} node={node} />
                        ))}
                    </div>
                </aside>

                {/* Main Content */}
                <main className="flex-1 lg:pl-72 py-8 px-4 lg:px-12 max-w-none prose prose-invert prose-zinc prose-headings:font-bold prose-headings:tracking-tight prose-a:text-blue-400 prose-img:rounded-lg">
                    {children}
                </main>
            </div>
        </div>
    );
}
