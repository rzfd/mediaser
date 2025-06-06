import api from '../utils/api';

export const currencyService = {
  // Convert currency
  async convertCurrency(from, to, amount) {
    try {
      const response = await api.post('/api/currency/convert', {
        from,
        to,
        amount: parseFloat(amount)
      });
      return response.data;
    } catch (error) {
      console.error('Currency conversion error:', error);
      throw new Error(error.response?.data?.message || 'Failed to convert currency');
    }
  },

  // Get exchange rate
  async getExchangeRate(from, to) {
    try {
      const response = await api.get(`/api/currency/rate?from=${from}&to=${to}`);
      return response.data;
    } catch (error) {
      console.error('Exchange rate error:', error);
      throw new Error(error.response?.data?.message || 'Failed to get exchange rate');
    }
  },

  // Get supported currencies
  async getSupportedCurrencies() {
    try {
      const response = await api.get('/api/currency/list');
      return response.data;
    } catch (error) {
      console.error('Get currencies error:', error);
      throw new Error(error.response?.data?.message || 'Failed to get supported currencies');
    }
  },

  // Update currency rates
  async updateCurrencyRates() {
    try {
      const response = await api.post('/api/currency/update');
      return response.data;
    } catch (error) {
      console.error('Update rates error:', error);
      throw new Error(error.response?.data?.message || 'Failed to update currency rates');
    }
  },

  // Get currency preferences
  async getCurrencyPreferences() {
    try {
      const response = await api.get('/api/currency/preferences');
      return response.data;
    } catch (error) {
      console.error('Get currency preferences error:', error);
      throw new Error(error.response?.data?.message || 'Failed to get currency preferences');
    }
  },

  // Update currency preferences
  async updateCurrencyPreferences(preferences) {
    try {
      const response = await api.put('/api/currency/preferences', preferences);
      return response.data;
    } catch (error) {
      console.error('Update currency preferences error:', error);
      throw new Error(error.response?.data?.message || 'Failed to update currency preferences');
    }
  },

  // Format currency amount
  formatCurrency(amount, currency = 'USD', locale = 'en-US') {
    try {
      return new Intl.NumberFormat(locale, {
        style: 'currency',
        currency: currency.toUpperCase(),
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      }).format(amount);
    } catch (error) {
      console.error('Format currency error:', error);
      return `${currency.toUpperCase()} ${amount.toFixed(2)}`;
    }
  },

  // Get currency symbol
  getCurrencySymbol(currency) {
    const symbols = {
      'USD': '$',
      'EUR': '€',
      'GBP': '£',
      'JPY': '¥',
      'CNY': '¥',
      'IDR': 'Rp',
      'SGD': 'S$',
      'MYR': 'RM',
      'THB': '฿',
      'VND': '₫',
      'PHP': '₱',
      'KRW': '₩',
      'INR': '₹'
    };
    return symbols[currency.toUpperCase()] || currency.toUpperCase();
  }
};

export default currencyService; 