import { useState, useEffect } from "react";
import Header from "../components/Header";
import useAuth from "../hooks/useAuth";

const Nutrition = () => {
  const { token } = useAuth();
  const [logs, setLogs] = useState([]);
  const [stats, setStats] = useState(null);
  const [suggestion, setSuggestion] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const [mealName, setMealName] = useState("");
  const [calories, setCalories] = useState("");

  const fetchData = async () => {
    try {
      const [logsRes, statsRes, suggestionRes] = await Promise.all([
        fetch("/api/nutrition", {
          headers: { Authorization: `Bearer ${token}` },
        }),
        fetch("/api/nutrition/stats", {
          headers: { Authorization: `Bearer ${token}` },
        }),
        fetch("/api/nutrition/suggestion", {
          headers: { Authorization: `Bearer ${token}` },
        }),
      ]);

      if (!logsRes.ok || !statsRes.ok || !suggestionRes.ok) {
        throw new Error("Failed to fetch nutrition data");
      }

      const logsData = await logsRes.json();
      const statsData = await statsRes.json();
      const suggestionData = await suggestionRes.json();

      setLogs(logsData || []);
      setStats(statsData);
      setSuggestion(suggestionData);
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
      const response = await fetch("/api/nutrition", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          meal_name: mealName,
          calories: parseInt(calories),
        }),
      });

      if (response.ok) {
        setMealName("");
        setCalories("");
        fetchData();
      } else {
        throw new Error("Failed to add nutrition log");
      }
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div className="bg-gray-900 text-white min-h-screen">
      <Header />
      <div className="p-4">
        <h1 className="text-2xl font-bold mb-4">Nutrition</h1>
        {loading && <p>Loading...</p>}
        {error && <p className="text-red-500">{error}</p>}
        {token ? (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div>
              <h2 className="text-xl font-bold mb-4">Log a Meal</h2>
              <form
                onSubmit={handleSubmit}
                className="bg-gray-800 p-6 rounded-lg shadow-md"
              >
                <div className="mb-4">
                  <label
                    className="block text-gray-400 text-sm font-bold mb-2"
                    htmlFor="mealName"
                  >
                    Meal Name
                  </label>
                  <input
                    id="mealName"
                    type="text"
                    value={mealName}
                    onChange={(e) => setMealName(e.target.value)}
                    className="shadow appearance-none border rounded w-full py-2 px-3 bg-gray-700 text-white leading-tight focus:outline-none focus:shadow-outline"
                    required
                  />
                </div>
                <div className="mb-6">
                  <label
                    className="block text-gray-400 text-sm font-bold mb-2"
                    htmlFor="calories"
                  >
                    Calories
                  </label>
                  <input
                    id="calories"
                    type="number"
                    value={calories}
                    onChange={(e) => setCalories(e.target.value)}
                    className="shadow appearance-none border rounded w-full py-2 px-3 bg-gray-700 text-white mb-3 leading-tight focus:outline-none focus:shadow-outline"
                    required
                  />
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
                    Average Daily Calories:{" "}
                    {stats.average_daily_calories.toFixed(2)}
                  </p>
                </div>
              )}

              {suggestion && (
                <div className="bg-gray-800 p-6 rounded-lg shadow-md mt-8">
                  <h2 className="text-xl font-bold mb-4">AI Suggestion</h2>
                  <p>{suggestion.suggestion}</p>
                </div>
              )}
            </div>

            <div>
              <h2 className="text-xl font-bold mb-4">Your Nutrition Logs</h2>
              <div className="bg-gray-800 p-6 rounded-lg shadow-md">
                {logs.length > 0 ? (
                  <ul>
                    {logs.map((log) => (
                      <li
                        key={log.ID}
                        className="border-b border-gray-700 py-2"
                      >
                        <p className="font-bold">{log.MealName}</p>
                        <p>{log.Calories} calories</p>
                        <p className="text-sm text-gray-400">
                          {new Date(log.CreatedAt).toLocaleString()}
                        </p>
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p>No nutrition logs yet.</p>
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
            to view your nutrition data.
          </p>
        )}
      </div>
    </div>
  );
};

export default Nutrition;
