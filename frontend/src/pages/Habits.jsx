import { useState, useEffect } from "react";
import Header from "../components/Header";
import useAuth from "../hooks/useAuth";

const Habits = () => {
  const { token } = useAuth();
  const [habits, setHabits] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [newHabitName, setNewHabitName] = useState("");

  const fetchData = async () => {
    try {
      const response = await fetch("/api/habits", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!response.ok) {
        throw new Error("Failed to fetch habits");
      }
      let data = await response.json();
      if (!Array.isArray(data)) {
        data = [];
      }

      const habitsWithStats = await Promise.all(
        data.map(async (habit) => {
          const statsRes = await fetch(`/api/habits/${habit.ID}/stats`, {
            headers: { Authorization: `Bearer ${token}` },
          });
          if (!statsRes.ok) {
            throw new Error(`Failed to fetch stats for habit ${habit.ID}`);
          }
          const stats = await statsRes.json();
          return { ...habit, ...stats };
        }),
      );

      setHabits(habitsWithStats);
    } catch (err) {
      setError(err.message);
      console.error("Error fetching habits:", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (token) {
      fetchData();
    }
  }, [token]);

  const handleCreateHabit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch("/api/habits", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ name: newHabitName }),
      });

      if (response.ok) {
        setNewHabitName("");
        fetchData();
      } else {
        throw new Error("Failed to create habit");
      }
    } catch (err) {
      setError(err.message);
      console.error("Error creating habit:", err);
    }
  };

  const handleTrackHabit = async (habitId) => {
    try {
      const response = await fetch(`/api/habits/${habitId}/track`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      });

      if (response.ok) {
        fetchData();
      } else {
        throw new Error("Failed to track habit");
      }
    } catch (err) {
      setError(err.message);
      console.error("Error tracking habit:", err);
    }
  };

  return (
    <div className="bg-gray-900 text-white min-h-screen">
      <Header />
      <div className="p-4">
        <h1 className="text-2xl font-bold mb-4">Habits</h1>
        {loading && <p>Loading...</p>}
        {error && <p className="text-red-500">{error}</p>}
        {token ? (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div>
              <h2 className="text-xl font-bold mb-4">Create a New Habit</h2>
              <form
                onSubmit={handleCreateHabit}
                className="bg-gray-800 p-6 rounded-lg shadow-md"
              >
                <div className="mb-4">
                  <label
                    className="block text-gray-400 text-sm font-bold mb-2"
                    htmlFor="habitName"
                  >
                    Habit Name
                  </label>
                  <input
                    id="habitName"
                    type="text"
                    value={newHabitName}
                    onChange={(e) => setNewHabitName(e.target.value)}
                    className="shadow appearance-none border rounded w-full py-2 px-3 bg-gray-700 text-white leading-tight focus:outline-none focus:shadow-outline"
                    required
                  />
                </div>
                <button
                  type="submit"
                  className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full"
                >
                  Create Habit
                </button>
              </form>
            </div>

            <div>
              <h2 className="text-xl font-bold mb-4">Your Habits</h2>
              <div className="bg-gray-800 p-6 rounded-lg shadow-md">
                {habits && habits.length > 0 ? (
                  <ul>
                    {habits.map((habit) => (
                      <li
                        key={habit.ID}
                        className="border-b border-gray-700 py-4 flex justify-between items-center"
                      >
                        <div>
                          <p className="font-bold">{habit.Name}</p>
                          <p>Streak: {habit.streak}</p>
                          <p>
                            Completion Rate (30d):{" "}
                            {habit.completion_rate.toFixed(2)}%
                          </p>
                        </div>
                        <button
                          onClick={() => handleTrackHabit(habit.ID)}
                          className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
                        >
                          Track Today
                        </button>
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p>No habits created yet.</p>
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
            to view your habits.
          </p>
        )}
      </div>
    </div>
  );
};

export default Habits;
