import React from 'react';
import { useTranslation } from 'react-i18next';

const ProfilePage = () => {
  const { t } = useTranslation();

  return (
    <div className="max-w-4xl mx-auto">
      <div className="text-center py-16">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          {t('profile.title')}
        </h1>
        <p className="text-xl text-gray-600 mb-8">
          Manage your account settings and preferences
        </p>
        <div className="card p-8">
          <p className="text-gray-600">
            Profile management functionality will be implemented here
          </p>
        </div>
      </div>
    </div>
  );
};

export default ProfilePage; 