import React, { useState, useContext } from 'react';
import { useTranslation } from 'react-i18next';
import { AuthContext } from '../contexts/AuthContext';
import { Youtube, Music, Upload, Play, Search, Grid, List } from 'lucide-react';
import YouTubeUpload from '../components/media/YouTubeUpload';
import TikTokUpload from '../components/media/TikTokUpload';
import MediaGallery from '../components/media/MediaGallery';
import LoadingSpinner from '../components/ui/LoadingSpinner';

const MediaPage = () => {
  const { t } = useTranslation();
  const { user } = useContext(AuthContext);
  const [activeTab, setActiveTab] = useState('gallery');
  const [viewMode, setViewMode] = useState('grid');
  const [searchTerm, setSearchTerm] = useState('');
  const [loading, setLoading] = useState(false);

  const tabs = [
    { id: 'gallery', name: t('media.gallery'), icon: Play },
    { id: 'youtube', name: t('media.youtube'), icon: Youtube },
    { id: 'tiktok', name: t('media.tiktok'), icon: Music }
  ];

  if (!user) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center py-16">
          <Upload className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h2 className="text-2xl font-bold text-gray-900 mb-4">
            {t('media.loginRequired')}
          </h2>
          <p className="text-gray-600 mb-6">
            {t('media.loginRequiredDesc')}
          </p>
          <a
            href="/login"
            className="btn-primary"
          >
            {t('navigation.login')}
          </a>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          {t('media.title')}
        </h1>
        <p className="text-gray-600">
          {t('media.description')}
        </p>
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200 mb-6">
        <nav className="-mb-px flex space-x-8">
          {tabs.map((tab) => {
            const Icon = tab.icon;
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`py-2 px-1 border-b-2 font-medium text-sm whitespace-nowrap flex items-center space-x-2 ${
                  activeTab === tab.id
                    ? 'border-primary-500 text-primary-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                <Icon className="w-4 h-4" />
                <span>{tab.name}</span>
              </button>
            );
          })}
        </nav>
      </div>

      {/* Search and View Controls */}
      {activeTab === 'gallery' && (
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-6 space-y-4 sm:space-y-0">
          <div className="relative flex-1 max-w-md">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
            <input
              type="text"
              placeholder={t('media.search')}
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-10 pr-4 py-2 w-full border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
            />
          </div>
          <div className="flex items-center space-x-2">
            <button
              onClick={() => setViewMode('grid')}
              className={`p-2 rounded-md ${
                viewMode === 'grid'
                  ? 'bg-primary-100 text-primary-600'
                  : 'text-gray-400 hover:text-gray-600'
              }`}
            >
              <Grid className="w-4 h-4" />
            </button>
            <button
              onClick={() => setViewMode('list')}
              className={`p-2 rounded-md ${
                viewMode === 'list'
                  ? 'bg-primary-100 text-primary-600'
                  : 'text-gray-400 hover:text-gray-600'
              }`}
            >
              <List className="w-4 h-4" />
            </button>
          </div>
        </div>
      )}

      {/* Content */}
      {loading && (
        <div className="flex justify-center py-8">
          <LoadingSpinner />
        </div>
      )}

      {!loading && (
        <div>
          {activeTab === 'gallery' && (
            <MediaGallery 
              searchTerm={searchTerm}
              viewMode={viewMode}
              setLoading={setLoading}
            />
          )}
          {activeTab === 'youtube' && (
            <YouTubeUpload setLoading={setLoading} />
          )}
          {activeTab === 'tiktok' && (
            <TikTokUpload setLoading={setLoading} />
          )}
        </div>
      )}
    </div>
  );
};

export default MediaPage; 