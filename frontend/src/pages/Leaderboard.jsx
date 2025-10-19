import { useState, useEffect } from "react";
import Header from "../components/Header";
import useAuth from "../hooks/useAuth";

const Leaderboard = () => {
  const { token } = useAuth();
  const [leaderboards, setLeaderboards] = useState({
    nutrition: [],
    sleep: [],
    habit: [],
  });
  const [activeTab, setActiveTab] = useState("nutrition");
  const [currentPage, setCurrentPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const itemsPerPage = 10;

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [nutritionRes, sleepRes, habitRes] = await Promise.all([
          fetch("/api/leaderboard/nutrition", {
            headers: { Authorization: `Bearer ${token}` },
          }),
          fetch("/api/leaderboard/sleep", {
            headers: { Authorization: `Bearer ${token}` },
          }),
          fetch("/api/leaderboard/habit", {
            headers: { Authorization: `Bearer ${token}` },
          }),
        ]);

        if (!nutritionRes.ok || !sleepRes.ok || !habitRes.ok) {
          throw new Error("Failed to fetch leaderboards");
        }

        const nutritionData = await nutritionRes.json();
        const sleepData = await sleepRes.json();
        const habitData = await habitRes.json();

        setLeaderboards({
          nutrition: nutritionData || [],
          sleep: sleepData || [],
          habit: habitData || [],
        });
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    if (token) {
      fetchData();
    }
  }, [token]);

  const handleTabClick = (tab) => {
    setActiveTab(tab);
    setCurrentPage(1);
  };

  const activeLeaderboard = leaderboards[activeTab];
  const indexOfLastItem = currentPage * itemsPerPage;
  const indexOfFirstItem = indexOfLastItem - itemsPerPage;
  const currentItems = activeLeaderboard.slice(
    indexOfFirstItem,
    indexOfLastItem
  );
  const totalPages = Math.max(1, Math.ceil(activeLeaderboard.length / itemsPerPage));

  const renderLeaderboard = () => (
    <div className="bg-gray-800 p-6 rounded-lg shadow-md">
      <table className="w-full text-left">
        <thead>
          <tr>
            <th className="p-2">Rank</th>
            <th className="p-2">Username</th>
            <th className="p-2">Score</th>
          </tr>
        </thead>
        <tbody>
          {currentItems.map((entry, index) => (
            <tr key={index} className="border-b border-gray-700">
              <td className="p-2">{indexOfFirstItem + index + 1}</td>
              <td className="p-2">{entry.username}</td>
              <td className="p-2">{entry.score.toFixed(2)}</td>
            </tr>
          ))}
        </tbody>
      </table>
      <div className="flex justify-between mt-4">
        <button
          onClick={() => setCurrentPage(currentPage - 1)}
          disabled={currentPage === 1}
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded disabled:bg-gray-600"
        >
          Previous
        </button>
        <span>
          Page {currentPage} of {totalPages}
        </span>
        <button
          onClick={() => setCurrentPage(currentPage + 1)}
          disabled={currentPage === totalPages}
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded disabled:bg-gray-600"
        >
          Next
        </button>
      </div>
    </div>
  );

  return (
    <div className="bg-gray-900 text-white min-h-screen">
      <Header />
      <div className="p-4">
        <h1 className="text-2xl font-bold mb-4">Leaderboards</h1>
        {loading && <p>Loading...</p>}
        {error && <p className="text-red-500">{error}</p>}
        {token ? (
          <div>
            <div className="flex border-b border-gray-700 mb-4">
              <button
                onClick={() => handleTabClick("nutrition")}
                className={`py-2 px-4 ${
                  activeTab === "nutrition" ? "border-b-2 border-blue-500" : ""
                }`}
              >
                Nutrition
              </button>
              <button
                onClick={() => handleTabClick("sleep")}
                className={`py-2 px-4 ${
                  activeTab === "sleep" ? "border-b-2 border-blue-500" : ""
                }`}
              >
                Sleep
              </button>
              <button
                onClick={() => handleTabClick("habit")}
                className={`py-2 px-4 ${
                  activeTab === "habit" ? "border-b-2 border-blue-500" : ""
                }`}
              >
                Habits
              </button>
            </div>
            {renderLeaderboard()}
          </div>
        ) : (
          <p>
            Please{" "}
            <a href="/login" className="text-blue-500 hover:underline">
              login
            </a>{" "}
            to view the leaderboards.
          </p>
        )}
      </div>
    </div>
  );
};

export default Leaderboard;
