import React from 'react';
import { useTranslation } from 'react-i18next';

const LoginPage = () => {
  const { t } = useTranslation();

  return (
    <div className="max-w-md mx-auto">
      <div className="text-center py-16">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          {t('auth.login')}
        </h1>
        <p className="text-xl text-gray-600 mb-8">
          Welcome back to MediaShar
        </p>
        <div className="card p-8">
          <p className="text-gray-600">
            Login form will be implemented here
          </p>
        </div>
      </div>
    </div>
  );
};

export default LoginPage; 