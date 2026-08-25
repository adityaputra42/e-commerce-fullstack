"use client";

import React, { createContext, useContext, useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { authService, setAccessToken, User } from '../services/api';

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (token: string, user: User) => void;
  logout: () => void;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
     const bootstrapAuth = async () => {
      const refreshResult = await authService.refresh();

      if (!refreshResult?.access_token) {
        setLoading(false);
        return;
      }

      setAccessToken(refreshResult.access_token);

      try {
        const userData = await authService.getCurrentUser();
        setUser(userData);
      } catch (e) {
        console.error("Failed to fetch user profile after refresh", e);
        setAccessToken(null);
        setUser(null);
      }

      setLoading(false);
    };

    bootstrapAuth();
  }, []);

  const login = (token: string, userData: User) => {
    setAccessToken(token);
    setUser(userData);
    router.push('/');
  };

  const logout = () => {
    authService.logout().finally(() => {
      setUser(null);
      router.push('/login');
    });
  };

  return (
    <AuthContext.Provider value={{
      user,
      loading,
      login,
      logout,
      isAuthenticated: !!user
    }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
