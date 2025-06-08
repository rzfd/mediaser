import React, { useState, useContext } from 'react';
import { useTranslation } from 'react-i18next';
import { AuthContext } from '../../contexts/AuthContext';
import { Music, Link, Plus, Trash2, Play, Eye, Calendar, Heart } from 'lucide-react';
import { toast } from 'react-hot-toast';

const TikTokUpload = ({ setLoading }) => {
  const { t } = useTranslation();
  const { user } = useContext(AuthContext);
  const [videoUrl, setVideoUrl] = useState('');
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [videos, setVideos] = useState([
    // Mock data for demonstration
    {
      id: 1,
      url: 'https://www.tiktok.com/@username/video/1234567890',
      videoId: '1234567890',
      title: 'Amazing Dance Video',
      description: 'Check out this amazing dance!',
      thumbnail: 'https://via.placeholder.com/300x400/ff0050/ffffff?text=TikTok+Dance',
      uploadedAt: new Date().toISOString(),
      views: 1200,
      likes: 89,
      author: '@username'
    }
  ]);

  // Extract TikTok video ID from URL
  const extractVideoId = (url) => {
    const regex = /(?:https?:\/\/)?(?:www\.)?tiktok\.com\/@[^\/]+\/video\/(\d+)/;
    const match = url.match(regex);
    return match ? match[1] : null;
  };

  // Validate TikTok URL
  const isValidTikTokUrl = (url) => {
    return /^(https?:\/\/)?(www\.)?tiktok\.com\/@[^\/]+\/video\/\d+/.test(url);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!videoUrl.trim()) {
      toast.error(t('media.tiktok.urlRequired'));
      return;
    }

    if (!isValidTikTokUrl(videoUrl)) {
      toast.error(t('media.tiktok.invalidUrl'));
      return;
    }

    const videoId = extractVideoId(videoUrl);
    if (!videoId) {
      toast.error(t('media.tiktok.invalidUrl'));
      return;
    }

    setSubmitting(true);
    setLoading(true);

    try {
      // Here you would normally make an API call to save the video
      // For now, we'll simulate it and add to local state
      
      const newVideo = {
        id: Date.now(),
        url: videoUrl,
        videoId: videoId,
        title: title || `TikTok Video ${videoId}`,
        description: description || '',
        thumbnail: `https://via.placeholder.com/300x400/ff0050/ffffff?text=TikTok+Video`,
        uploadedAt: new Date().toISOString(),
        views: 0,
        likes: 0,
        author: videoUrl.match(/@([^\/]+)/)?.[1] || 'unknown'
      };

      setVideos(prev => [newVideo, ...prev]);
      
      // Reset form
      setVideoUrl('');
      setTitle('');
      setDescription('');
      
      toast.success(t('media.tiktok.uploadSuccess'));
      
    } catch (error) {
      console.error('Error uploading TikTok video:', error);
      toast.error(t('media.tiktok.uploadError'));
    } finally {
      setSubmitting(false);
      setLoading(false);
    }
  };

  const handleDelete = async (videoId) => {
    try {
      // Here you would make an API call to delete the video
      setVideos(prev => prev.filter(video => video.id !== videoId));
      toast.success(t('media.tiktok.deleteSuccess'));
    } catch (error) {
      console.error('Error deleting video:', error);
      toast.error(t('media.tiktok.deleteError'));
    }
  };

  const openVideo = (url) => {
    window.open(url, '_blank');
  };

  return (
    <div className="space-y-8">
      {/* Upload Form */}
      <div className="card p-6">
        <div className="flex items-center mb-4">
          <Music className="w-6 h-6 text-pink-600 mr-2" />
          <h2 className="text-xl font-semibold text-gray-900">
            {t('media.tiktok.addVideo')}
          </h2>
        </div>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {t('media.tiktok.videoUrl')} *
            </label>
            <div className="relative">
              <Link className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
              <input
                type="url"
                value={videoUrl}
                onChange={(e) => setVideoUrl(e.target.value)}
                placeholder="https://www.tiktok.com/@username/video/1234567890"
                className="pl-10 w-full p-3 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
                required
              />
            </div>
            <p className="text-xs text-gray-500 mt-1">
              {t('media.tiktok.urlHelper')}
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {t('media.tiktok.customTitle')}
            </label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder={t('media.tiktok.titlePlaceholder')}
              className="w-full p-3 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {t('media.tiktok.description')}
            </label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder={t('media.tiktok.descriptionPlaceholder')}
              rows={3}
              className="w-full p-3 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
            />
          </div>

          <button
            type="submit"
            disabled={submitting}
            className="btn-primary flex items-center space-x-2"
          >
            <Plus className="w-4 h-4" />
            <span>
              {submitting ? t('media.tiktok.adding') : t('media.tiktok.addVideo')}
            </span>
          </button>
        </form>
      </div>

      {/* Videos List */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">
          {t('media.tiktok.myVideos')} ({videos.length})
        </h3>
        
        {videos.length === 0 ? (
          <div className="text-center py-12">
            <Music className="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-500">
              {t('media.tiktok.noVideos')}
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
            {videos.map((video) => (
              <div key={video.id} className="card overflow-hidden">
                <div className="relative aspect-[3/4]">
                  <img
                    src={video.thumbnail}
                    alt={video.title}
                    className="w-full h-full object-cover"
                  />
                  <div className="absolute inset-0 bg-black bg-opacity-50 opacity-0 hover:opacity-100 transition-opacity flex items-center justify-center space-x-2">
                    <button
                      onClick={() => openVideo(video.url)}
                      className="p-2 bg-white rounded-full text-gray-900 hover:bg-gray-100"
                      title={t('media.tiktok.watch')}
                    >
                      <Play className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => handleDelete(video.id)}
                      className="p-2 bg-red-600 rounded-full text-white hover:bg-red-700"
                      title={t('media.tiktok.delete')}
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                  
                  {/* TikTok style overlay info */}
                  <div className="absolute bottom-0 left-0 right-0 p-3 bg-gradient-to-t from-black to-transparent text-white">
                    <p className="text-xs font-medium mb-1">@{video.author}</p>
                    <p className="text-xs line-clamp-2">{video.title}</p>
                  </div>
                </div>
                
                <div className="p-3">
                  <div className="flex items-center justify-between text-xs text-gray-500">
                    <div className="flex items-center space-x-3">
                      <div className="flex items-center space-x-1">
                        <Eye className="w-3 h-3" />
                        <span>{video.views}</span>
                      </div>
                      <div className="flex items-center space-x-1">
                        <Heart className="w-3 h-3" />
                        <span>{video.likes}</span>
                      </div>
                    </div>
                    <div className="flex items-center space-x-1">
                      <Calendar className="w-3 h-3" />
                      <span>{new Date(video.uploadedAt).toLocaleDateString()}</span>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default TikTokUpload; 