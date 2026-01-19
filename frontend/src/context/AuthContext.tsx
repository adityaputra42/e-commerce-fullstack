"use client";

import React, { createContext, useContext, useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { authService, User, LoginResponse } from '../services/api';

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
    // Check for existing token on mount
    const checkAuth = async () => {
      const token = localStorage.getItem('token');
      const storedUser = localStorage.getItem('user');

      // Optimistic load: if we have a token and stored user, set it immediately
      if (token && storedUser && storedUser !== "undefined") {
        try {
           setUser(JSON.parse(storedUser));
        } catch (e) {
           console.error("Error parsing stored user", e);
        }
      }
      
      if (token) {
        try {
          // Verify token and get fresh user data
          const userData = await authService.getCurrentUser();
          setUser(userData);
          // Update stored user data
          localStorage.setItem('user', JSON.stringify(userData));
        } catch (e: any) {
          console.error("Failed to fetch user profile", e);
          
          // Only log out on authentication errors (401/403), not network errors
          if (e.response && (e.response.status === 401 || e.response.status === 403)) {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            setUser(null);
          }
        }
      }
      setLoading(false);
    };

    checkAuth();
  }, []);

  const login = (token: string, userData: User) => {
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(userData));
    setUser(userData);
    router.push('/');
  };

  const logout = () => {
    authService.logout();
    setUser(null);
    router.push('/login');
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
