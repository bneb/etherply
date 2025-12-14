
import fs from 'fs';
import path from 'path';
import matter from 'gray-matter';
import { compileMDX } from 'next-mdx-remote/rsc';
import rehypeSlug from 'rehype-slug';
import rehypeAutolinkHeadings from 'rehype-autolink-headings';

const DOCS_ROOT = path.join(process.cwd(), 'src/content/docs');

export type DocMetadata = {
    title: string;
    description?: string;
    sidebar_label?: string;
    sidebar_position?: number;
    slug: string;
};

export type DocNode = {
    name: string;
    path: string;
    isDirectory: boolean;
    metadata?: DocMetadata;
    children?: DocNode[];
};

// Helper to pretty print filenames
function prettifyName(name: string): string {
    return name
        .replace(/\.mdx?$/, '') // Remove extension
        .replace(/[-_]/g, ' ')   // Replace separators with spaces
        .replace(/\b\w/g, c => c.toUpperCase()); // Title Case
}

// Recursive function to get all docs
export function getDocsTree(dir: string = DOCS_ROOT, relativePath: string = ''): DocNode[] {
    const items = fs.readdirSync(dir, { withFileTypes: true });

    const nodes: DocNode[] = items.map((item) => {
        const fullPath = path.join(dir, item.name);
        const itemRelativePath = path.join(relativePath, item.name);

        if (item.isDirectory()) {
            const children = getDocsTree(fullPath, itemRelativePath);
            if (children.length === 0) return null;

            return {
                name: item.name,
                path: itemRelativePath,
                isDirectory: true,
                children: children,
            };
        } else if (item.name.endsWith('.md') || item.name.endsWith('.mdx')) {
            const fileContent = fs.readFileSync(fullPath, 'utf8');
            const { data } = matter(fileContent);

            return {
                name: item.name,
                path: itemRelativePath.replace(/\.mdx?$/, ''),
                isDirectory: false,
                metadata: {
                    title: data.title || prettifyName(item.name),
                    description: data.description,
                    sidebar_label: data.sidebar_label,
                    sidebar_position: data.sidebar_position,
                    slug: itemRelativePath.replace(/\.mdx?$/, ''),
                } as DocMetadata
            };
        }
        return null;
    }).filter(Boolean) as DocNode[];

    return nodes.sort((a, b) => {
        // Sort logic: directories first, then sidebar_position, then alpha
        const aPos = a.metadata?.sidebar_position ?? 999;
        const bPos = b.metadata?.sidebar_position ?? 999;
        if (aPos !== bPos) return aPos - bPos;
        return a.name.localeCompare(b.name);
    });
}
