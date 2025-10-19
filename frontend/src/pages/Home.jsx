import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import useAuth from "../hooks/useAuth";
import Header from "../components/Header";

const Home = () => {
  const { user, token } = useAuth();
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (user) {
      const fetchData = async () => {
        try {
          const [habitsRes, nutritionRes, sleepRes, moodRes] =
            await Promise.all([
              fetch("/api/habits", {
                headers: { Authorization: `Bearer ${token}` },
              }),
              fetch("/api/nutrition", {
                headers: { Authorization: `Bearer ${token}` },
              }),
              fetch("/api/sleep", {
                headers: { Authorization: `Bearer ${token}` },
              }),
              fetch("/api/mood", {
                headers: { Authorization: `Bearer ${token}` },
              }),
            ]);

          if (
            !habitsRes.ok ||
            !nutritionRes.ok ||
            !sleepRes.ok ||
            !moodRes.ok
          ) {
            throw new Error("Failed to fetch data");
          }

          let habits = await habitsRes.json();
          if (!Array.isArray(habits)) habits = [];
          let nutrition = await nutritionRes.json();
          if (!Array.isArray(nutrition)) nutrition = [];
          let sleep = await sleepRes.json();
          if (!Array.isArray(sleep)) sleep = [];
          let mood = await moodRes.json();
          if (!Array.isArray(mood)) mood = [];

          setData({ habits, nutrition, sleep, mood });
        } catch (err) {
          setError(err.message);
        } finally {
          setLoading(false);
        }
      };

      fetchData();
    }
  }, [user, token]);

  return (
    <div className="bg-gray-900 text-white min-h-screen">
      <Header />
      <div className="p-4">
        {user ? (
          <div>
            <h1 className="text-2xl font-bold mb-4">Dashboard</h1>
            {loading && <p>Loading...</p>}
            {error && <p className="text-red-500">{error}</p>}
            {data && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                {/* Habits Card */}
                <div className="bg-gray-800 p-4 rounded-lg shadow">
                  <h2 className="text-xl font-bold mb-2">Habits</h2>
                  <p>{data.habits.length} habits tracked.</p>
                  <Link to="/habits" className="text-blue-500 hover:underline">
                    View More
                  </Link>
                </div>
                {/* Nutrition Card */}
                <div className="bg-gray-800 p-4 rounded-lg shadow">
                  <h2 className="text-xl font-bold mb-2">Nutrition</h2>
                  <p>{data.nutrition.length} nutrition logs.</p>
                  <Link
                    to="/nutrition"
                    className="text-blue-500 hover:underline"
                  >
                    View More
                  </Link>
                </div>
                {/* Sleep Card */}
                <div className="bg-gray-800 p-4 rounded-lg shadow">
                  <h2 className="text-xl font-bold mb-2">Sleep</h2>
                  <p>{data.sleep.length} sleep logs.</p>
                  <Link to="/sleep" className="text-blue-500 hover:underline">
                    View More
                  </Link>
                </div>
                {/* Mood Card */}
                <div className="bg-gray-800 p-4 rounded-lg shadow">
                  <h2 className="text-xl font-bold mb-2">Mood</h2>
                  <p>{data.mood.length} mood logs.</p>
                  <Link to="/mood" className="text-blue-500 hover:underline">
                    View More
                  </Link>
                </div>
              </div>
            )}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center text-center pt-20">
            <h1 className="text-4xl font-bold mb-4">
              Welcome to Habit Tracker
            </h1>
            <p className="text-lg mb-8">
              Start your journey to a healthier and more productive life by
              tracking your habits.
            </p>
            <div className="flex gap-4">
              <Link
                to="/login"
                className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition-transform hover:scale-110"
              >
                Login
              </Link>
              <Link
                to="/register"
                className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded transition-transform hover:scale-110"
              >
                Register
              </Link>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Home;
