
"use client";

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card';
import { Check } from 'lucide-react';

export default function BillingPage() {
    const [currentPlan, setCurrentPlan] = useState('free');

    // Plans data
    const plans = [
        {
            id: 'free',
            name: 'Hobby',
            price: '$0',
            description: 'For passion projects',
            features: ['100 Concurrent Connections', '24h History', 'Community Support']
        },
        {
            id: 'scale_up',
            name: 'Scale-Up',
            price: '$499',
            period: '/mo',
            popular: true,
            description: 'For growing SaaS',
            features: ['10k Concurrent Connections', '30d History', 'Email Support', 'Multi-Region']
        },
        {
            id: 'enterprise',
            name: 'Enterprise',
            price: 'Custom',
            description: 'For mission critical apps',
            features: ['Unlimited Connections', 'Unlimited History', 'Dedicated Slack', 'VPC Peering']
        }
    ];

    return (
        <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
            <div>
                <h1 className="text-3xl font-bold tracking-tight text-white mb-1">Billing & Plans</h1>
                <p className="text-zinc-500">Upgrade your workspace to handle more scale.</p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                {plans.map((plan) => {
                    const isSelected = currentPlan === plan.id;
                    const isPopular = plan.popular;

                    return (
                        <Card
                            key={plan.id}
                            className={`
                                relative flex flex-col transition-all duration-300
                                ${isSelected
                                    ? 'bg-blue-600/[0.03] border-blue-500/50 shadow-[0_0_40px_-10px_rgba(37,99,235,0.2)]'
                                    : 'bg-white/[0.02] border-white/5 hover:border-white/10 hover:bg-white/[0.03]'
                                }
                            `}
                        >
                            {/* Popular Badge */}
                            {isPopular && !isSelected && (
                                <div className="absolute -top-3 left-1/2 -translate-x-1/2 bg-blue-600 text-white text-[10px] font-bold uppercase tracking-widest py-1 px-3 rounded-full border border-blue-400 shadow-lg shadow-blue-900/50">
                                    Most Popular
                                </div>
                            )}

                            {/* Active Badge */}
                            {isSelected && (
                                <div className="absolute -top-3 left-1/2 -translate-x-1/2 bg-emerald-500/10 text-emerald-400 text-[10px] font-bold uppercase tracking-widest py-1 px-3 rounded-full border border-emerald-500/20 backdrop-blur-sm">
                                    Current Plan
                                </div>
                            )}

                            <CardHeader>
                                <CardTitle className={`text-xl ${isSelected ? 'text-blue-400' : 'text-white'}`}>{plan.name}</CardTitle>
                                <div className="flex items-baseline gap-1 mt-4">
                                    <span className="text-4xl font-bold text-white tracking-tight">{plan.price}</span>
                                    {plan.period && <span className="text-sm text-zinc-500">{plan.period}</span>}
                                </div>
                                <CardDescription className="text-zinc-500">{plan.description}</CardDescription>
                            </CardHeader>
                            <CardContent className="flex-1">
                                <ul className="space-y-4 text-sm">
                                    {plan.features.map((feature) => (
                                        <li key={feature} className="flex items-start gap-3 text-zinc-400">
                                            <div className={`mt-0.5 rounded-full p-0.5 ${isSelected ? 'bg-blue-500/20 text-blue-400' : 'bg-zinc-800 text-zinc-500'}`}>
                                                <Check className="h-3 w-3" />
                                            </div>
                                            {feature}
                                        </li>
                                    ))}
                                </ul>
                            </CardContent>
                            <CardFooter>
                                <Button
                                    className={`w-full h-10 ${isSelected ? 'border-blue-500/30 text-blue-400 hover:text-blue-300 hover:bg-blue-500/10' : ''}`}
                                    variant={isSelected ? "outline" : "default"}
                                    onClick={() => setCurrentPlan(plan.id)}
                                >
                                    {isSelected ? 'Your Plan' : 'Upgrade Plan'}
                                </Button>
                            </CardFooter>
                        </Card>
                    )
                })}
            </div>
        </div>
    );
}
