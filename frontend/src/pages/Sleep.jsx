import { useState, useEffect } from "react";
import Header from "../components/Header";
import useAuth from "../hooks/useAuth";

const Sleep = () => {
  const { token } = useAuth();
  const [logs, setLogs] = useState([]);
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const [durationHours, setDurationHours] = useState("");
  const [qualityRating, setQualityRating] = useState(5);

  const fetchData = async () => {
    try {
      const [logsRes, statsRes] = await Promise.all([
        fetch("/api/sleep", { headers: { Authorization: `Bearer ${token}` } }),
        fetch("/api/sleep/stats", {
          headers: { Authorization: `Bearer ${token}` },
        }),
      ]);

      if (!logsRes.ok || !statsRes.ok) {
        throw new Error("Failed to fetch sleep data");
      }

      const logsData = await logsRes.json();
      const statsData = await statsRes.json();

      setLogs(logsData || []);
      setStats(statsData);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (token) {
      fetchData();
    }
  }, [token]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch("/api/sleep", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          duration_hours: parseFloat(durationHours),
          quality_rating: parseInt(qualityRating),
        }),
      });

      if (response.ok) {
        setDurationHours("");
        setQualityRating(5);
        fetchData();
      } else {
        throw new Error("Failed to add sleep log");
      }
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div className="bg-gray-900 text-white min-h-screen">
      <Header />
      <div className="p-4">
        <h1 className="text-2xl font-bold mb-4">Sleep</h1>
        {loading && <p>Loading...</p>}
        {error && <p className="text-red-500">{error}</p>}
        {token ? (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div>
              <h2 className="text-xl font-bold mb-4">Log Your Sleep</h2>
              <form
                onSubmit={handleSubmit}
                className="bg-gray-800 p-6 rounded-lg shadow-md"
              >
                <div className="mb-4">
                  <label
                    className="block text-gray-400 text-sm font-bold mb-2"
                    htmlFor="durationHours"
                  >
                    Duration (hours)
                  </label>
                  <input
                    id="durationHours"
                    type="number"
                    step="0.1"
                    value={durationHours}
                    onChange={(e) => setDurationHours(e.target.value)}
                    className="shadow appearance-none border rounded w-full py-2 px-3 bg-gray-700 text-white leading-tight focus:outline-none focus:shadow-outline"
                    required
                  />
                </div>
                <div className="mb-6">
                  <label
                    className="block text-gray-400 text-sm font-bold mb-2"
                    htmlFor="qualityRating"
                  >
                    Quality Rating (1-10)
                  </label>
                  <input
                    id="qualityRating"
                    type="range"
                    min="1"
                    max="10"
                    value={qualityRating}
                    onChange={(e) => setQualityRating(e.target.value)}
                    className="w-full"
                  />
                  <div className="text-center">{qualityRating}</div>
                </div>
                <button
                  type="submit"
                  className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full"
                >
                  Add Log
                </button>
              </form>

              {stats && (
                <div className="bg-gray-800 p-6 rounded-lg shadow-md mt-8">
                  <h2 className="text-xl font-bold mb-4">
                    Your Stats (Last 30 Days)
                  </h2>
                  <p>
                    Average Sleep Duration: {stats.average_duration.toFixed(2)}{" "}
                    hours
                  </p>
                  <p>
                    Average Sleep Quality: {stats.average_quality.toFixed(2)} /
                    10
                  </p>
                </div>
              )}
            </div>

            <div>
              <h2 className="text-xl font-bold mb-4">Your Sleep Logs</h2>
              <div className="bg-gray-800 p-6 rounded-lg shadow-md">
                {logs.length > 0 ? (
                  <ul>
                    {logs.map((log) => (
                      <li
                        key={log.ID}
                        className="border-b border-gray-700 py-2"
                      >
                        <p className="font-bold">
                          Duration: {log.DurationHours} hours
                        </p>
                        <p>Quality: {log.QualityRating}/10</p>
                        <p className="text-sm text-gray-400">
                          {new Date(log.CreatedAt).toLocaleString()}
                        </p>
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p>No sleep logs yet.</p>
                )}
              </div>
            </div>
          </div>
        ) : (
          <p>
            Please{" "}
            <a href="/login" className="text-blue-500 hover:underline">
              login
            </a>{" "}
            to view your sleep data.
          </p>
        )}
      </div>
    </div>
  );
};

export default Sleep;
