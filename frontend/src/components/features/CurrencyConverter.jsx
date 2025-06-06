import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { ArrowRightLeft, RefreshCw, TrendingUp, Clock } from 'lucide-react';
import toast from 'react-hot-toast';
import { currencyService } from '../../services/currencyService';

const CurrencyConverter = () => {
  const { t } = useTranslation();
  const [fromCurrency, setFromCurrency] = useState('USD');
  const [toCurrency, setToCurrency] = useState('IDR');
  const [amount, setAmount] = useState('100');
  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(false);
  const [currencies, setCurrencies] = useState([]);
  const [exchangeRate, setExchangeRate] = useState(null);
  const [lastUpdated, setLastUpdated] = useState(null);

  // Load supported currencies
  useEffect(() => {
    const loadCurrencies = async () => {
      try {
        const response = await currencyService.getSupportedCurrencies();
        setCurrencies(response.currencies || []);
      } catch (error) {
        console.error('Failed to load currencies:', error);
        // Use default currencies if API fails
        setCurrencies([
          { code: 'USD', name: 'US Dollar' },
          { code: 'IDR', name: 'Indonesian Rupiah' },
          { code: 'EUR', name: 'Euro' },
          { code: 'JPY', name: 'Japanese Yen' },
          { code: 'SGD', name: 'Singapore Dollar' },
          { code: 'MYR', name: 'Malaysian Ringgit' },
          { code: 'CNY', name: 'Chinese Yuan' }
        ]);
      }
    };

    loadCurrencies();
  }, []);

  // Load exchange rate when currencies change
  useEffect(() => {
    if (fromCurrency && toCurrency && fromCurrency !== toCurrency) {
      loadExchangeRate();
    }
  }, [fromCurrency, toCurrency]);

  const loadExchangeRate = async () => {
    try {
      const response = await currencyService.getExchangeRate(fromCurrency, toCurrency);
      setExchangeRate(response.rate);
      setLastUpdated(new Date(response.updated_at));
    } catch (error) {
      console.error('Failed to load exchange rate:', error);
    }
  };

  const handleConvert = async () => {
    if (!amount || isNaN(amount) || parseFloat(amount) <= 0) {
      toast.error(t('currency.error'));
      return;
    }

    setLoading(true);
    try {
      const response = await currencyService.convertCurrency(
        fromCurrency,
        toCurrency,
        parseFloat(amount)
      );
      
      setResult(response);
      setExchangeRate(response.rate);
      setLastUpdated(new Date(response.updated_at));
      
      toast.success(t('currency.result'));
    } catch (error) {
      console.error('Conversion failed:', error);
      toast.error(error.message || t('currency.error'));
    } finally {
      setLoading(false);
    }
  };

  const handleSwapCurrencies = () => {
    setFromCurrency(toCurrency);
    setToCurrency(fromCurrency);
    setResult(null);
  };

  const formatCurrency = (value, currency) => {
    return currencyService.formatCurrency(value, currency);
  };

  const formatNumber = (value) => {
    return new Intl.NumberFormat('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 6
    }).format(value);
  };

  return (
    <div className="max-w-2xl mx-auto">
      <div className="card p-6">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            {t('currency.converter')}
          </h1>
          <p className="text-gray-600">
            {t('app.description')}
          </p>
        </div>

        {/* Converter Form */}
        <div className="space-y-6">
          {/* Amount Input */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              {t('currency.amount')}
            </label>
            <input
              type="number"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              placeholder="0.00"
              className="input text-lg font-mono"
              min="0"
              step="0.01"
            />
          </div>

          {/* Currency Selection */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-end">
            {/* From Currency */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                {t('currency.from')}
              </label>
              <select
                value={fromCurrency}
                onChange={(e) => setFromCurrency(e.target.value)}
                className="input"
              >
                {currencies.map((currency) => (
                  <option key={currency.code} value={currency.code}>
                    {currency.code} - {currency.name}
                  </option>
                ))}
              </select>
            </div>

            {/* Swap Button */}
            <div className="flex justify-center">
              <button
                onClick={handleSwapCurrencies}
                className="p-3 rounded-full bg-gray-100 hover:bg-gray-200 transition-colors duration-200"
                aria-label="Swap currencies"
              >
                <ArrowRightLeft className="w-5 h-5 text-gray-600" />
              </button>
            </div>

            {/* To Currency */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                {t('currency.to')}
              </label>
              <select
                value={toCurrency}
                onChange={(e) => setToCurrency(e.target.value)}
                className="input"
              >
                {currencies.map((currency) => (
                  <option key={currency.code} value={currency.code}>
                    {currency.code} - {currency.name}
                  </option>
                ))}
              </select>
            </div>
          </div>

          {/* Convert Button */}
          <button
            onClick={handleConvert}
            disabled={loading || !amount || fromCurrency === toCurrency}
            className="w-full btn-primary flex items-center justify-center space-x-2 py-3"
          >
            {loading ? (
              <>
                <RefreshCw className="w-5 h-5 animate-spin" />
                <span>{t('currency.loading')}</span>
              </>
            ) : (
              <>
                <TrendingUp className="w-5 h-5" />
                <span>{t('currency.convert')}</span>
              </>
            )}
          </button>
        </div>

        {/* Results */}
        {result && (
          <div className="mt-8 p-6 bg-gradient-to-r from-primary-50 to-blue-50 rounded-lg border border-primary-200">
            <div className="text-center">
              <div className="text-sm text-gray-600 mb-2">
                {t('currency.result')}
              </div>
              <div className="text-3xl font-bold text-gray-900 mb-4">
                {formatCurrency(result.converted_amount, toCurrency)}
              </div>
              
              {/* Exchange Rate Info */}
              <div className="space-y-2 text-sm text-gray-600">
                <div className="flex items-center justify-center space-x-2">
                  <TrendingUp className="w-4 h-4" />
                  <span>
                    1 {fromCurrency} = {formatNumber(result.rate)} {toCurrency}
                  </span>
                </div>
                
                {lastUpdated && (
                  <div className="flex items-center justify-center space-x-2">
                    <Clock className="w-4 h-4" />
                    <span>
                      {t('currency.updated')}: {lastUpdated.toLocaleString()}
                    </span>
                  </div>
                )}
              </div>
            </div>
          </div>
        )}

        {/* Exchange Rate Preview */}
        {exchangeRate && !result && fromCurrency !== toCurrency && (
          <div className="mt-6 p-4 bg-gray-50 rounded-lg">
            <div className="text-center text-sm text-gray-600">
              <div className="flex items-center justify-center space-x-2">
                <TrendingUp className="w-4 h-4" />
                <span>
                  1 {fromCurrency} = {formatNumber(exchangeRate)} {toCurrency}
                </span>
              </div>
              
              {lastUpdated && (
                <div className="flex items-center justify-center space-x-2 mt-1">
                  <Clock className="w-4 h-4" />
                  <span>
                    {t('currency.updated')}: {lastUpdated.toLocaleString()}
                  </span>
                </div>
              )}
            </div>
          </div>
        )}

        {/* Quick Amount Buttons */}
        <div className="mt-6">
          <div className="text-sm font-medium text-gray-700 mb-3">
            Quick amounts:
          </div>
          <div className="grid grid-cols-4 gap-2">
            {['1', '10', '100', '1000'].map((quickAmount) => (
              <button
                key={quickAmount}
                onClick={() => setAmount(quickAmount)}
                className="px-3 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors duration-200"
              >
                {quickAmount}
              </button>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default CurrencyConverter; 