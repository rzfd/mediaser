import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { ChevronDown, Check } from 'lucide-react';
import { currencyService } from '../../services/currencyService';

const CurrencySelector = () => {
  const { t } = useTranslation();
  const [isOpen, setIsOpen] = useState(false);
  const [currencies, setCurrencies] = useState([]);
  const [selectedCurrency, setSelectedCurrency] = useState(
    localStorage.getItem('mediashar-currency') || 'USD'
  );

  const defaultCurrencies = [
    { code: 'USD', name: 'US Dollar', symbol: '$' },
    { code: 'IDR', name: 'Indonesian Rupiah', symbol: 'Rp' },
    { code: 'EUR', name: 'Euro', symbol: '€' },
    { code: 'JPY', name: 'Japanese Yen', symbol: '¥' },
    { code: 'SGD', name: 'Singapore Dollar', symbol: 'S$' },
    { code: 'MYR', name: 'Malaysian Ringgit', symbol: 'RM' },
    { code: 'CNY', name: 'Chinese Yuan', symbol: '¥' }
  ];

  useEffect(() => {
    const loadCurrencies = async () => {
      try {
        const response = await currencyService.getSupportedCurrencies();
        const supportedCurrencies = response.currencies?.map(curr => ({
          ...curr,
          symbol: currencyService.getCurrencySymbol(curr.code)
        })) || defaultCurrencies;
        setCurrencies(supportedCurrencies);
      } catch (error) {
        console.warn('Failed to load currencies from API, using defaults:', error);
        setCurrencies(defaultCurrencies);
      }
    };

    loadCurrencies();
  }, []);

  const currentCurrency = currencies.find(curr => curr.code === selectedCurrency) || 
                         defaultCurrencies.find(curr => curr.code === selectedCurrency) ||
                         defaultCurrencies[0];

  const handleCurrencyChange = async (currencyCode) => {
    try {
      setSelectedCurrency(currencyCode);
      setIsOpen(false);
      
      // Save preference to localStorage
      localStorage.setItem('mediashar-currency', currencyCode);
      
      // Dispatch custom event for other components to listen
      window.dispatchEvent(new CustomEvent('currencyChanged', {
        detail: { currency: currencyCode }
      }));
      
    } catch (error) {
      console.error('Failed to change currency:', error);
    }
  };

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (isOpen && !event.target.closest('.currency-selector')) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [isOpen]);

  return (
    <div className="relative currency-selector">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center space-x-2 px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 transition-colors duration-200"
        aria-label={t('currency.title')}
        aria-expanded={isOpen}
        aria-haspopup="listbox"
      >
        <span className="currency-symbol text-green-600 font-semibold">
          {currentCurrency?.symbol || '$'}
        </span>
        <span className="hidden sm:inline-block font-mono">
          {currentCurrency?.code || 'USD'}
        </span>
        <ChevronDown 
          className={`w-4 h-4 transition-transform duration-200 ${
            isOpen ? 'rotate-180' : ''
          }`} 
        />
      </button>

      {isOpen && (
        <div className="absolute right-0 mt-2 w-64 bg-white border border-gray-200 rounded-lg shadow-lg z-50 max-h-60 overflow-y-auto">
          <div className="py-2">
            <div className="px-3 py-2 text-xs font-medium text-gray-500 uppercase tracking-wider border-b border-gray-100">
              {t('currency.title')}
            </div>
            
            {currencies.map((currency) => (
              <button
                key={currency.code}
                onClick={() => handleCurrencyChange(currency.code)}
                className={`w-full flex items-center justify-between px-3 py-2 text-sm hover:bg-gray-50 transition-colors duration-200 ${
                  selectedCurrency === currency.code
                    ? 'bg-primary-50 text-primary-700'
                    : 'text-gray-900'
                }`}
                role="option"
                aria-selected={selectedCurrency === currency.code}
              >
                <div className="flex items-center space-x-3">
                  <span className="currency-symbol text-green-600 font-semibold w-6 text-center">
                    {currency.symbol}
                  </span>
                  <div className="text-left">
                    <div className="font-medium font-mono">{currency.code}</div>
                    <div className="text-xs text-gray-500">{currency.name}</div>
                  </div>
                </div>
                
                {selectedCurrency === currency.code && (
                  <Check className="w-4 h-4 text-primary-600" />
                )}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default CurrencySelector; 