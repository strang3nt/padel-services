import { initData, useSignal } from '@tma.js/sdk-react';
import { createContext, FC, ReactNode, useContext, useEffect, useState } from 'react';

interface ContextContent {
  bearerToken: string | null
  loading: boolean,
  error: string | null
}

interface ResponseContent {
  token: string
  id: string
}

const AuthContext = createContext<ContextContent>({
  bearerToken: null,
  loading: true,
  error: null
});

export const AuthProvider: FC<{ children: ReactNode }> = ({ children }) => {
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [bearerToken, setBearerToken] = useState<string | null>(null);
  const initDataRaw = useSignal(initData.raw);

  useEffect(() => {
    const initAuth = async () => {
      try {
        // Get data from Telegram WebApp global

        if (!initData) {
          throw new Error("Not a Telegram client");
        }

        const response: Response = await fetch('/auth', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ initDataRaw }),
        })

        if (response.status === 403) {
          setError("Access denied: user not allowed.")
        }

        if (!response.ok) {
          setError("Authentication failed")
        } else {
          const data = await response.json().then(x => x as ResponseContent)
          setBearerToken(data.token)
        }
      } catch (err) {
        setError(`Access denied: ${err}.`)

      } finally {
        setLoading(false);
      }
    };

    initAuth();
  }, []);

  return (
    <AuthContext.Provider value={{ bearerToken, loading, error }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
