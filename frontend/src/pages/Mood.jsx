import { useState, useEffect } from "react";
import Header from "../components/Header";
import useAuth from "../hooks/useAuth";

const Mood = () => {
  const { token } = useAuth();
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const [moodRating, setMoodRating] = useState(5);
  const [notes, setNotes] = useState("");

  const fetchData = async () => {
    try {
      const response = await fetch("/api/mood", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!response.ok) {
        throw new Error("Failed to fetch mood data");
      }
      const data = await response.json();
      setLogs(data || []);
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
      const response = await fetch("/api/mood", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ mood_rating: parseInt(moodRating), notes }),
      });

      if (response.ok) {
        setMoodRating(5);
        setNotes("");
        fetchData();
      } else {
        throw new Error("Failed to add mood log");
      }
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div className="bg-gray-900 text-white min-h-screen">
      <Header />
      <div className="p-4">
        <h1 className="text-2xl font-bold mb-4">Mood</h1>
        {loading && <p>Loading...</p>}
        {error && <p className="text-red-500">{error}</p>}
        {token ? (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div>
              <h2 className="text-xl font-bold mb-4">Log Your Mood</h2>
              <form
                onSubmit={handleSubmit}
                className="bg-gray-800 p-6 rounded-lg shadow-md"
              >
                <div className="mb-4">
                  <label
                    className="block text-gray-400 text-sm font-bold mb-2"
                    htmlFor="moodRating"
                  >
                    Mood Rating (1-10)
                  </label>
                  <input
                    id="moodRating"
                    type="range"
                    min="1"
                    max="10"
                    value={moodRating}
                    onChange={(e) => setMoodRating(e.target.value)}
                    className="w-full"
                  />
                  <div className="text-center">{moodRating}</div>
                </div>
                <div className="mb-6">
                  <label
                    className="block text-gray-400 text-sm font-bold mb-2"
                    htmlFor="notes"
                  >
                    Notes
                  </label>
                  <textarea
                    id="notes"
                    value={notes}
                    onChange={(e) => setNotes(e.target.value)}
                    className="shadow appearance-none border rounded w-full py-2 px-3 bg-gray-700 text-white mb-3 leading-tight focus:outline-none focus:shadow-outline"
                  />
                </div>
                <button
                  type="submit"
                  className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full"
                >
                  Add Log
                </button>
              </form>
            </div>

            <div>
              <h2 className="text-xl font-bold mb-4">Your Mood Logs</h2>
              <div className="bg-gray-800 p-6 rounded-lg shadow-md">
                {logs.length > 0 ? (
                  <ul>
                    {logs.map((log) => (
                      <li
                        key={log.ID}
                        className="border-b border-gray-700 py-4"
                      >
                        <p className="font-bold">Mood: {log.MoodRating}/10</p>
                        {log.Notes && <p>Notes: {log.Notes}</p>}
                        {log.ai_analysis && (
                          <p className="text-green-400 mt-2">
                            AI Analysis: {log.ai_analysis}
                          </p>
                        )}
                        <p className="text-sm text-gray-400 mt-2">
                          {new Date(log.CreatedAt).toLocaleString()}
                        </p>
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p>No mood logs yet.</p>
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
            to view your mood data.
          </p>
        )}
      </div>
    </div>
  );
};

export default Mood;
