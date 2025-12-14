import { CursorOverlay } from "@/components/CursorOverlay";

export default function Home() {
  return (
    <main className="min-h-screen flex flex-col items-center justify-center relative overflow-hidden">
      <CursorOverlay />

      <div className="z-10 text-center max-w-2xl px-4 pointer-events-none select-none">
        <div className="inline-block animate-bounce mb-8">
          <span className="text-6xl">👇</span>
        </div>
        <h1 className="text-5xl md:text-7xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-500 to-purple-600 mb-6 drop-shadow-sm">
          Move Your Mouse
        </h1>
        <p className="text-xl md:text-2xl text-gray-600 dark:text-gray-300 mb-8">
          Open this page in a second window to see multiplayer cursors in action.
        </p>
        <div className="inline-flex items-center gap-2 px-4 py-2 bg-white/50 dark:bg-black/50 backdrop-blur rounded-full border border-gray-200 dark:border-gray-800 text-sm text-gray-500">
          <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse" />
          Powered by EtherPly Presence
        </div>
      </div>

      {/* Decorative background grid */}
      <div className="absolute inset-0 -z-10 bg-[linear-gradient(to_right,#80808012_1px,transparent_1px),linear-gradient(to_bottom,#80808012_1px,transparent_1px)] bg-[size:24px_24px]" />
    </main>
  );
}
