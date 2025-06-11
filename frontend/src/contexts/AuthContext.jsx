import React, { createContext, useContext, useState, useEffect } from 'react';
import authService from '../services/authService';
import { getValidToken, cleanupInvalidToken, isAuthenticated as checkAuth } from '../utils/tokenUtils';

const AuthContext = createContext();

export { AuthContext };

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    // Check if user is already logged in on app start
    const initializeAuth = () => {
      try {
        const token = getValidToken();
        const currentUser = authService.getCurrentUser();
        
        if (currentUser && token) {
          setUser(currentUser);
          setIsAuthenticated(true);
        } else {
          // Clean up invalid tokens
          cleanupInvalidToken();
        }
      } catch (error) {
        console.error('Failed to initialize auth:', error);
        // Clear invalid tokens
        cleanupInvalidToken();
      } finally {
        setIsLoading(false);
      }
    };

    initializeAuth();
  }, []);

  const login = async (credentials) => {
    try {
      const response = await authService.login(credentials);
      console.log('Login response:', response);
      setUser(response.user);
      setIsAuthenticated(true);
      console.log('Auth state updated:', { user: response.user, isAuthenticated: true });
      return response;
    } catch (error) {
      throw error;
    }
  };

  const register = async (userData) => {
    try {
      const response = await authService.register(userData);
      return response;
    } catch (error) {
      throw error;
    }
  };

  const logout = () => {
    authService.logout();
    setUser(null);
    setIsAuthenticated(false);
  };

  const loginWithGoogle = async (credential) => {
    try {
      const response = await authService.loginWithGoogle(credential);
      console.log('Google login response:', response);
      setUser(response.user);
      setIsAuthenticated(true);
      console.log('Auth state updated with Google:', { user: response.user, isAuthenticated: true });
      return response;
    } catch (error) {
      throw error;
    }
  };

  const value = {
    user,
    isAuthenticated,
    isLoading,
    login,
    register,
    logout,
    loginWithGoogle
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}; 