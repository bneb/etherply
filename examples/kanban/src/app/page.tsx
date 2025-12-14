import { KanbanBoard } from "@/components/KanbanBoard";

export default function Home() {
  return (
    <main className="min-h-screen p-8">
      <div className="max-w-7xl mx-auto mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
          Project Roadmap
        </h1>
        <p className="text-gray-500 dark:text-gray-400 mt-2">
          Collaborative Kanban Board powered by EtherPly
        </p>
      </div>
      <KanbanBoard />
    </main>
  );
}
