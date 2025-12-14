import type { Metadata } from "next";
import { Inter } from "next/font/google"; // Using Inter as requested standards
import "./globals.css";
import { cn } from "@/lib/utils";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "EtherPly - Postgres for Realtime",
  description: "The sync engine for professional SaaS applications.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body className={cn(inter.className, "min-h-screen bg-background antialiased selection:bg-blue-500/30")}>
        <div className="fixed inset-0 bg-[url('/grid.svg')] bg-center [mask-image:linear-gradient(180deg,white,rgba(255,255,255,0))]" />
        <div className="relative">
          {children}
        </div>
      </body>
    </html>
  );
}

