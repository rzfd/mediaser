import React from 'react';
import { useTranslation } from 'react-i18next';

const DonatePage = () => {
  const { t } = useTranslation();

  return (
    <div className="max-w-4xl mx-auto">
      <div className="text-center py-16">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          {t('donation.title')}
        </h1>
        <p className="text-xl text-gray-600 mb-8">
          Support your favorite creators with multi-currency donations
        </p>
        <div className="card p-8">
          <p className="text-gray-600">
            Donation functionality will be implemented here with Midtrans integration
          </p>
        </div>
      </div>
    </div>
  );
};

export default DonatePage; 