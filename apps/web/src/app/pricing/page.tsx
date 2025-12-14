
import { Navbar } from '@/components/navbar';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Check } from 'lucide-react';

export default function PricingPage() {
    return (
        <main className="min-h-screen bg-black text-white selection:bg-blue-500/20">
            <Navbar />

            <section className="pt-32 pb-20 px-4">
                <div className="container mx-auto text-center mb-16">
                    <h1 className="text-4xl md:text-6xl font-bold mb-6 tracking-tight">
                        Simple, unified <span className="text-blue-500 glow-text">pricing</span>.
                    </h1>
                    <p className="text-xl text-zinc-500 max-w-2xl mx-auto">
                        Start for free, scale with your growth. No hidden fees for concurrent connections.
                    </p>
                </div>

                <div className="container mx-auto grid grid-cols-1 md:grid-cols-3 gap-8 max-w-6xl">
                    {/* Free Tier */}
                    <PricingCard
                        title="Developer"
                        price="$0"
                        period="/month"
                        description="Perfect for side projects and prototypes."
                        features={[
                            "Up to 100 concurrent connections",
                            "10,000 messages / month",
                            "24h message retention",
                            "Community Support",
                            "1 Project"
                        ]}
                        cta="Start for Free"
                        ctaVariant="outline"
                    />

                    {/* Pro Tier */}
                    <PricingCard
                        title="Pro"
                        price="$49"
                        period="/month"
                        description="For growing applications with production needs."
                        features={[
                            "Up to 10,000 concurrent connections",
                            "1M messages / month",
                            "30-day message retention",
                            "Priority Email Support",
                            "Unlimited Projects",
                            "Custom Domains",
                            "SLA Guarantees"
                        ]}
                        featured={true}
                        cta="Get Started"
                        ctaVariant="default"
                    />

                    {/* Enterprise Tier */}
                    <PricingCard
                        title="Enterprise"
                        price="Custom"
                        period=""
                        description="For large-scale mission-critical deployments."
                        features={[
                            "Unlimited concurrent connections",
                            "Unlimited messages",
                            "Unlimited retention",
                            "24/7 Dedicated Support",
                            "VPC Peering",
                            "SSO / SAML",
                            "Audit Logs"
                        ]}
                        cta="Contact Sales"
                        ctaVariant="outline"
                    />
                </div>
            </section>
        </main>
    );
}

function PricingCard({
    title,
    price,
    period,
    description,
    features,
    cta,
    ctaVariant = "outline",
    featured = false
}: {
    title: string,
    price: string,
    period: string,
    description: string,
    features: string[],
    cta: string,
    ctaVariant?: "default" | "outline" | "secondary" | "ghost" | "link" | "destructive" | "glass",
    featured?: boolean
}) {
    return (
        <Card className={`bg-black border-white/10 flex flex-col ${featured ? 'border-blue-500/50 shadow-[0_0_40px_-10px_rgba(59,130,246,0.3)]' : ''}`}>
            <CardHeader>
                <CardTitle className="text-2xl font-bold text-white mb-2">{title}</CardTitle>
                <div className="flex items-baseline gap-1 mb-2">
                    <span className="text-4xl font-bold text-white">{price}</span>
                    <span className="text-zinc-500">{period}</span>
                </div>
                <CardDescription className="text-zinc-500">{description}</CardDescription>
            </CardHeader>
            <CardContent className="flex-1">
                <ul className="space-y-3">
                    {features.map((feature, i) => (
                        <li key={i} className="flex items-center gap-2 text-zinc-300">
                            <div className={`p-1 rounded-full ${featured ? 'bg-blue-500/20 text-blue-400' : 'bg-white/10 text-white'}`}>
                                <Check className="h-3 w-3" />
                            </div>
                            <span className="text-sm">{feature}</span>
                        </li>
                    ))}
                </ul>
            </CardContent>
            <CardFooter>
                <Button
                    className={`w-full ${featured ? 'bg-blue-600 hover:bg-blue-500' : ''}`}
                    variant={ctaVariant == "default" ? "default" : "outline"}
                >
                    {cta}
                </Button>
            </CardFooter>
        </Card>
    )
}
