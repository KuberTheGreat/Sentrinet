"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";

interface Scan {
  id: number;
  target: string;
  port: number;
  is_open: boolean;
  duration_ms: number;
  created_at: string;
}

export default function DashboardPage() {
  const { token } = useAuth();
  const router = useRouter();
  const [scans, setScans] = useState<Scan[]>([]);
  const [target, setTarget] = useState("");
  const [startPort, setStartPort] = useState(1);
  const [endPort, setEndPort] = useState(1024);

  useEffect(() => {
    if (!token) {
      router.push("/login");
      return;
    }

    fetch("http://localhost:8080/scans", {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => res.json())
      .then(setScans)
      .catch(console.error);
  }, [token]);

  const handleScan = async (e: React.FormEvent) => {
    e.preventDefault();
    const res = await fetch("http://localhost:8080/scan", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ target, start_port: startPort, end_port: endPort }),
    });
    const data = await res.json();
    alert(`Scan completed! Found ${data.filter((r: any) => r.is_open).length} open ports.`);
  };

  return (
    <div className="p-10 max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-6">SentriNet Dashboard</h1>

      <form onSubmit={handleScan} className="flex gap-2 mb-6">
        <input
          type="text"
          placeholder="Target (e.g. scanme.nmap.org)"
          value={target}
          onChange={(e) => setTarget(e.target.value)}
          className="border p-2 rounded w-1/3 bg-gray-900 text-gray-200"
        />
        <input
          type="number"
          placeholder="Start Port"
          value={startPort}
          onChange={(e) => setStartPort(Number(e.target.value))}
          className="border p-2 rounded w-24 bg-gray-900 text-gray-200"
        />
        <input
          type="number"
          placeholder="End Port"
          value={endPort}
          onChange={(e) => setEndPort(Number(e.target.value))}
          className="border p-2 rounded w-24 bg-gray-900 text-gray-200"
        />
        <button
          type="submit"
          className="bg-blue-600 px-4 py-2 rounded hover:bg-blue-700"
        >
          Start Scan
        </button>
      </form>

      <table className="w-full border border-gray-700 text-left">
        <thead>
          <tr className="bg-gray-800 text-gray-300">
            <th className="p-2">Target</th>
            <th className="p-2">Port</th>
            <th className="p-2">Status</th>
            <th className="p-2">Duration (ms)</th>
          </tr>
        </thead>
        <tbody>
          {scans.map((scan) => (
            <tr key={scan.id} className="border-t border-gray-700">
              <td className="p-2">{scan.target}</td>
              <td className="p-2">{scan.port}</td>
              <td className="p-2">{scan.is_open ? "🟢 Open" : "🔴 Closed"}</td>
              <td className="p-2">{scan.duration_ms}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
