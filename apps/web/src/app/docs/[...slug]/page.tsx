
import fs from 'fs';
import path from 'path';
import { notFound } from 'next/navigation';
import { compileMDX } from 'next-mdx-remote/rsc';
import rehypeSlug from 'rehype-slug';
import rehypeAutolinkHeadings from 'rehype-autolink-headings';

export default async function DocPage({ params }: { params: { slug: string[] } }) {
    const slug = params.slug.join('/');
    const filePath = path.join(process.cwd(), 'src/content/docs', `${slug}.md`);
    const mdxPath = path.join(process.cwd(), 'src/content/docs', `${slug}.mdx`);

    let source = '';

    if (fs.existsSync(filePath)) {
        source = fs.readFileSync(filePath, 'utf8');
    } else if (fs.existsSync(mdxPath)) {
        source = fs.readFileSync(mdxPath, 'utf8');
    } else {
        notFound();
    }

    const { content, frontmatter } = await compileMDX<{ title: string; description?: string }>({
        source,
        options: {
            parseFrontmatter: true,
            mdxOptions: {
                rehypePlugins: [rehypeSlug, [rehypeAutolinkHeadings, { behavior: 'wrap' }]]
            }
        },
    });

    return (
        <article className="min-h-screen pb-20">
            <h1 className="text-4xl font-extrabold mb-4">{frontmatter.title}</h1>
            {frontmatter.description && <p className="text-xl text-zinc-400 mb-8">{frontmatter.description}</p>}
            <hr className="border-white/10 mb-8" />
            <div className="docs-content">
                {content}
            </div>
        </article>
    );
}
