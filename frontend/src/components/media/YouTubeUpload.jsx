import React, { useState, useContext } from 'react';
import { useTranslation } from 'react-i18next';
import { AuthContext } from '../../contexts/AuthContext';
import { Youtube, Link, Plus, Trash2, Play, Eye, Calendar } from 'lucide-react';
import { toast } from 'react-hot-toast';

const YouTubeUpload = ({ setLoading }) => {
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
      url: 'https://www.youtube.com/watch?v=dQw4w9WgXcQ',
      videoId: 'dQw4w9WgXcQ',
      title: 'Never Gonna Give You Up',
      description: 'Classic Rick Roll video',
      thumbnail: 'https://img.youtube.com/vi/dQw4w9WgXcQ/maxresdefault.jpg',
      uploadedAt: new Date().toISOString(),
      views: 0
    }
  ]);

  // Extract YouTube video ID from URL
  const extractVideoId = (url) => {
    const regex = /(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/watch\?v=|youtu\.be\/)([^&\n?#]+)/;
    const match = url.match(regex);
    return match ? match[1] : null;
  };

  // Validate YouTube URL
  const isValidYouTubeUrl = (url) => {
    return /^(https?:\/\/)?(www\.)?(youtube\.com\/watch\?v=|youtu\.be\/)[a-zA-Z0-9_-]+/.test(url);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!videoUrl.trim()) {
      toast.error(t('media.youtube.urlRequired'));
      return;
    }

    if (!isValidYouTubeUrl(videoUrl)) {
      toast.error(t('media.youtube.invalidUrl'));
      return;
    }

    const videoId = extractVideoId(videoUrl);
    if (!videoId) {
      toast.error(t('media.youtube.invalidUrl'));
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
        title: title || `YouTube Video ${videoId}`,
        description: description || '',
        thumbnail: `https://img.youtube.com/vi/${videoId}/maxresdefault.jpg`,
        uploadedAt: new Date().toISOString(),
        views: 0
      };

      setVideos(prev => [newVideo, ...prev]);
      
      // Reset form
      setVideoUrl('');
      setTitle('');
      setDescription('');
      
      toast.success(t('media.youtube.uploadSuccess'));
      
    } catch (error) {
      console.error('Error uploading YouTube video:', error);
      toast.error(t('media.youtube.uploadError'));
    } finally {
      setSubmitting(false);
      setLoading(false);
    }
  };

  const handleDelete = async (videoId) => {
    try {
      // Here you would make an API call to delete the video
      setVideos(prev => prev.filter(video => video.id !== videoId));
      toast.success(t('media.youtube.deleteSuccess'));
    } catch (error) {
      console.error('Error deleting video:', error);
      toast.error(t('media.youtube.deleteError'));
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
          <Youtube className="w-6 h-6 text-red-600 mr-2" />
          <h2 className="text-xl font-semibold text-gray-900">
            {t('media.youtube.addVideo')}
          </h2>
        </div>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {t('media.youtube.videoUrl')} *
            </label>
            <div className="relative">
              <Link className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
              <input
                type="url"
                value={videoUrl}
                onChange={(e) => setVideoUrl(e.target.value)}
                placeholder="https://www.youtube.com/watch?v=..."
                className="pl-10 w-full p-3 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
                required
              />
            </div>
            <p className="text-xs text-gray-500 mt-1">
              {t('media.youtube.urlHelper')}
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {t('media.youtube.customTitle')}
            </label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder={t('media.youtube.titlePlaceholder')}
              className="w-full p-3 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {t('media.youtube.description')}
            </label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder={t('media.youtube.descriptionPlaceholder')}
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
              {submitting ? t('media.youtube.adding') : t('media.youtube.addVideo')}
            </span>
          </button>
        </form>
      </div>

      {/* Videos List */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">
          {t('media.youtube.myVideos')} ({videos.length})
        </h3>
        
        {videos.length === 0 ? (
          <div className="text-center py-12">
            <Youtube className="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-500">
              {t('media.youtube.noVideos')}
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {videos.map((video) => (
              <div key={video.id} className="card overflow-hidden">
                <div className="relative">
                  <img
                    src={video.thumbnail}
                    alt={video.title}
                    className="w-full h-48 object-cover"
                  />
                  <div className="absolute inset-0 bg-black bg-opacity-50 opacity-0 hover:opacity-100 transition-opacity flex items-center justify-center space-x-4">
                    <button
                      onClick={() => openVideo(video.url)}
                      className="p-2 bg-white rounded-full text-gray-900 hover:bg-gray-100"
                      title={t('media.youtube.watch')}
                    >
                      <Play className="w-5 h-5" />
                    </button>
                    <button
                      onClick={() => handleDelete(video.id)}
                      className="p-2 bg-red-600 rounded-full text-white hover:bg-red-700"
                      title={t('media.youtube.delete')}
                    >
                      <Trash2 className="w-5 h-5" />
                    </button>
                  </div>
                </div>
                
                <div className="p-4">
                  <h4 className="font-semibold text-gray-900 mb-2 line-clamp-2">
                    {video.title}
                  </h4>
                  {video.description && (
                    <p className="text-sm text-gray-600 mb-3 line-clamp-2">
                      {video.description}
                    </p>
                  )}
                  <div className="flex items-center justify-between text-xs text-gray-500">
                    <div className="flex items-center space-x-1">
                      <Eye className="w-3 h-3" />
                      <span>{video.views} {t('media.views')}</span>
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

export default YouTubeUpload; 