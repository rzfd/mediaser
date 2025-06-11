import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import GoogleLoginButton from '../../components/auth/GoogleLoginButton';

const LoginPage = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { login } = useAuth();
  const [formData, setFormData] = useState({
    email: '',
    password: ''
  });
  const [errors, setErrors] = useState({});
  const [isLoading, setIsLoading] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    // Clear error when user starts typing
    if (errors[name]) {
      setErrors(prev => ({
        ...prev,
        [name]: ''
      }));
    }
  };

  const validateForm = () => {
    const newErrors = {};

    if (!formData.email.trim()) {
      newErrors.email = t('auth.errors.emailRequired');
    } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
      newErrors.email = t('auth.errors.emailInvalid');
    }

    if (!formData.password) {
      newErrors.password = t('auth.errors.passwordRequired');
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setIsLoading(true);
    try {
      await login(formData);
      
      // Show success message
      alert(t('auth.loginSuccess'));
      
      // Redirect to home page
      navigate('/');
    } catch (error) {
      setErrors({
        general: error.message || t('auth.errors.loginFailed')
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleGoogleSuccess = () => {
    // Show success message
    alert(t('auth.loginSuccess'));
    
    // Redirect to home page
    navigate('/');
  };

  const handleGoogleError = (error) => {
    setErrors({
      general: error.message || t('auth.errors.googleLoginFailed')
    });
  };

  return (
    <div className="max-w-md mx-auto">
      <div className="text-center py-8">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          {t('auth.login')}
        </h1>
        <p className="text-xl text-gray-600 mb-8">
          Welcome back to MediaShar
        </p>
        
        <div className="card p-8">
          {/* Google Login Button */}
          <div className="mb-6">
            <GoogleLoginButton 
              onSuccess={handleGoogleSuccess}
              onError={handleGoogleError}
            />
          </div>

          {/* Divider */}
          <div className="relative mb-6">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-gray-300"></div>
            </div>
            <div className="relative flex justify-center text-sm">
              <span className="px-2 bg-white text-gray-500">{t('auth.orLoginWith')}</span>
            </div>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            {errors.general && (
              <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg">
                {errors.general}
              </div>
            )}

            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
                {t('auth.email')}
              </label>
              <input
                type="email"
                id="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                className={`w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                  errors.email ? 'border-red-300' : 'border-gray-300'
                }`}
                placeholder={t('auth.emailPlaceholder')}
                autoComplete="email"
              />
              {errors.email && (
                <p className="text-red-500 text-sm mt-1">{errors.email}</p>
              )}
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-1">
                {t('auth.password')}
              </label>
              <input
                type="password"
                id="password"
                name="password"
                value={formData.password}
                onChange={handleChange}
                className={`w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                  errors.password ? 'border-red-300' : 'border-gray-300'
                }`}
                placeholder={t('auth.passwordPlaceholder')}
                autoComplete="current-password"
              />
              {errors.password && (
                <p className="text-red-500 text-sm mt-1">{errors.password}</p>
              )}
            </div>

            <button
              type="submit"
              disabled={isLoading}
              className={`w-full py-2 px-4 rounded-lg font-medium text-white transition-colors ${
                isLoading
                  ? 'bg-gray-400 cursor-not-allowed'
                  : 'bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500'
              }`}
            >
              {isLoading ? t('auth.loggingIn') : t('auth.login')}
            </button>
          </form>

          <div className="mt-6 text-center">
          <p className="text-gray-600">
              {t('auth.dontHaveAccount')}{' '}
              <Link to="/register" className="text-blue-600 hover:text-blue-800 font-medium">
                {t('auth.register')}
              </Link>
            </p>
          </div>

          {/* Quick Login for Testing */}
          <div className="mt-6 pt-6 border-t border-gray-200">
            <p className="text-sm text-gray-500 mb-3 text-center">
              {t('auth.quickLogin')}
            </p>
            <div className="flex space-x-2">
              <button
                type="button"
                onClick={() => setFormData({ email: 'streamer@test.com', password: 'password123' })}
                className="flex-1 py-2 px-3 text-xs bg-purple-100 text-purple-800 rounded-lg hover:bg-purple-200 transition-colors"
              >
                {t('auth.testStreamer')}
              </button>
              <button
                type="button"
                onClick={() => setFormData({ email: 'donator@test.com', password: 'password123' })}
                className="flex-1 py-2 px-3 text-xs bg-green-100 text-green-800 rounded-lg hover:bg-green-200 transition-colors"
              >
                {t('auth.testDonator')}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LoginPage; 