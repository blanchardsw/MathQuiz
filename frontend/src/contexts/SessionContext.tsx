// SessionContext.tsx
import React, { createContext, useContext, useState, useEffect } from 'react';
import { apiService } from '../services/api';

interface SessionContextType {
  initialized: boolean;
}

const SessionContext = createContext<SessionContextType>({ initialized: false });

export const SessionProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [initialized, setInitialized] = useState(false);

  useEffect(() => {
    const init = async () => {
      try {
        await apiService.initSession();
        setInitialized(true);
      } catch (error) {
        console.error('Session init failed:', error);
      }
    };
    init();
  }, []);

  return (
    <SessionContext.Provider value={{ initialized }}>
      {children}
    </SessionContext.Provider>
  );
};

export const useSession = () => useContext(SessionContext);