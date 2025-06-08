import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Youtube, Music, Play, Eye, Calendar, Heart, Filter, SortAsc } from 'lucide-react';

const MediaGallery = ({ searchTerm, viewMode, setLoading }) => {
  const { t } = useTranslation();
  const [allMedia, setAllMedia] = useState([]);
  const [filteredMedia, setFilteredMedia] = useState([]);
  const [sortBy, setSortBy] = useState('date');
  const [filterType, setFilterType] = useState('all');

  // Mock data combining YouTube and TikTok videos
  const mockMedia = [
    {
      id: 1,
      type: 'youtube',
      url: 'https://www.youtube.com/watch?v=dQw4w9WgXcQ',
      videoId: 'dQw4w9WgXcQ',
      title: 'Never Gonna Give You Up',
      description: 'Classic Rick Roll video',
      thumbnail: 'https://img.youtube.com/vi/dQw4w9WgXcQ/maxresdefault.jpg',
      uploadedAt: new Date('2024-01-15').toISOString(),
      views: 1500,
      likes: 120,
      author: 'Rick Astley'
    },
    {
      id: 2,
      type: 'tiktok',
      url: 'https://www.tiktok.com/@username/video/1234567890',
      videoId: '1234567890',
      title: 'Amazing Dance Video',
      description: 'Check out this amazing dance!',
      thumbnail: 'https://via.placeholder.com/300x400/ff0050/ffffff?text=TikTok+Dance',
      uploadedAt: new Date('2024-01-20').toISOString(),
      views: 2300,
      likes: 189,
      author: '@username'
    },
    {
      id: 3,
      type: 'youtube',
      url: 'https://www.youtube.com/watch?v=oHg5SJYRHA0',
      videoId: 'oHg5SJYRHA0',
      title: 'RickRoll\'d',
      description: 'Another classic video',
      thumbnail: 'https://img.youtube.com/vi/oHg5SJYRHA0/maxresdefault.jpg',
      uploadedAt: new Date('2024-01-10').toISOString(),
      views: 890,
      likes: 67,
      author: 'Various Artists'
    },
    {
      id: 4,
      type: 'tiktok',
      url: 'https://www.tiktok.com/@creator/video/9876543210',
      videoId: '9876543210',
      title: 'Cooking Tips & Tricks',
      description: 'Learn amazing cooking techniques!',
      thumbnail: 'https://via.placeholder.com/300x400/25f4ee/000000?text=Cooking+Tips',
      uploadedAt: new Date('2024-01-18').toISOString(),
      views: 3200,
      likes: 245,
      author: '@creator'
    }
  ];

  useEffect(() => {
    // In a real app, you would fetch media from your API here
    setLoading(true);
    setTimeout(() => {
      setAllMedia(mockMedia);
      setLoading(false);
    }, 1000);
  }, [setLoading]);

  useEffect(() => {
    let filtered = allMedia;

    // Filter by type
    if (filterType !== 'all') {
      filtered = filtered.filter(media => media.type === filterType);
    }

    // Filter by search term
    if (searchTerm) {
      filtered = filtered.filter(media => 
        media.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
        media.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
        media.author.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    // Sort
    switch (sortBy) {
      case 'date':
        filtered.sort((a, b) => new Date(b.uploadedAt) - new Date(a.uploadedAt));
        break;
      case 'views':
        filtered.sort((a, b) => b.views - a.views);
        break;
      case 'likes':
        filtered.sort((a, b) => b.likes - a.likes);
        break;
      case 'title':
        filtered.sort((a, b) => a.title.localeCompare(b.title));
        break;
      default:
        break;
    }

    setFilteredMedia(filtered);
  }, [allMedia, searchTerm, sortBy, filterType]);

  const openVideo = (url) => {
    window.open(url, '_blank');
  };

  const getMediaIcon = (type) => {
    return type === 'youtube' ? 
      <Youtube className="w-4 h-4 text-red-600" /> : 
      <Music className="w-4 h-4 text-pink-600" />;
  };

  const getMediaBadge = (type) => {
    return type === 'youtube' ? 
      <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-red-100 text-red-800">
        <Youtube className="w-3 h-3 mr-1" />
        YouTube
      </span> :
      <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-pink-100 text-pink-800">
        <Music className="w-3 h-3 mr-1" />
        TikTok
      </span>;
  };

  return (
    <div className="space-y-6">
      {/* Controls */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-4 sm:space-y-0">
        <div className="flex items-center space-x-4">
          {/* Filter by type */}
          <div className="flex items-center space-x-2">
            <Filter className="w-4 h-4 text-gray-500" />
            <select
              value={filterType}
              onChange={(e) => setFilterType(e.target.value)}
              className="border border-gray-300 rounded-md px-3 py-1 text-sm focus:ring-primary-500 focus:border-primary-500"
            >
              <option value="all">{t('media.filter.all')}</option>
              <option value="youtube">{t('media.filter.youtube')}</option>
              <option value="tiktok">{t('media.filter.tiktok')}</option>
            </select>
          </div>

          {/* Sort */}
          <div className="flex items-center space-x-2">
            <SortAsc className="w-4 h-4 text-gray-500" />
            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
              className="border border-gray-300 rounded-md px-3 py-1 text-sm focus:ring-primary-500 focus:border-primary-500"
            >
              <option value="date">{t('media.sort.date')}</option>
              <option value="views">{t('media.sort.views')}</option>
              <option value="likes">{t('media.sort.likes')}</option>
              <option value="title">{t('media.sort.title')}</option>
            </select>
          </div>
        </div>

        <div className="text-sm text-gray-500">
          {filteredMedia.length} {t('media.results')}
        </div>
      </div>

      {/* Media Grid/List */}
      {filteredMedia.length === 0 ? (
        <div className="text-center py-12">
          <div className="flex justify-center space-x-2 mb-4">
            <Youtube className="w-8 h-8 text-gray-400" />
            <Music className="w-8 h-8 text-gray-400" />
          </div>
          <p className="text-gray-500">
            {searchTerm ? t('media.noResults') : t('media.noMedia')}
          </p>
        </div>
      ) : (
        <div className={
          viewMode === 'grid' 
            ? 'grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6'
            : 'space-y-4'
        }>
          {filteredMedia.map((media) => (
            <div 
              key={`${media.type}-${media.id}`} 
              className={`card overflow-hidden cursor-pointer hover:shadow-lg transition-shadow ${
                viewMode === 'list' ? 'flex items-center space-x-4 p-4' : ''
              }`}
              onClick={() => openVideo(media.url)}
            >
              {viewMode === 'grid' ? (
                <>
                  <div className="relative">
                    <img
                      src={media.thumbnail}
                      alt={media.title}
                      className={`w-full object-cover ${
                        media.type === 'tiktok' ? 'h-64' : 'h-48'
                      }`}
                    />
                    <div className="absolute inset-0 bg-black bg-opacity-0 hover:bg-opacity-50 transition-opacity flex items-center justify-center">
                      <Play className="w-12 h-12 text-white opacity-0 hover:opacity-100 transition-opacity" />
                    </div>
                    <div className="absolute top-2 left-2">
                      {getMediaBadge(media.type)}
                    </div>
                  </div>
                  
                  <div className="p-4">
                    <h4 className="font-semibold text-gray-900 mb-2 line-clamp-2">
                      {media.title}
                    </h4>
                    <p className="text-sm text-gray-600 mb-2">
                      {media.author}
                    </p>
                    <p className="text-sm text-gray-600 mb-3 line-clamp-2">
                      {media.description}
                    </p>
                    <div className="flex items-center justify-between text-xs text-gray-500">
                      <div className="flex items-center space-x-3">
                        <div className="flex items-center space-x-1">
                          <Eye className="w-3 h-3" />
                          <span>{media.views}</span>
                        </div>
                        <div className="flex items-center space-x-1">
                          <Heart className="w-3 h-3" />
                          <span>{media.likes}</span>
                        </div>
                      </div>
                      <div className="flex items-center space-x-1">
                        <Calendar className="w-3 h-3" />
                        <span>{new Date(media.uploadedAt).toLocaleDateString()}</span>
                      </div>
                    </div>
                  </div>
                </>
              ) : (
                <>
                  <div className="relative flex-shrink-0">
                    <img
                      src={media.thumbnail}
                      alt={media.title}
                      className="w-24 h-16 object-cover rounded"
                    />
                    <div className="absolute inset-0 flex items-center justify-center">
                      <Play className="w-6 h-6 text-white opacity-75" />
                    </div>
                  </div>
                  
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center space-x-2 mb-1">
                          {getMediaIcon(media.type)}
                          <h4 className="font-semibold text-gray-900 truncate">
                            {media.title}
                          </h4>
                        </div>
                        <p className="text-sm text-gray-600 mb-1">
                          {media.author}
                        </p>
                        <p className="text-sm text-gray-600 line-clamp-1">
                          {media.description}
                        </p>
                      </div>
                      
                      <div className="flex-shrink-0 text-right ml-4">
                        <div className="flex items-center space-x-3 text-xs text-gray-500 mb-1">
                          <div className="flex items-center space-x-1">
                            <Eye className="w-3 h-3" />
                            <span>{media.views}</span>
                          </div>
                          <div className="flex items-center space-x-1">
                            <Heart className="w-3 h-3" />
                            <span>{media.likes}</span>
                          </div>
                        </div>
                        <div className="flex items-center space-x-1 text-xs text-gray-500">
                          <Calendar className="w-3 h-3" />
                          <span>{new Date(media.uploadedAt).toLocaleDateString()}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default MediaGallery; 