import { useState, useEffect } from "react";
import { jwtDecode } from "jwt-decode";
import AuthContext from "./AuthContext";

const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(localStorage.getItem("token"));
  const [username, setUsername] = useState(localStorage.getItem("username"));

  useEffect(() => {
    if (token && username) {
      try {
        const decodedToken = jwtDecode(token);
        if (decodedToken.exp * 1000 < Date.now()) {
          localStorage.removeItem("token");
          localStorage.removeItem("username");
          setToken(null);
          setUsername(null);
          setUser(null);
        } else {
          setUser({ id: decodedToken.sub, token, username });
        }
      } catch (error) {
        console.error("Invalid token", error);
        localStorage.removeItem("token");
        localStorage.removeItem("username");
        setToken(null);
        setUsername(null);
        setUser(null);
      }
    } else {
      setUser(null);
    }
  }, [token, username]);

  const login = async (username, password) => {
    const response = await fetch("/api/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });

    if (response.ok) {
      const { token, username: returnedUsername } = await response.json();
      localStorage.setItem("token", token);
      localStorage.setItem("username", returnedUsername);
      setToken(token);
      setUsername(returnedUsername);
      return true;
    } else {
      return false;
    }
  };

  const register = async (username, password) => {
    const response = await fetch("/api/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });

    return response.ok;
  };

  const logout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("username");
    setToken(null);
    setUsername(null);
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, register, logout, token }}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthProvider;
