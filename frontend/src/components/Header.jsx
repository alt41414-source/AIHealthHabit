import { useState } from "react";
import { Link } from "react-router-dom";
import useAuth from "../hooks/useAuth";

const Header = () => {
  const { user, logout } = useAuth();
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const handleLogout = () => {
    logout();
    setDropdownOpen(false);
  };

  return (
    <header className="bg-gray-900 text-white p-8 flex justify-between items-center">
      <div className="text-xl font-bold">
        <Link to="/">Habit Tracker</Link>
      </div>
      <nav className="flex gap-4">
        <Link
          to="/nutrition"
          className="hover:text-gray-300 transition-transform hover:scale-110"
        >
          Nutrition
        </Link>
        <Link
          to="/mood"
          className="hover:text-gray-300 transition-transform hover:scale-110"
        >
          Mood
        </Link>
        <Link
          to="/sleep"
          className="hover:text-gray-300 transition-transform hover:scale-110"
        >
          Sleep
        </Link>
        <Link
          to="/leaderboard"
          className="hover:text-gray-300 transition-transform hover:scale-110"
        >
          Leaderboard
        </Link>
        <Link
          to="/habits"
          className="hover:text-gray-300 transition-transform hover:scale-110"
        >
          Habits
        </Link>
        {user ? (
          <div className="relative">
            <button
              onClick={() => setDropdownOpen(!dropdownOpen)}
              className="hover:text-gray-300 transition-transform hover:scale-110 focus:outline-none"
            >
              {user.username}
            </button>
            {dropdownOpen && (
              <div className="absolute right-0 mt-2 w-48 bg-gray-800 rounded-md shadow-lg py-1 z-20">
                <button
                  onClick={handleLogout}
                  className="block px-4 py-2 text-sm text-white hover:bg-gray-700 w-full text-left"
                >
                  Logout
                </button>
              </div>
            )}
          </div>
        ) : (
          <>
            <Link
              to="/login"
              className="hover:text-gray-300 transition-transform hover:scale-110"
            >
              Login
            </Link>
            <Link
              to="/register"
              className="hover:text-gray-300 transition-transform hover:scale-110"
            >
              Register
            </Link>
          </>
        )}
      </nav>
    </header>
  );
};

export default Header;
