import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Youtube, Music, Link, Upload, Check, X, Info } from 'lucide-react';
import { toast } from 'react-hot-toast';

const DonationMediaShare = ({ 
  donationAmount, 
  streamerSettings, 
  onMediaSubmit,
  className = ""
}) => {
  const { t } = useTranslation();
  const [mediaType, setMediaType] = useState('youtube');
  const [mediaUrl, setMediaUrl] = useState('');
  const [mediaTitle, setMediaTitle] = useState('');
  const [mediaMessage, setMediaMessage] = useState('');
  const [isEligible, setIsEligible] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  // Mock streamer settings if not provided
  const defaultSettings = {
    mediaShareEnabled: true,
    minDonationAmount: 5000,
    currency: 'IDR',
    allowYoutube: true,
    allowTiktok: true,
    welcomeMessage: 'Terima kasih atas donasi Anda! Silakan bagikan media favorit Anda.'
  };

  const settings = streamerSettings || defaultSettings;

  useEffect(() => {
    // Check if donation amount is eligible for media share
    const eligible = donationAmount >= settings.minDonationAmount && settings.mediaShareEnabled;
    setIsEligible(eligible);
  }, [donationAmount, settings]);

  // URL validation functions
  const isValidYouTubeUrl = (url) => {
    return /^(https?:\/\/)?(www\.)?(youtube\.com\/watch\?v=|youtu\.be\/)[a-zA-Z0-9_-]+/.test(url);
  };

  const isValidTikTokUrl = (url) => {
    return /^(https?:\/\/)?(www\.)?tiktok\.com\/@[^\/]+\/video\/\d+/.test(url);
  };

  const validateMediaUrl = () => {
    if (!mediaUrl.trim()) return false;
    
    if (mediaType === 'youtube') {
      return isValidYouTubeUrl(mediaUrl);
    } else if (mediaType === 'tiktok') {
      return isValidTikTokUrl(mediaUrl);
    }
    
    return false;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateMediaUrl()) {
      toast.error(t('mediaShare.invalidUrl'));
      return;
    }

    setSubmitting(true);
    
    try {
      const mediaData = {
        type: mediaType,
        url: mediaUrl,
        title: mediaTitle || `${mediaType === 'youtube' ? 'YouTube' : 'TikTok'} Video`,
        message: mediaMessage,
        donationAmount: donationAmount,
        timestamp: new Date().toISOString()
      };

      // Call parent component's callback
      if (onMediaSubmit) {
        await onMediaSubmit(mediaData);
      }

      toast.success(t('mediaShare.submitSuccess'));
      
      // Reset form
      setMediaUrl('');
      setMediaTitle('');
      setMediaMessage('');
      
    } catch (error) {
      console.error('Error submitting media:', error);
      toast.error(t('mediaShare.submitError'));
    } finally {
      setSubmitting(false);
    }
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: settings.currency,
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount);
  };

  if (!settings.mediaShareEnabled) {
    return null;
  }

  return (
    <div className={`card p-6 ${className}`}>
      <div className="flex items-center mb-4">
        <Upload className="w-6 h-6 text-primary-600 mr-2" />
        <h3 className="text-lg font-semibold text-gray-900">
          {t('mediaShare.shareMedia')}
        </h3>
      </div>

      {/* Eligibility Status */}
      <div className={`p-3 rounded-lg mb-4 ${isEligible ? 'bg-green-50 border border-green-200' : 'bg-yellow-50 border border-yellow-200'}`}>
        <div className="flex items-start">
          {isEligible ? (
            <Check className="w-5 h-5 text-green-600 mr-2 mt-0.5" />
          ) : (
            <Info className="w-5 h-5 text-yellow-600 mr-2 mt-0.5" />
          )}
          <div className="flex-1">
            {isEligible ? (
              <div>
                <p className="text-sm font-medium text-green-800">
                  {t('mediaShare.eligible')}
                </p>
                <p className="text-xs text-green-600 mt-1">
                  {settings.welcomeMessage}
                </p>
              </div>
            ) : (
              <div>
                <p className="text-sm font-medium text-yellow-800">
                  {t('mediaShare.notEligible')}
                </p>
                <p className="text-xs text-yellow-600 mt-1">
                  {t('mediaShare.minimumRequired')}: {formatCurrency(settings.minDonationAmount)}
                  <br />
                  {t('mediaShare.currentDonation')}: {formatCurrency(donationAmount)}
                </p>
              </div>
            )}
          </div>
        </div>
      </div>

      {isEligible && (
        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Platform Selection */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              {t('mediaShare.selectPlatform')}
            </label>
            <div className="flex space-x-4">
              {settings.allowYoutube && (
                <button
                  type="button"
                  onClick={() => setMediaType('youtube')}
                  className={`flex items-center px-4 py-2 rounded-lg border transition-colors ${
                    mediaType === 'youtube'
                      ? 'border-red-500 bg-red-50 text-red-700'
                      : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50'
                  }`}
                >
                  <Youtube className="w-4 h-4 mr-2" />
                  YouTube
                </button>
              )}
              {settings.allowTiktok && (
                <button
                  type="button"
                  onClick={() => setMediaType('tiktok')}
                  className={`flex items-center px-4 py-2 rounded-lg border transition-colors ${
                    mediaType === 'tiktok'
                      ? 'border-pink-500 bg-pink-50 text-pink-700'
                      : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50'
                  }`}
                >
                  <Music className="w-4 h-4 mr-2" />
                  TikTok
                </button>
              )}
            </div>
          </div>

          {/* Media URL */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {mediaType === 'youtube' ? t('mediaShare.youtubeUrl') : t('mediaShare.tiktokUrl')} *
            </label>
            <div className="relative">
              <Link className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
              <input
                type="url"
                value={mediaUrl}
                onChange={(e) => setMediaUrl(e.target.value)}
                placeholder={
                  mediaType === 'youtube' 
                    ? 'https://www.youtube.com/watch?v=...' 
                    : 'https://www.tiktok.com/@username/video/...'
                }
                className="pl-10 w-full p-3 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
                required
              />
            </div>
            {mediaUrl && !validateMediaUrl() && (
              <p className="text-xs text-red-600 mt-1 flex items-center">
                <X className="w-3 h-3 mr-1" />
                {t('mediaShare.invalidUrlFormat')}
              </p>
            )}
          </div>

          {/* Media Title */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {t('mediaShare.mediaTitle')}
            </label>
            <input
              type="text"
              value={mediaTitle}
              onChange={(e) => setMediaTitle(e.target.value)}
              placeholder={t('mediaShare.mediaTitlePlaceholder')}
              className="w-full p-3 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
            />
          </div>

          {/* Message to Streamer */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {t('mediaShare.messageToStreamer')}
            </label>
            <textarea
              value={mediaMessage}
              onChange={(e) => setMediaMessage(e.target.value)}
              placeholder={t('mediaShare.messageToStreamerPlaceholder')}
              rows={3}
              className="w-full p-3 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
            />
          </div>

          {/* Submit Button */}
          <button
            type="submit"
            disabled={submitting || !validateMediaUrl()}
            className="w-full btn-primary flex items-center justify-center space-x-2"
          >
            <Upload className="w-4 h-4" />
            <span>
              {submitting ? t('mediaShare.submitting') : t('mediaShare.shareWithStreamer')}
            </span>
          </button>
        </form>
      )}
    </div>
  );
};

export default DonationMediaShare; 