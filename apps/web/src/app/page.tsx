
import Link from 'next/link';
import Image from 'next/image';
import { Button } from '@/components/ui/button';
import { Navbar } from '@/components/navbar';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { ArrowRight, Database, Globe, ShieldCheck, Zap } from 'lucide-react';

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col bg-black selection:bg-blue-500/20">
      <Navbar />

      {/* Hero Section */}
      <section className="pt-40 pb-32 px-4 relative overflow-hidden">
        <div className="container mx-auto text-center relative z-10">
          <div className="inline-flex items-center rounded-full border border-white/5 bg-white/[0.02] px-3 py-1 text-sm font-medium text-zinc-400 mb-8 backdrop-blur-md">
            <span className="flex h-1.5 w-1.5 rounded-full bg-blue-500 mr-2 shadow-[0_0_10px_rgb(59,130,246)]"></span>
            Version 1.0 Now Available
          </div>

          <h1 className="text-5xl md:text-7xl font-bold tracking-tighter mb-8 text-white leading-tight">
            Stop building <span className="text-blue-500 glow-text">WebSockets</span> from scratch.
          </h1>

          <div className="mt-4 text-xl md:text-2xl text-zinc-500 max-w-4xl mx-auto mb-12 leading-relaxed font-light flex flex-col gap-2">
            <p>Stop debugging CRDT merging logic.</p>
            <p>We manage concurrent state at the edge.</p>
            <p className="text-white font-medium">Deploy your core collaboration feature in 15 minutes.</p>
          </div>

          <div className="flex flex-col sm:flex-row items-center justify-center gap-6">
            <Link href="/login">
              <Button size="lg" className="h-14 px-10 text-lg rounded-full bg-blue-600 hover:bg-blue-500 text-white shadow-[0_0_30px_-5px_rgba(37,99,235,0.4)] hover:shadow-[0_0_40px_-5px_rgba(37,99,235,0.5)] transition-all duration-300">
                Start Building
                <ArrowRight className="ml-2 h-5 w-5" />
              </Button>
            </Link>
            <Link href="/docs" className="text-zinc-400 hover:text-white transition-colors flex items-center gap-2 text-lg font-medium px-6 py-4 rounded-full hover:bg-white/5">
              Read the Docs
            </Link>
          </div>
        </div>

        {/* Ambient Glow */}
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[1000px] h-[600px] bg-blue-600/10 blur-[150px] rounded-full -z-10 opacity-50" />
      </section>

      {/* Value Props */}
      <section className="py-32 border-t border-white/5">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <FeatureCard
              icon={<Database className="h-6 w-6 text-blue-500" />}
              title="Mathematical Consistency"
              description="Idempotent state management. We handle the complex math of merging concurrent edits."
            />
            <FeatureCard
              icon={<Globe className="h-6 w-6 text-blue-500" />}
              title="50ms Global Latency"
              description="Replicated to the edge via NATS JetStream. Your users are never waiting."
            />
            <FeatureCard
              icon={<ShieldCheck className="h-6 w-6 text-blue-500" />}
              title="Sovereign Infrastructure"
              description="Run the binary in your VPC. Open-core, auditable, and enterprise-ready."
            />
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-12 border-t border-white/5">
        <div className="container mx-auto px-4 text-center">
          <div className="mb-6 flex justify-center">
            <div className="p-2 bg-white/5 rounded-lg">
              <Image src="/logo.svg" alt="EtherPly Logo" width={24} height={24} className="h-6 w-6" />
            </div>
          </div>
          <p className="text-zinc-600 text-sm">© 2025 EtherPly Inc. Built for the builders.</p>
        </div>
      </footer>
    </main>
  );
}

function FeatureCard({ icon, title, description }: { icon: React.ReactNode, title: string, description: string }) {
  return (
    <Card className="bg-transparent border-white/5 hover:border-blue-500/20 transition-all duration-500 group">
      <CardHeader>
        <div className="mb-6 p-3 bg-blue-500/10 w-fit rounded-lg group-hover:bg-blue-500/20 transition-colors">
          {icon}
        </div>
        <CardTitle className="text-xl font-medium text-white group-hover:text-blue-400 transition-colors">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <CardDescription className="text-base text-zinc-500 leading-relaxed group-hover:text-zinc-400 transition-colors">
          {description}
        </CardDescription>
      </CardContent>
    </Card>
  )
}
