import React from 'react';
import { Routes, Route } from 'react-router-dom';
import { Toaster } from 'react-hot-toast';
import { useTranslation } from 'react-i18next';

// Context Providers
import { AuthProvider, useAuth } from './contexts/AuthContext';

// Layout Components
import Navbar from './components/layout/Navbar';
import Footer from './components/layout/Footer';

// UI Components
import LoadingSpinner from './components/ui/LoadingSpinner';

// Page Components
import HomePage from './pages/HomePage';
import DonatePage from './pages/DonatePage';
import StreamersPage from './pages/StreamersPage';
import ProfilePage from './pages/ProfilePage';
import LoginPage from './pages/auth/LoginPage';
import RegisterPage from './pages/auth/RegisterPage';

// Feature Components
import CurrencyConverter from './components/features/CurrencyConverter';
import LanguageTranslator from './components/features/LanguageTranslator';

// App Content Component (needs to be inside AuthProvider)
function AppContent() {
  const { i18n } = useTranslation();
  const { isLoading } = useAuth();

  // Set document language for accessibility
  React.useEffect(() => {
    document.documentElement.lang = i18n.language;
  }, [i18n.language]);

  // Show loading spinner while auth is initializing
  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <LoadingSpinner size="xl" />
          <p className="mt-4 text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="App min-h-screen bg-gray-50 flex flex-col">
      {/* Toast notifications */}
      <Toaster
        position="top-right"
        toastOptions={{
          duration: 4000,
          style: {
            background: '#363636',
            color: '#fff',
          },
          success: {
            duration: 3000,
            iconTheme: {
              primary: '#22c55e',
              secondary: '#fff',
            },
          },
          error: {
            duration: 5000,
            iconTheme: {
              primary: '#ef4444',
              secondary: '#fff',
            },
          },
        }}
      />

      {/* Navigation */}
      <Navbar />

      {/* Main Content */}
      <main className="flex-1 container mx-auto px-4 py-8">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/donate" element={<DonatePage />} />
          <Route path="/donate/:streamerId" element={<DonatePage />} />
          <Route path="/streamers" element={<StreamersPage />} />
          <Route path="/profile" element={<ProfilePage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/currency" element={<CurrencyConverter />} />
          <Route path="/language" element={<LanguageTranslator />} />
          
          {/* 404 Route */}
          <Route path="*" element={
            <div className="text-center py-16">
              <h1 className="text-4xl font-bold text-gray-900 mb-4">404</h1>
              <p className="text-gray-600 mb-8">Page not found</p>
              <a 
                href="/" 
                className="btn-primary"
              >
                Go Home
              </a>
            </div>
          } />
        </Routes>
      </main>

      {/* Footer */}
      <Footer />
    </div>
  );
}

// Main App Component
function App() {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  );
}

export default App; 