import React, { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../../contexts/AuthContext';

const GoogleLoginButton = ({ onSuccess, onError }) => {
  const { t } = useTranslation();
  const { loginWithGoogle } = useAuth();

  useEffect(() => {
    // Check if Google Client ID is configured
    if (!process.env.REACT_APP_GOOGLE_CLIENT_ID) {
      console.warn('Google Client ID not configured');
      return;
    }

    // Check if script already exists
    if (document.querySelector('script[src="https://accounts.google.com/gsi/client"]')) {
      initializeGoogleSignIn();
      return;
    }

    // Load Google Identity Services
    const script = document.createElement('script');
    script.src = 'https://accounts.google.com/gsi/client';
    script.async = true;
    script.defer = true;
    script.onload = initializeGoogleSignIn;
    script.onerror = () => {
      console.error('Failed to load Google Identity Services');
    };
    document.head.appendChild(script);

    return () => {
      // Only remove if we added it
      const existingScript = document.querySelector('script[src="https://accounts.google.com/gsi/client"]');
      if (existingScript && existingScript === script) {
        document.head.removeChild(script);
      }
    };
  }, []);

  const initializeGoogleSignIn = () => {
    if (window.google && process.env.REACT_APP_GOOGLE_CLIENT_ID) {
      try {
        window.google.accounts.id.initialize({
          client_id: process.env.REACT_APP_GOOGLE_CLIENT_ID,
          callback: handleCredentialResponse,
          auto_select: false,
          cancel_on_tap_outside: true,
        });

        const buttonContainer = document.getElementById('google-signin-button');
        if (buttonContainer) {
          window.google.accounts.id.renderButton(
            buttonContainer,
            {
              theme: 'outline',
              size: 'large',
              text: 'signin_with',
              shape: 'rectangular',
            }
          );
        } else {
          // Show fallback button if container not found
          const fallbackButton = document.getElementById('manual-google-button');
          if (fallbackButton) {
            fallbackButton.style.display = 'block';
          }
        }
      } catch (error) {
        console.error('Failed to initialize Google Sign-In:', error);
        // Show fallback button on error
        const fallbackButton = document.getElementById('manual-google-button');
        if (fallbackButton) {
          fallbackButton.style.display = 'block';
        }
      }
    }
  };

  const handleCredentialResponse = async (response) => {
    try {
      await loginWithGoogle(response.credential);
      if (onSuccess) {
        onSuccess();
      }
    } catch (error) {
      console.error('Google login error:', error);
      if (onError) {
        onError(error);
      }
    }
  };

  const handleManualClick = () => {
    if (window.google) {
      window.google.accounts.id.prompt();
    }
  };

  // Don't render if Google Client ID is not configured
  if (!process.env.REACT_APP_GOOGLE_CLIENT_ID) {
    return null;
  }

  return (
    <div className="w-full">
      <div id="google-signin-button" className="w-full"></div>
      {/* Fallback button if Google button doesn't render */}
      <button
        type="button"
        onClick={handleManualClick}
        className="w-full mt-2 flex items-center justify-center px-4 py-2 border border-gray-300 rounded-lg shadow-sm bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        style={{ display: 'none' }}
        id="manual-google-button"
      >
        <svg className="w-5 h-5 mr-2" viewBox="0 0 24 24">
          <path
            fill="#4285F4"
            d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
          />
          <path
            fill="#34A853"
            d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
          />
          <path
            fill="#FBBC05"
            d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
          />
          <path
            fill="#EA4335"
            d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
          />
        </svg>
        {t('auth.loginWithGoogle')}
      </button>
    </div>
  );
};

export default GoogleLoginButton; 