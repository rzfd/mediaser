import React, { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Menu, X, Globe, DollarSign, User, LogOut } from 'lucide-react';

// UI Components
import LanguageSelector from '../ui/LanguageSelector';
import CurrencySelector from '../ui/CurrencySelector';

const Navbar = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);
  const { t } = useTranslation();
  const location = useLocation();

  // Mock user data - replace with actual auth context
  const user = JSON.parse(localStorage.getItem('mediashar-user') || 'null');
  const isAuthenticated = !!user;

  const navItems = [
    { path: '/', label: t('navigation.home') },
    { path: '/donate', label: t('navigation.donate') },
    { path: '/streamers', label: t('navigation.streamers') },
  ];

  const isActiveRoute = (path) => {
    return location.pathname === path;
  };

  const handleLogout = () => {
    localStorage.removeItem('mediashar-token');
    localStorage.removeItem('mediashar-user');
    window.location.href = '/';
  };

  return (
    <nav className="bg-white shadow-lg border-b border-gray-200 sticky top-0 z-50">
      <div className="container mx-auto px-4">
        <div className="flex justify-between items-center h-16">
          {/* Logo and Brand */}
          <Link to="/" className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-gradient-to-br from-primary-500 to-primary-700 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-sm">MS</span>
            </div>
            <span className="text-xl font-bold gradient-text">
              {t('app.title')}
            </span>
          </Link>

          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center space-x-8">
            {/* Main Navigation */}
            <div className="flex items-center space-x-6">
              {navItems.map((item) => (
                <Link
                  key={item.path}
                  to={item.path}
                  className={`font-medium transition-colors duration-200 ${
                    isActiveRoute(item.path)
                      ? 'text-primary-600 border-b-2 border-primary-600 pb-1'
                      : 'text-gray-700 hover:text-primary-600'
                  }`}
                >
                  {item.label}
                </Link>
              ))}
            </div>

            {/* Utility Navigation */}
            <div className="flex items-center space-x-4">
              {/* Language Selector */}
              <div className="flex items-center space-x-2">
                <Globe className="w-4 h-4 text-gray-500" />
                <LanguageSelector />
              </div>

              {/* Currency Selector */}
              <div className="flex items-center space-x-2">
                <DollarSign className="w-4 h-4 text-gray-500" />
                <CurrencySelector />
              </div>

              {/* Authentication */}
              {isAuthenticated ? (
                <div className="relative">
                  <button
                    onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
                    className="flex items-center space-x-2 p-2 rounded-lg hover:bg-gray-100 transition-colors duration-200"
                  >
                    <div className="w-8 h-8 bg-primary-100 rounded-full flex items-center justify-center">
                      <User className="w-4 h-4 text-primary-600" />
                    </div>
                    <span className="text-sm font-medium text-gray-700">
                      {user?.username || user?.email}
                    </span>
                  </button>

                  {/* User Dropdown Menu */}
                  {isUserMenuOpen && (
                    <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 py-2">
                      <Link
                        to="/profile"
                        className="flex items-center px-4 py-2 text-gray-700 hover:bg-gray-50 transition-colors duration-200"
                        onClick={() => setIsUserMenuOpen(false)}
                      >
                        <User className="w-4 h-4 mr-2" />
                        {t('navigation.profile')}
                      </Link>
                      <hr className="my-2" />
                      <button
                        onClick={handleLogout}
                        className="w-full flex items-center px-4 py-2 text-red-600 hover:bg-red-50 transition-colors duration-200"
                      >
                        <LogOut className="w-4 h-4 mr-2" />
                        {t('navigation.logout')}
                      </button>
                    </div>
                  )}
                </div>
              ) : (
                <div className="flex items-center space-x-4">
                  <Link
                    to="/login"
                    className="font-medium text-gray-700 hover:text-primary-600 transition-colors duration-200"
                  >
                    {t('navigation.login')}
                  </Link>
                  <Link
                    to="/register"
                    className="btn-primary"
                  >
                    {t('navigation.register')}
                  </Link>
                </div>
              )}
            </div>
          </div>

          {/* Mobile Menu Button */}
          <div className="md:hidden">
            <button
              onClick={() => setIsMenuOpen(!isMenuOpen)}
              className="p-2 rounded-lg hover:bg-gray-100 transition-colors duration-200"
            >
              {isMenuOpen ? (
                <X className="w-6 h-6 text-gray-700" />
              ) : (
                <Menu className="w-6 h-6 text-gray-700" />
              )}
            </button>
          </div>
        </div>

        {/* Mobile Navigation Menu */}
        {isMenuOpen && (
          <div className="md:hidden py-4 border-t border-gray-200">
            <div className="space-y-4">
              {/* Main Navigation */}
              <div className="space-y-2">
                {navItems.map((item) => (
                  <Link
                    key={item.path}
                    to={item.path}
                    onClick={() => setIsMenuOpen(false)}
                    className={`block px-4 py-2 rounded-lg font-medium transition-colors duration-200 ${
                      isActiveRoute(item.path)
                        ? 'bg-primary-50 text-primary-600'
                        : 'text-gray-700 hover:bg-gray-50'
                    }`}
                  >
                    {item.label}
                  </Link>
                ))}
              </div>

              {/* Utility Tools */}
              <div className="space-y-3 px-4 py-3 bg-gray-50 rounded-lg">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <Globe className="w-4 h-4 text-gray-500" />
                    <span className="text-sm font-medium text-gray-700">
                      {t('language.title')}
                    </span>
                  </div>
                  <LanguageSelector />
                </div>

                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <DollarSign className="w-4 h-4 text-gray-500" />
                    <span className="text-sm font-medium text-gray-700">
                      {t('currency.title')}
                    </span>
                  </div>
                  <CurrencySelector />
                </div>
              </div>

              {/* Authentication */}
              <div className="space-y-2">
                {isAuthenticated ? (
                  <>
                    <Link
                      to="/profile"
                      onClick={() => setIsMenuOpen(false)}
                      className="flex items-center px-4 py-2 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors duration-200"
                    >
                      <User className="w-4 h-4 mr-2" />
                      {t('navigation.profile')}
                    </Link>
                    <button
                      onClick={handleLogout}
                      className="w-full flex items-center px-4 py-2 rounded-lg text-red-600 hover:bg-red-50 transition-colors duration-200"
                    >
                      <LogOut className="w-4 h-4 mr-2" />
                      {t('navigation.logout')}
                    </button>
                  </>
                ) : (
                  <>
                    <Link
                      to="/login"
                      onClick={() => setIsMenuOpen(false)}
                      className="block px-4 py-2 rounded-lg font-medium text-gray-700 hover:bg-gray-50 transition-colors duration-200"
                    >
                      {t('navigation.login')}
                    </Link>
                    <Link
                      to="/register"
                      onClick={() => setIsMenuOpen(false)}
                      className="block px-4 py-2 rounded-lg font-medium bg-primary-600 text-white text-center hover:bg-primary-700 transition-colors duration-200"
                    >
                      {t('navigation.register')}
                    </Link>
                  </>
                )}
              </div>
            </div>
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navbar; 