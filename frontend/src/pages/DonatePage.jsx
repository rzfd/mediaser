import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams, useNavigate } from 'react-router-dom';
import { Heart, User, DollarSign, CreditCard, MessageCircle, AlertCircle } from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import DonationMediaShare from '../components/media/DonationMediaShare';
import { apiRequest, getValidToken } from '../utils/tokenUtils';
import { processSnapPayment, fallbackToRedirect, waitForSnap } from '../utils/midtransUtils';

const DonatePage = () => {
  const { t } = useTranslation();
  const { streamerId } = useParams();
  const navigate = useNavigate();
  const { user, isAuthenticated, isLoading: authLoading } = useAuth();

  // State management
  const [selectedStreamer, setSelectedStreamer] = useState(null);
  const [streamers, setStreamers] = useState([]);
  const [formData, setFormData] = useState({
    amount: 50000,
    currency: 'IDR',
    message: '',
    display_name: '',
    is_anonymous: false,
    streamer_id: streamerId || ''
  });
  const [loading, setLoading] = useState(false);
  const [streamersLoading, setStreamersLoading] = useState(true);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(false);
  const [sharedMedia, setSharedMedia] = useState(null);

  // API base URL
  const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

  // Predefined amounts
  const predefinedAmounts = [25000, 50000, 100000, 200000, 500000];

  // Load streamers and selected streamer
  useEffect(() => {
    loadStreamers();
  }, []);

  useEffect(() => {
    if (streamerId && streamers.length > 0) {
      const streamer = streamers.find(s => s.id === parseInt(streamerId));
      if (streamer) {
        setSelectedStreamer(streamer);
        setFormData(prev => ({ ...prev, streamer_id: streamerId }));
      }
    }
  }, [streamerId, streamers]);

  // Set default display name when user is authenticated
  useEffect(() => {
    if (isAuthenticated && user) {
      setFormData(prev => ({
        ...prev,
        display_name: prev.display_name || user.username || user.email
      }));
    }
  }, [isAuthenticated, user]);

  const loadStreamers = async () => {
    try {
      setStreamersLoading(true);
      
      const response = await apiRequest(`${API_BASE_URL}/streamers`, {
        method: 'GET'
      });
      
      if (!response.ok) {
        throw new Error('Failed to fetch streamers');
      }
      
      const data = await response.json();
      const streamersData = data.data || data || [];
      setStreamers(streamersData);
      
      // Set mock data if no streamers
      if (streamersData.length === 0) {
        setStreamers([
          {
            id: 1,
            username: 'gaming_master',
            full_name: 'Gaming Master',
            description: 'Professional gamer and content creator'
          },
          {
            id: 2,
            username: 'music_lover',
            full_name: 'Music Lover',
            description: 'Singer and music producer'
          }
        ]);
      }
    } catch (err) {
      console.error('Error loading streamers:', err);
      // Set mock data on error
      setStreamers([
        {
          id: 1,
          username: 'gaming_master',
          full_name: 'Gaming Master',
          description: 'Professional gamer and content creator'
        },
        {
          id: 2,
          username: 'music_lover',
          full_name: 'Music Lover',
          description: 'Singer and music producer'
        }
      ]);
    } finally {
      setStreamersLoading(false);
    }
  };

  const handleInputChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
    setError(null);
  };

  const handleAmountSelect = (amount) => {
    setFormData(prev => ({ ...prev, amount }));
    setError(null);
  };

  const handleStreamerSelect = (streamerId) => {
    const streamer = streamers.find(s => s.id === parseInt(streamerId));
    setSelectedStreamer(streamer);
    setFormData(prev => ({ ...prev, streamer_id: streamerId }));
    setError(null);
  };

  const handleMediaSubmit = async (mediaData) => {
    console.log('Media shared:', mediaData);
    setSharedMedia(mediaData);
    // Here you would typically send the media data to your backend
    // For now, we'll just store it in state
  };

  const validateForm = () => {
    if (!formData.streamer_id) {
      setError(t('donation.selectStreamer'));
      return false;
    }
    
    if (!formData.amount || formData.amount < 10000) {
      setError(t('donation.minimumAmount') + ' 10,000 IDR');
      return false;
    }
    
    if (formData.amount > 10000000) {
      setError(t('donation.maximumAmount') + ' 10,000,000 IDR');
      return false;
    }
    
    if (!formData.display_name.trim()) {
      setError(t('donation.enterDisplayName'));
      return false;
    }
    
    return true;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (authLoading) {
      return;
    }
    
    if (!isAuthenticated) {
      alert('Please login first to make a donation');
      navigate('/login');
      return;
    }
    
    if (!validateForm()) {
      return;
    }
    
    setLoading(true);
    setError(null);
    
    try {
      // Check authentication using context first
      if (!isAuthenticated) {
        alert('Please login first to make a donation');
        navigate('/login');
        return;
      }
      
      // Create donation - token will be handled by apiRequest
      
      const donationResponse = await apiRequest(`${API_BASE_URL}/donations`, {
        method: 'POST',
        body: JSON.stringify({
          streamer_id: parseInt(formData.streamer_id),
          amount: formData.amount,
          currency: formData.currency,
          message: formData.message,
          display_name: formData.display_name,
          is_anonymous: formData.is_anonymous
        })
      });
      
      if (!donationResponse.ok) {
        throw new Error('Failed to create donation');
      }
      
      const donationData = await donationResponse.json();
      const donationId = donationData.data?.id || donationData.id;
      
      console.log('Donation created:', donationData);
      console.log('Donation ID:', donationId);
      
      // Create payment using Midtrans
      const paymentResponse = await apiRequest(`${API_BASE_URL}/midtrans/payment/${donationId}`, {
        method: 'POST'
      });
      
      console.log('Payment response status:', paymentResponse.status);
      
      const paymentData = await paymentResponse.json();
      console.log('Payment response data:', paymentData);
      
      if (!paymentResponse.ok) {
        const errorMsg = paymentData.message || paymentData.error || 'Failed to create payment';
        console.error('Payment creation failed:', errorMsg);
        throw new Error(errorMsg);
      }
      // Parse snap token - Midtrans returns token in data.token field
      const snapToken = paymentData.data?.token || paymentData.token || paymentData.data?.snap_token || paymentData.snap_token;
      const redirectUrl = paymentData.data?.redirect_url || paymentData.redirect_url;
      
      console.log('Snap token:', snapToken);
      console.log('Redirect URL:', redirectUrl);
      console.log('Window.snap available:', !!window.snap);
      console.log('Full payment data:', paymentData);
      
      // Try to process payment with Snap
      if (snapToken) {
        try {
          const snapAvailable = await waitForSnap(3000); // Wait up to 3 seconds for Snap
          console.log('Snap available after check:', snapAvailable);
          
          if (snapAvailable) {
            // Use Snap payment
            await processSnapPayment(snapToken, {
              onSuccess: (result) => {
                setSuccess(true);
                setFormData({
                  amount: 50000,
                  currency: 'IDR',
                  message: '',
                  display_name: user?.username || '',
                  is_anonymous: false,
                  streamer_id: ''
                });
                setSelectedStreamer(null);
                alert(t('donation.success'));
              },
              onPending: (result) => {
                alert('Payment is pending. Please complete your payment.');
              },
              onError: (result) => {
                setError('Payment failed. Please try again.');
              },
              onClose: () => {
                console.log('Payment popup closed by user');
              }
            });
          } else {
            // Fallback to redirect if Snap not available
            if (fallbackToRedirect(redirectUrl)) {
              setSuccess(true);
              alert('Opening payment page in new tab...');
            } else {
              // Last resort - show token info
              alert('Payment created successfully! Token: ' + snapToken.substring(0, 20) + '...');
              setSuccess(true);
            }
          }
        } catch (error) {
          console.error('Payment processing error:', error);
          // Try redirect fallback
          if (fallbackToRedirect(redirectUrl)) {
            setSuccess(true);
            alert('Opening payment page in new tab...');
          } else {
            setError('Failed to open payment interface: ' + error.message);
          }
        }
      } else {
        // No snap token, but donation was created successfully
        alert('Donation created successfully! Payment will be processed separately.');
        setSuccess(true);
      }
      
    } catch (err) {
      console.error('Error creating donation:', err);
      setError(err.message || t('donation.failed'));
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0
    }).format(amount);
  };

  if (streamersLoading || authLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <LoadingSpinner size="xl" />
          <p className="mt-4 text-gray-600">
            {authLoading ? 'Loading authentication...' : 'Loading donation page...'}
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      {/* Header */}
      <div className="text-center py-8">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          {t('donation.title')}
        </h1>
        <p className="text-xl text-gray-600 mb-8">
          Support your favorite creators with multi-currency donations
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Left Column: Donation Form */}
        <div className="bg-white rounded-lg shadow-lg p-6">
          <h2 className="text-2xl font-bold text-gray-900 mb-6 flex items-center">
            <Heart className="w-6 h-6 text-red-500 mr-2" />
            Make a Donation
          </h2>

          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Streamer Selection */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                {t('donation.selectStreamer')}
              </label>
              <select
                name="streamer_id"
                value={formData.streamer_id}
                onChange={(e) => handleStreamerSelect(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                required
              >
                <option value="">{t('donation.selectStreamer')}</option>
                {streamers.map(streamer => (
                  <option key={streamer.id} value={streamer.id}>
                    {streamer.full_name || streamer.username} (@{streamer.username})
                  </option>
                ))}
              </select>
            </div>

            {/* Amount Selection */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                {t('donation.amount')}
              </label>
              <div className="grid grid-cols-3 gap-2 mb-4">
                {predefinedAmounts.map(amount => (
                  <button
                    key={amount}
                    type="button"
                    onClick={() => handleAmountSelect(amount)}
                    className={`py-2 px-3 rounded-lg text-sm font-medium transition-colors ${
                      formData.amount === amount
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                    }`}
                  >
                    {formatCurrency(amount)}
                  </button>
                ))}
              </div>
              <input
                type="number"
                name="amount"
                value={formData.amount}
                onChange={handleInputChange}
                min="10000"
                max="10000000"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder={t('donation.enterAmount')}
                required
              />
            </div>

            {/* Currency Selection */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                {t('currency.title')}
              </label>
              <select
                name="currency"
                value={formData.currency}
                onChange={handleInputChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="IDR">IDR - Indonesian Rupiah</option>
                <option value="USD">USD - US Dollar</option>
              </select>
            </div>

            {/* Display Name */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                {t('donation.displayName')}
              </label>
              <input
                type="text"
                name="display_name"
                value={formData.display_name}
                onChange={handleInputChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder={t('donation.enterDisplayName')}
                required
              />
            </div>

            {/* Message */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                {t('donation.message')} (Optional)
              </label>
              <textarea
                name="message"
                value={formData.message}
                onChange={handleInputChange}
                rows="3"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder={t('donation.enterMessage')}
              />
            </div>

            {/* Media Share Option - Always show for testing */}
            <div className="border border-blue-200 rounded-lg p-4 bg-blue-50">
              <h4 className="text-sm font-medium text-blue-800 mb-2">🎬 Media Share Feature</h4>
              <p className="text-xs text-blue-600 mb-3">Share your favorite YouTube or TikTok videos with streamers!</p>
              <DonationMediaShare
                donationAmount={formData.amount}
                streamerSettings={{
                  mediaShareEnabled: true,
                  minDonationAmount: 5000,
                  currency: 'IDR',
                  allowYoutube: true,
                  allowTiktok: true,
                  welcomeMessage: 'Terima kasih atas donasi Anda! Silakan bagikan media favorit Anda.'
                }}
                onMediaSubmit={handleMediaSubmit}
                className="mb-4"
              />
            </div>

            {/* Anonymous Option */}
            <div className="flex items-center">
              <input
                type="checkbox"
                id="is_anonymous"
                name="is_anonymous"
                checked={formData.is_anonymous}
                onChange={handleInputChange}
                className="mr-2 h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
              />
              <label htmlFor="is_anonymous" className="text-sm text-gray-700">
                {t('donation.anonymous')}
              </label>
            </div>

            {/* Error Message */}
            {error && (
              <div className="p-4 bg-red-50 border border-red-200 rounded-lg flex items-center">
                <AlertCircle className="w-5 h-5 text-red-500 mr-2" />
                <p className="text-red-700 text-sm">{error}</p>
              </div>
            )}

            {/* Submit Button */}
            <button
              type="submit"
              disabled={loading || !isAuthenticated}
              className={`w-full py-3 px-4 rounded-lg font-medium transition-colors flex items-center justify-center ${
                loading || !isAuthenticated
                  ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                  : 'bg-gradient-to-r from-red-500 to-pink-600 text-white hover:from-red-600 hover:to-pink-700'
              }`}
            >
              {loading ? (
                <>
                  <LoadingSpinner size="sm" className="mr-2" />
                  {t('donation.processing')}
                </>
              ) : (
                <>
                  <CreditCard className="w-5 h-5 mr-2" />
                  {t('donation.submit')}
                </>
              )}
            </button>

            {!isAuthenticated && (
              <p className="text-center text-sm text-gray-500">
                Please <a href="/login" className="text-blue-600 hover:underline">login</a> to make a donation
              </p>
            )}
          </form>
        </div>

        {/* Right Column: Selected Streamer Info */}
        <div className="bg-white rounded-lg shadow-lg p-6">
          <h3 className="text-xl font-bold text-gray-900 mb-6">
            {selectedStreamer ? 'Supporting' : 'Select a Streamer'}
          </h3>

          {selectedStreamer ? (
            <div className="space-y-4">
              {/* Streamer Avatar */}
              <div className="flex items-center">
                <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center">
                  {selectedStreamer.avatar_url ? (
                    <img
                      src={selectedStreamer.avatar_url}
                      alt={selectedStreamer.username}
                      className="w-16 h-16 rounded-full object-cover"
                    />
                  ) : (
                    <User className="w-8 h-8 text-white" />
                  )}
                </div>
                <div className="ml-4">
                  <h4 className="text-lg font-semibold text-gray-900">
                    {selectedStreamer.full_name || selectedStreamer.username}
                  </h4>
                  <p className="text-gray-500">@{selectedStreamer.username}</p>
                </div>
              </div>

              {/* Description */}
              {selectedStreamer.description && (
                <div>
                  <h5 className="text-sm font-medium text-gray-700 mb-2">About</h5>
                  <p className="text-gray-600 text-sm">{selectedStreamer.description}</p>
                </div>
              )}

              {/* Donation Summary */}
              <div className="border-t pt-4">
                <h5 className="text-sm font-medium text-gray-700 mb-2">Donation Summary</h5>
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-gray-600">Amount:</span>
                    <span className="font-medium">{formatCurrency(formData.amount)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Display Name:</span>
                    <span className="font-medium">
                      {formData.is_anonymous ? 'Anonymous' : formData.display_name || 'Not set'}
                    </span>
                  </div>
                  {formData.message && (
                    <div>
                      <span className="text-gray-600">Message:</span>
                      <p className="text-sm italic text-gray-600 mt-1">"{formData.message}"</p>
                    </div>
                  )}
                </div>
              </div>
            </div>
          ) : (
            <div className="text-center py-8">
              <Heart className="w-16 h-16 text-gray-300 mx-auto mb-4" />
              <p className="text-gray-500">Select a streamer to see their information</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default DonatePage; 