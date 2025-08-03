import React, { useEffect } from 'react';
import Home from './pages/Home';
import './index.css';
import { SessionProvider } from './contexts/SessionContext';
import { apiService } from './services/api';
import { toast, ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

function App() {
  useEffect(() => {
    const initSession = async () => {
      try {
        await apiService.initSession();
        console.log("Session initialized");
      } catch (err) {
        console.error("Session init failed:", err);
        toast.error("Failed to initialize session. Please try again later.", {
          position: "top-right",
          autoClose: 5000,
        });
      }
    };

    initSession();
  }, []);

  return (
    <SessionProvider>
      <div className="App">
        <Home />
        <ToastContainer
          position="top-right"
          autoClose={5000}
          aria-label="Notification messages"
        />
      </div>
    </SessionProvider>
  );
}

export default App;
