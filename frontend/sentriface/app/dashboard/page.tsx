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
  const { token, authFetch, logout } = useAuth();
  const router = useRouter();

  const [scans, setScans] = useState<Scan[]>([]);
  const [target, setTarget] = useState("");
  const [startPort, setStartPort] = useState(1);
  const [endPort, setEndPort] = useState(1024);
  const [loading, setLoading] = useState(false);

  // filters
  const [searchTerm, setSearchTerm] = useState("");
  const [showOpenOnly, setShowOpenOnly] = useState(false);
  const [sortOrder, setSortOrder] = useState<"newest" | "oldest">("newest");

  useEffect(() => {
    if (!token) {
      router.push("/login");
      return;
    }
    fetchScans();
  }, [token]);

  const fetchScans = async () => {
    try {
      const res = await authFetch("https://sentrinet.onrender.com/scans");
      const data = await res.json();
      if (Array.isArray(data)) setScans(data);
      else if (Array.isArray(data.scans)) setScans(data.scans);
    } catch (err) {
      console.error("Error fetching scans:", err);
    }
  };

  const handleScan = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!target) return alert("Enter a target to scan!");

    setLoading(true);
    try {
      const res = await authFetch("https://sentrinet.onrender.com/scan", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ target, start_port: startPort, end_port: endPort }),
      });

      const data = await res.json();
      const openCount = data.filter((r: any) => r.is_open).length;
      alert(`Scan completed! Found ${openCount} open ports.`);
      setTarget("");
      await fetchScans();
    } catch (err) {
      console.error(err);
      alert("Scan failed. Check console for details.");
    } finally {
      setLoading(false);
    }
  };

  // derived view
  const filtered = scans
    .filter((scan) =>
      scan.target.toLowerCase().includes(searchTerm.toLowerCase())
    )
    .filter((scan) => (showOpenOnly ? scan.is_open : true))
    .sort((a, b) =>
      sortOrder === "newest"
        ? new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
        : new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
    );

  return (
    <div className="min-h-screen bg-gray-950 text-white p-8">
      <div className="max-w-5xl mx-auto space-y-8">
        {/* HEADER */}
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold text-blue-400">SentriNet Dashboard</h1>
          <button
            onClick={logout}
            className="text-sm text-red-400 hover:text-red-500 transition"
          >
            Logout
          </button>
        </div>

        {/* SCAN FORM */}
        <form
          onSubmit={handleScan}
          className="bg-gray-900 p-6 rounded-xl shadow-md space-y-4"
        >
          <h2 className="text-xl font-semibold">Run a New Scan</h2>
          <div className="flex flex-wrap gap-3">
            <input
              type="text"
              placeholder="Target (e.g. scanme.nmap.org)"
              value={target}
              onChange={(e) => setTarget(e.target.value)}
              className="flex-1 min-w-[200px] p-2 rounded-md bg-gray-800 border border-gray-700 text-white focus:ring focus:ring-blue-600"
            />
            <input
              type="number"
              placeholder="Start Port"
              value={startPort}
              onChange={(e) => setStartPort(Number(e.target.value))}
              className="w-28 p-2 rounded-md bg-gray-800 border border-gray-700 text-white focus:ring focus:ring-blue-600"
            />
            <input
              type="number"
              placeholder="End Port"
              value={endPort}
              onChange={(e) => setEndPort(Number(e.target.value))}
              className="w-28 p-2 rounded-md bg-gray-800 border border-gray-700 text-white focus:ring focus:ring-blue-600"
            />
            <button
              type="submit"
              disabled={loading}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-md font-semibold disabled:bg-gray-600"
            >
              {loading ? "Scanning..." : "Start Scan"}
            </button>
          </div>
        </form>

        {/* FILTERS */}
        <div className="flex flex-wrap gap-3 items-center bg-gray-900 p-4 rounded-xl shadow-md">
          <input
            type="text"
            placeholder="Search target..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="flex-grow p-2 bg-gray-800 border border-gray-700 rounded-md focus:ring focus:ring-blue-600"
          />

          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={showOpenOnly}
              onChange={() => setShowOpenOnly((p) => !p)}
            />
            <span>Show open ports only</span>
          </label>

          <select
            value={sortOrder}
            onChange={(e) => setSortOrder(e.target.value as "newest" | "oldest")}
            className="bg-gray-800 border border-gray-700 rounded-md p-2"
          >
            <option value="newest">Newest first</option>
            <option value="oldest">Oldest first</option>
          </select>
        </div>

        {/* SCAN TABLE */}
        <div className="bg-gray-900 p-6 rounded-xl shadow-lg overflow-x-auto">
          <h2 className="text-xl font-semibold mb-4 text-blue-300">
            Scan History
          </h2>

          {filtered.length === 0 ? (
            <p className="text-gray-500 text-center py-6">
              No scans found — run one above or adjust filters.
            </p>
          ) : (
            <table className="w-full text-left border border-gray-800">
              <thead>
                <tr className="bg-gray-800 text-gray-300">
                  <th className="p-2">Target</th>
                  <th className="p-2">Port</th>
                  <th className="p-2">Status</th>
                  <th className="p-2">Duration (ms)</th>
                  <th className="p-2">Date</th>
                </tr>
              </thead>
              <tbody>
                {filtered.map((scan) => (
                  <tr
                    key={scan.id}
                    className="border-t border-gray-800 hover:bg-gray-800/50 transition"
                  >
                    <td className="p-2 font-medium">{scan.target}</td>
                    <td className="p-2">{scan.port}</td>
                    <td className="p-2">
                      {scan.is_open ? (
                        <span className="text-green-400">🟢 Open</span>
                      ) : (
                        <span className="text-red-500">🔴 Closed</span>
                      )}
                    </td>
                    <td className="p-2">{scan.duration_ms}</td>
                    <td className="p-2 text-gray-400 text-sm">
                      {new Date(scan.created_at).toLocaleString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  );
}
