
"use client";

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card';
import { Badge } from 'lucide-react'; // Using icon as badge mock

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
        <div className="space-y-8">
            <div>
                <h1 className="text-3xl font-bold tracking-tight">Billing & Plans</h1>
                <p className="text-zinc-400">Upgrade your workspace to handle more scale.</p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                {plans.map((plan) => (
                    <Card
                        key={plan.id}
                        className={`bg-zinc-900/50 border-white/10 flex flex-col ${currentPlan === plan.id ? 'border-blue-500 ring-1 ring-blue-500' : ''}`}
                    >
                        <CardHeader>
                            <CardTitle className="text-xl">{plan.name}</CardTitle>
                            <div className="flex items-baseline gap-1 mt-2">
                                <span className="text-3xl font-bold">{plan.price}</span>
                                {plan.period && <span className="text-sm text-zinc-500">{plan.period}</span>}
                            </div>
                            <CardDescription>{plan.description}</CardDescription>
                        </CardHeader>
                        <CardContent className="flex-1">
                            <ul className="space-y-2 text-sm text-zinc-400">
                                {plan.features.map((feature) => (
                                    <li key={feature} className="flex items-center gap-2">
                                        <div className="h-1.5 w-1.5 rounded-full bg-blue-500" />
                                        {feature}
                                    </li>
                                ))}
                            </ul>
                        </CardContent>
                        <CardFooter>
                            <Button
                                className="w-full"
                                variant={currentPlan === plan.id ? "outline" : "default"}
                                onClick={() => setCurrentPlan(plan.id)}
                            >
                                {currentPlan === plan.id ? 'Current Plan' : 'Upgrade'}
                            </Button>
                        </CardFooter>
                    </Card>
                ))}
            </div>
        </div>
    );
}
