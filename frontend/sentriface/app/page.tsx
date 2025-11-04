"use client";
import Link from "next/link";

export default function HomePage() {
  return (
    <main className="flex flex-col items-center justify-center min-h-screen">
      <h1 className="text-5xl font-bold mb-6">👁️‍🗨️ SentriNet</h1>
      <p className="text-lg mb-10 text-gray-400">Secure. Smart. Self-evolving.</p>
      <div className="flex gap-4">
        <Link href="/login" className="bg-blue-600 px-6 py-2 rounded-md hover:bg-blue-700">
          Login
        </Link>
        <Link href="/register" className="bg-gray-700 px-6 py-2 rounded-md hover:bg-gray-800">
          Register
        </Link>
      </div>
    </main>
  );
}
