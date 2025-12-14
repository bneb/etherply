
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Navbar } from '@/components/navbar';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { ArrowRight, CheckCircle2, Database, Globe, ShieldCheck } from 'lucide-react';

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col">
      <Navbar />

      {/* Hero Section */}
      <section className="pt-32 pb-20 px-4 md:pt-48 md:pb-32 relative overflow-hidden">
        <div className="container mx-auto text-center relative z-10">
          <div className="inline-flex items-center rounded-full border border-blue-500/30 bg-blue-500/10 px-3 py-1 text-sm font-medium text-blue-400 mb-8 backdrop-blur-sm">
            <span className="flex h-2 w-2 rounded-full bg-blue-500 mr-2 animate-pulse"></span>
            Version 1.0 Now Available
          </div>

          <h1 className="text-5xl md:text-7xl font-bold tracking-tight mb-6 bg-clip-text text-transparent bg-gradient-to-b from-white to-white/40">
            Postgres for <span className="text-blue-500">Realtime</span>.
          </h1>

          <p className="mt-4 text-xl text-zinc-400 max-w-2xl mx-auto mb-10 leading-relaxed">
            The sync engine for professional SaaS. Add Googledocs-style collaboration to your app correctly, with conflict-free merging and on-premise control.
          </p>

          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link href="/login">
              <Button size="lg" className="h-12 px-8 text-base bg-blue-600 hover:bg-blue-500 shadow-blue-500/25 shadow-lg">
                Start Building
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </Link>
            <Link href="/docs" className="text-zinc-400 hover:text-white transition-colors flex items-center gap-2">
              Read the Docs
            </Link>
          </div>
        </div>

        {/* Background Glow */}
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[500px] bg-blue-500/20 blur-[120px] rounded-full -z-10" />
      </section>

      {/* Value Props */}
      <section className="py-24 bg-black/40 border-y border-white/5">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <FeatureCard
              icon={<Database className="h-8 w-8 text-blue-400" />}
              title="Conflict-Free State"
              description="Built on Automerge CRDTs. We handle the complex math of merging concurrent edits so you don't have to."
            />
            <FeatureCard
              icon={<Globe className="h-8 w-8 text-indigo-400" />}
              title="Global Replication"
              description="Multi-region by default. Your data is replicated to the edge via NATS JetStream for <50ms latency."
            />
            <FeatureCard
              icon={<ShieldCheck className="h-8 w-8 text-emerald-400" />}
              title="On-Premise Ready"
              description="Don't trust the cloud? Run the EtherPly binary in your own VPC. We are open-core and enterprise ready."
            />
          </div>
        </div>
      </section>

      {/* Footer Stub */}
      <footer className="py-12 border-t border-white/10 bg-black">
        <div className="container mx-auto px-4 text-center text-zinc-600">
          <p>© 2025 EtherPly Inc. Built for the builders.</p>
        </div>
      </footer>
    </main>
  );
}

function FeatureCard({ icon, title, description }: { icon: React.ReactNode, title: string, description: string }) {
  return (
    <Card className="bg-white/5 border-white/10 hover:border-blue-500/50 transition-colors">
      <CardHeader>
        <div className="mb-4">{icon}</div>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <CardDescription className="text-base text-zinc-400">
          {description}
        </CardDescription>
      </CardContent>
    </Card>
  )
}
