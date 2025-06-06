import React from 'react';
import { useTranslation } from 'react-i18next';

const StreamersPage = () => {
  const { t } = useTranslation();

  return (
    <div className="max-w-6xl mx-auto">
      <div className="text-center py-16">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          {t('streamers.title')}
        </h1>
        <p className="text-xl text-gray-600 mb-8">
          Discover and support amazing creators from around the world
        </p>
        <div className="card p-8">
          <p className="text-gray-600">
            Streamers listing and search functionality will be implemented here
          </p>
        </div>
      </div>
    </div>
  );
};

export default StreamersPage; 