import React from 'react';
import Home from './pages/Home';
import './index.css';
import { useEffect } from "react";
import axios from 'axios';
import { SessionProvider } from './contexts/SessionContext';

function App() {
  useEffect(() => {
    const initSession = async () => {
      const baseURL = process.env.REACT_APP_API_URL || 'http://localhost:4000/api';

      try {
        const res = await axios.post(baseURL + "/init-session", {
          withCredentials: true,
        });

        if (res.status === 200) {
          console.log("Session initialized");
        } else {
          console.warn("Unexpected status:", res.status);
        }
      } catch (err) {
        console.error("Session init failed:", err);
      }
    };

    initSession();
  }, []);

  return (
    <SessionProvider>
      <div className="App">
        <Home />
      </div>
    </SessionProvider>
  );
}

export default App;
