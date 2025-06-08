import React, { useState, useContext } from 'react';
import { useTranslation } from 'react-i18next';
import { AuthContext } from '../../contexts/AuthContext';
import { Settings, DollarSign, Save, Youtube, Music, ToggleLeft, ToggleRight } from 'lucide-react';
import { toast } from 'react-hot-toast';

const StreamerMediaSettings = () => {
  const { t } = useTranslation();
  const { user } = useContext(AuthContext);
  const [settings, setSettings] = useState({
    mediaShareEnabled: true,
    minDonationAmount: 5000, // IDR
    currency: 'IDR',
    allowYoutube: true,
    allowTiktok: true,
    autoApprove: false,
    maxDurationYoutube: 300, // seconds
    maxDurationTiktok: 180, // seconds
    welcomeMessage: 'Terima kasih atas donasi Anda! Silakan bagikan media favorit Anda.'
  });
  const [saving, setSaving] = useState(false);

  // Only show for streamers
  if (!user || user.userType !== 'streamer') {
    return (
      <div className="card p-8 text-center">
        <Settings className="w-16 h-16 text-gray-400 mx-auto mb-4" />
        <h2 className="text-xl font-semibold text-gray-900 mb-2">
          {t('mediaShare.streamerOnly')}
        </h2>
        <p className="text-gray-600">
          {t('mediaShare.streamerOnlyDesc')}
        </p>
      </div>
    );
  }

  const handleSave = async () => {
    setSaving(true);
    try {
      // Here you would make an API call to save settings
      // For now, we'll simulate it
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      toast.success(t('mediaShare.settingsSaved'));
    } catch (error) {
      console.error('Error saving settings:', error);
      toast.error(t('mediaShare.settingsError'));
    } finally {
      setSaving(false);
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

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="card p-6">
        <div className="flex items-center mb-4">
          <Settings className="w-6 h-6 text-primary-600 mr-2" />
          <h2 className="text-xl font-semibold text-gray-900">
            {t('mediaShare.settingsTitle')}
          </h2>
        </div>
        <p className="text-gray-600">
          {t('mediaShare.settingsDesc')}
        </p>
      </div>

      {/* Basic Settings */}
      <div className="card p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">
          {t('mediaShare.basicSettings')}
        </h3>
        
        <div className="space-y-4">
          {/* Enable Media Share */}
          <div className="flex items-center justify-between">
            <div>
              <label className="text-sm font-medium text-gray-700">
                {t('mediaShare.enableMediaShare')}
              </label>
              <p className="text-xs text-gray-500">
                {t('mediaShare.enableMediaShareDesc')}
              </p>
            </div>
            <button
              onClick={() => setSettings(prev => ({ ...prev, mediaShareEnabled: !prev.mediaShareEnabled }))}
              className="flex items-center"
            >
              {settings.mediaShareEnabled ? (
                <ToggleRight className="w-8 h-8 text-green-600" />
              ) : (
                <ToggleLeft className="w-8 h-8 text-gray-400" />
              )}
            </button>
          </div>

          {/* Minimum Donation Amount */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {t('mediaShare.minDonationAmount')}
            </label>
            <div className="relative">
              <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
              <input
                type="number"
                value={settings.minDonationAmount}
                onChange={(e) => setSettings(prev => ({ ...prev, minDonationAmount: parseInt(e.target.value) || 0 }))}
                className="pl-10 w-full p-3 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
                min="0"
                step="1000"
              />
            </div>
            <p className="text-xs text-gray-500 mt-1">
              {t('mediaShare.currentAmount')}: {formatCurrency(settings.minDonationAmount)}
            </p>
          </div>

          {/* Currency */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {t('mediaShare.currency')}
            </label>
            <select
              value={settings.currency}
              onChange={(e) => setSettings(prev => ({ ...prev, currency: e.target.value }))}
              className="w-full p-3 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
            >
              <option value="IDR">Indonesian Rupiah (IDR)</option>
              <option value="USD">US Dollar (USD)</option>
              <option value="EUR">Euro (EUR)</option>
            </select>
          </div>
        </div>
      </div>

      {/* Platform Settings */}
      <div className="card p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">
          {t('mediaShare.platformSettings')}
        </h3>
        
        <div className="space-y-4">
          {/* Allow YouTube */}
          <div className="flex items-center justify-between">
            <div className="flex items-center">
              <Youtube className="w-5 h-5 text-red-600 mr-2" />
              <div>
                <label className="text-sm font-medium text-gray-700">
                  {t('mediaShare.allowYoutube')}
                </label>
                <p className="text-xs text-gray-500">
                  {t('mediaShare.allowYoutubeDesc')}
                </p>
              </div>
            </div>
            <button
              onClick={() => setSettings(prev => ({ ...prev, allowYoutube: !prev.allowYoutube }))}
              className="flex items-center"
            >
              {settings.allowYoutube ? (
                <ToggleRight className="w-8 h-8 text-green-600" />
              ) : (
                <ToggleLeft className="w-8 h-8 text-gray-400" />
              )}
            </button>
          </div>

          {/* YouTube Duration Limit */}
          {settings.allowYoutube && (
            <div className="ml-7">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                {t('mediaShare.maxDurationYoutube')} ({t('mediaShare.seconds')})
              </label>
              <input
                type="number"
                value={settings.maxDurationYoutube}
                onChange={(e) => setSettings(prev => ({ ...prev, maxDurationYoutube: parseInt(e.target.value) || 0 }))}
                className="w-full p-2 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
                min="30"
                max="600"
                step="30"
              />
            </div>
          )}

          {/* Allow TikTok */}
          <div className="flex items-center justify-between">
            <div className="flex items-center">
              <Music className="w-5 h-5 text-pink-600 mr-2" />
              <div>
                <label className="text-sm font-medium text-gray-700">
                  {t('mediaShare.allowTiktok')}
                </label>
                <p className="text-xs text-gray-500">
                  {t('mediaShare.allowTiktokDesc')}
                </p>
              </div>
            </div>
            <button
              onClick={() => setSettings(prev => ({ ...prev, allowTiktok: !prev.allowTiktok }))}
              className="flex items-center"
            >
              {settings.allowTiktok ? (
                <ToggleRight className="w-8 h-8 text-green-600" />
              ) : (
                <ToggleLeft className="w-8 h-8 text-gray-400" />
              )}
            </button>
          </div>

          {/* TikTok Duration Limit */}
          {settings.allowTiktok && (
            <div className="ml-7">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                {t('mediaShare.maxDurationTiktok')} ({t('mediaShare.seconds')})
              </label>
              <input
                type="number"
                value={settings.maxDurationTiktok}
                onChange={(e) => setSettings(prev => ({ ...prev, maxDurationTiktok: parseInt(e.target.value) || 0 }))}
                className="w-full p-2 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
                min="15"
                max="300"
                step="15"
              />
            </div>
          )}
        </div>
      </div>

      {/* Moderation Settings */}
      <div className="card p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">
          {t('mediaShare.moderationSettings')}
        </h3>
        
        <div className="space-y-4">
          {/* Auto Approve */}
          <div className="flex items-center justify-between">
            <div>
              <label className="text-sm font-medium text-gray-700">
                {t('mediaShare.autoApprove')}
              </label>
              <p className="text-xs text-gray-500">
                {t('mediaShare.autoApproveDesc')}
              </p>
            </div>
            <button
              onClick={() => setSettings(prev => ({ ...prev, autoApprove: !prev.autoApprove }))}
              className="flex items-center"
            >
              {settings.autoApprove ? (
                <ToggleRight className="w-8 h-8 text-green-600" />
              ) : (
                <ToggleLeft className="w-8 h-8 text-gray-400" />
              )}
            </button>
          </div>

          {/* Welcome Message */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {t('mediaShare.welcomeMessage')}
            </label>
            <textarea
              value={settings.welcomeMessage}
              onChange={(e) => setSettings(prev => ({ ...prev, welcomeMessage: e.target.value }))}
              rows={3}
              className="w-full p-3 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
              placeholder={t('mediaShare.welcomeMessagePlaceholder')}
            />
          </div>
        </div>
      </div>

      {/* Save Button */}
      <div className="card p-6">
        <button
          onClick={handleSave}
          disabled={saving}
          className="btn-primary flex items-center space-x-2 w-full justify-center"
        >
          <Save className="w-4 h-4" />
          <span>
            {saving ? t('mediaShare.saving') : t('mediaShare.saveSettings')}
          </span>
        </button>
      </div>
    </div>
  );
};

export default StreamerMediaSettings; 