import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import AuthProvider from "./contexts/AuthProvider";
import Home from "./pages/Home";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Nutrition from "./pages/Nutrition";
import Mood from "./pages/Mood";
import Sleep from "./pages/Sleep";
import Leaderboard from "./pages/Leaderboard";
import Habits from "./pages/Habits";

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/nutrition" element={<Nutrition />} />
          <Route path="/mood" element={<Mood />} />
          <Route path="/sleep" element={<Sleep />} />
          <Route path="/leaderboard" element={<Leaderboard />} />
          <Route path="/habits" element={<Habits />} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
