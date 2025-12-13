"use client";

import { useEffect, useState } from "react";
import CollaborativeEditor from "@/components/CollaborativeEditor";
import PresenceWidget from "@/components/PresenceWidget";
import { EtherPlyClient } from "@/lib/etherply-client";

export default function Home() {
    const [client, setClient] = useState<EtherPlyClient | null>(null);

    useEffect(() => {
        const init = async () => {
            // Generate a random user ID for demo
            const userId = "user_" + Math.floor(Math.random() * 10000);

            try {
                // Fetch secure token from our own backend (Next.js server)
                const res = await fetch(`/api/auth/token?userId=${userId}`);
                if (!res.ok) {
                    console.error("Failed to fetch auth token:", res.statusText);
                    return;
                }
                const data = await res.json();

                if (data.token) {
                    const c = new EtherPlyClient({
                        workspaceId: "demo-workspace",
                        userId: userId,
                        token: data.token
                    });
                    setClient(c);
                }
            } catch (err) {
                console.error("Auth setup failed:", err);
            }
        };

        init();
    }, []);

    if (!client) return null;

    return (
        <main className="flex min-h-screen flex-col items-center justify-between p-24 bg-brand-50">
            <div className="z-10 max-w-5xl w-full items-center justify-between font-mono text-sm lg:flex mb-12">
                <p className="fixed left-0 top-0 flex w-full justify-center border-b border-gray-300 bg-gradient-to-b from-zinc-200 pb-6 pt-8 backdrop-blur-2xl dark:border-neutral-800 dark:bg-zinc-800/30 dark:from-inherit lg:static lg:w-auto  lg:rounded-xl lg:border lg:bg-gray-200 lg:p-4 lg:dark:bg-zinc-800/30">
                    <span className="font-bold mr-2 tracking-tight">EtherPly</span> <span className="font-light italic">The Synchronized World</span>
                </p>
                <div className="fixed bottom-0 left-0 flex h-48 w-full items-end justify-center bg-gradient-to-t from-white via-white dark:from-black dark:via-black lg:static lg:h-auto lg:w-auto lg:bg-none">
                    <PresenceWidget client={client} />
                </div>
            </div>

            <div className="relative flex place-items-center w-full max-w-3xl">
                <div className="w-full">
                    <h2 className="text-3xl font-bold mb-4 text-brand-900 tracking-tight">Don&apos;t Just Watch. Be Present.</h2>
                    <p className="mb-8 text-lg text-gray-600 leading-relaxed max-w-lg">
                        A collaborative canvas that breathes with you. Open this page in a second window to feel the pulse of real-time state.
                    </p>
                    <CollaborativeEditor client={client} />
                </div>
            </div>

            <div className="mb-32 grid text-center lg:max-w-5xl lg:w-full lg:mb-0 lg:grid-cols-4 lg:text-left mt-24">
                <a
                    href="https://etherply.com/docs"
                    className="group rounded-lg border border-transparent px-5 py-4 transition-colors hover:border-gray-300 hover:bg-gray-100 hover:dark:border-neutral-700 hover:dark:bg-neutral-800/30"
                    target="_blank"
                    rel="noopener noreferrer"
                >
                    <h2 className={`mb-3 text-2xl font-semibold tracking-tight`}>
                        The Ritual{" "}
                        <span className="inline-block transition-transform group-hover:translate-x-1 motion-reduce:transform-none">
                            -&gt;
                        </span>
                    </h2>
                    <p className={`m-0 max-w-[30ch] text-sm opacity-60`}>
                        Five minutes to multiplayer. Read the liturgy.
                    </p>
                </a>

                <a
                    href="https://etherply.com/pricing"
                    className="group rounded-lg border border-transparent px-5 py-4 transition-colors hover:border-gray-300 hover:bg-gray-100 hover:dark:border-neutral-700 hover:dark:bg-neutral-800/30"
                    target="_blank"
                    rel="noopener noreferrer"
                >
                    <h2 className={`mb-3 text-2xl font-semibold tracking-tight`}>
                        Economics{" "}
                        <span className="inline-block transition-transform group-hover:translate-x-1 motion-reduce:transform-none">
                            -&gt;
                        </span>
                    </h2>
                    <p className={`m-0 max-w-[30ch] text-sm opacity-60`}>
                        Scale without ceiling. From garage to IPO.
                    </p>
                </a>
            </div>
        </main>
    );
}
