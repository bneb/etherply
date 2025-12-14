import { VotingApp } from "@/components/VotingApp";

export default function Home() {
  return (
    <main className="min-h-screen bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8 flex flex-col items-center justify-center">
      <div className="text-center mb-12">
        <h1 className="text-4xl font-extrabold text-gray-900 dark:text-white sm:text-5xl sm:tracking-tight lg:text-6xl">
          Start Voting
        </h1>
        <p className="mt-5 max-w-xl mx-auto text-xl text-gray-500 dark:text-gray-400">
          Real-time cross-language demo.
          <br />
          <span className="text-indigo-600 dark:text-indigo-400 font-medium">Next.js</span> frontend + <span className="text-yellow-600 dark:text-yellow-400 font-medium">Python</span> backend.
        </p>
      </div>

      <VotingApp />

      <div className="mt-12 grid grid-cols-1 gap-8 sm:grid-cols-2 max-w-2xl w-full">
        <div className="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-sm border border-gray-200 dark:border-gray-700">
          <h3 className="text-lg font-medium text-gray-900 dark:text-white">Frontend (TypeScript)</h3>
          <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
            Optimistically updates local state and sends vote ops via <code className="bg-gray-100 dark:bg-gray-900 px-1 py-0.5 rounded">@etherply/sdk</code>.
          </p>
        </div>
        <div className="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-sm border border-gray-200 dark:border-gray-700">
          <h3 className="text-lg font-medium text-gray-900 dark:text-white">Backend (Python)</h3>
          <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
            Observes votes, calculates aggregates, and updates the winner using <code className="bg-gray-100 dark:bg-gray-900 px-1 py-0.5 rounded">etherply-python</code>.
          </p>
        </div>
      </div>
    </main>
  );
}
