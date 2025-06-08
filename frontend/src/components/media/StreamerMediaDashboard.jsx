import React, { useState, useContext, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { AuthContext } from '../../contexts/AuthContext';
import { 
  Monitor, 
  Youtube, 
  Music, 
  Play, 
  Check, 
  X, 
  Clock, 
  DollarSign,
  User,
  Calendar,
  Filter,
  Search,
  Eye
} from 'lucide-react';
import { toast } from 'react-hot-toast';

const StreamerMediaDashboard = () => {
  const { t } = useTranslation();
  const { user } = useContext(AuthContext);
  const [mediaQueue, setMediaQueue] = useState([]);
  const [filteredQueue, setFilteredQueue] = useState([]);
  const [filterStatus, setFilterStatus] = useState('all');
  const [searchTerm, setSearchTerm] = useState('');
  const [loading, setLoading] = useState(true);

  // Mock data for demonstration
  const mockMediaQueue = [
    {
      id: 1,
      type: 'youtube',
      url: 'https://www.youtube.com/watch?v=dQw4w9WgXcQ',
      title: 'Never Gonna Give You Up',
      message: 'Classic song untuk streamer favorit!',
      donatorName: 'music_lover123',
      donationAmount: 10000,
      currency: 'IDR',
      status: 'pending',
      submittedAt: new Date(Date.now() - 1000 * 60 * 30).toISOString(),
      thumbnail: 'https://img.youtube.com/vi/dQw4w9WgXcQ/maxresdefault.jpg'
    },
    {
      id: 2,
      type: 'tiktok',
      url: 'https://www.tiktok.com/@username/video/1234567890',
      title: 'Funny Dance Challenge',
      message: 'Ini lucu banget, harus ditonton!',
      donatorName: 'dancer_pro',
      donationAmount: 15000,
      currency: 'IDR',
      status: 'approved',
      submittedAt: new Date(Date.now() - 1000 * 60 * 60).toISOString(),
      thumbnail: 'https://via.placeholder.com/300x400/ff0050/ffffff?text=TikTok+Dance'
    }
  ];

  useEffect(() => {
    setTimeout(() => {
      setMediaQueue(mockMediaQueue);
      setLoading(false);
    }, 1000);
  }, []);

  useEffect(() => {
    let filtered = mediaQueue;

    if (filterStatus !== 'all') {
      filtered = filtered.filter(item => item.status === filterStatus);
    }

    if (searchTerm) {
      filtered = filtered.filter(item => 
        item.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
        item.donatorName.toLowerCase().includes(searchTerm.toLowerCase()) ||
        item.message.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    setFilteredQueue(filtered);
  }, [mediaQueue, filterStatus, searchTerm]);

  if (!user || user.userType !== 'streamer') {
    return (
      <div className="card p-8 text-center">
        <Monitor className="w-16 h-16 text-gray-400 mx-auto mb-4" />
        <h2 className="text-xl font-semibold text-gray-900 mb-2">
          {t('mediaShare.streamerOnly')}
        </h2>
        <p className="text-gray-600">
          {t('mediaShare.streamerOnlyDesc')}
        </p>
      </div>
    );
  }

  const handleApprove = async (mediaId) => {
    try {
      setMediaQueue(prev => 
        prev.map(item => 
          item.id === mediaId 
            ? { ...item, status: 'approved' }
            : item
        )
      );
      toast.success(t('mediaShare.approved'));
    } catch (error) {
      toast.error(t('mediaShare.approveError'));
    }
  };

  const handleReject = async (mediaId) => {
    try {
      setMediaQueue(prev => 
        prev.map(item => 
          item.id === mediaId 
            ? { ...item, status: 'rejected' }
            : item
        )
      );
      toast.success(t('mediaShare.rejected'));
    } catch (error) {
      toast.error(t('mediaShare.rejectError'));
    }
  };

  if (loading) {
    return (
      <div className="card p-8 text-center">
        <div className="animate-spin w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full mx-auto mb-4"></div>
        <p className="text-gray-600">{t('common.loading')}</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="card p-6">
        <div className="flex items-center mb-4">
          <Monitor className="w-6 h-6 text-primary-600 mr-2" />
          <h2 className="text-xl font-semibold text-gray-900">
            {t('mediaShare.dashboard')}
          </h2>
        </div>
        <p className="text-gray-600">
          {t('mediaShare.dashboardDesc')}
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="card p-4">
          <div className="flex items-center">
            <Clock className="w-8 h-8 text-yellow-600 mr-3" />
            <div>
              <p className="text-sm text-gray-600">{t('mediaShare.pending')}</p>
              <p className="text-xl font-semibold text-gray-900">
                {mediaQueue.filter(m => m.status === 'pending').length}
              </p>
            </div>
          </div>
        </div>
        <div className="card p-4">
          <div className="flex items-center">
            <Check className="w-8 h-8 text-green-600 mr-3" />
            <div>
              <p className="text-sm text-gray-600">{t('mediaShare.approved')}</p>
              <p className="text-xl font-semibold text-gray-900">
                {mediaQueue.filter(m => m.status === 'approved').length}
              </p>
            </div>
          </div>
        </div>
        <div className="card p-4">
          <div className="flex items-center">
            <X className="w-8 h-8 text-red-600 mr-3" />
            <div>
              <p className="text-sm text-gray-600">{t('mediaShare.rejected')}</p>
              <p className="text-xl font-semibold text-gray-900">
                {mediaQueue.filter(m => m.status === 'rejected').length}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div className="card p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">
          {t('mediaShare.mediaQueue')} ({filteredQueue.length})
        </h3>
        
        {filteredQueue.length === 0 ? (
          <div className="text-center py-12">
            <Monitor className="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-500">
              {t('mediaShare.noMediaInQueue')}
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {filteredQueue.map((media) => (
              <div key={media.id} className="border border-gray-200 rounded-lg p-4">
                <div className="flex items-start space-x-4">
                  <div className="relative flex-shrink-0">
                    <img
                      src={media.thumbnail}
                      alt={media.title}
                      className="w-20 h-16 object-cover rounded"
                    />
                  </div>
                  
                  <div className="flex-1">
                    <div className="flex items-center space-x-2 mb-1">
                      {media.type === 'youtube' ? (
                        <Youtube className="w-4 h-4 text-red-600" />
                      ) : (
                        <Music className="w-4 h-4 text-pink-600" />
                      )}
                      <h4 className="font-semibold text-gray-900">
                        {media.title}
                      </h4>
                    </div>
                    
                    <p className="text-sm text-gray-600 mb-2">
                      {media.message}
                    </p>
                    
                    <div className="flex items-center space-x-4 text-xs text-gray-500">
                      <span>{media.donatorName}</span>
                      <span>Rp {media.donationAmount.toLocaleString()}</span>
                    </div>
                  </div>
                  
                  {media.status === 'pending' && (
                    <div className="flex space-x-2">
                      <button
                        onClick={() => handleApprove(media.id)}
                        className="p-2 text-green-600 hover:bg-green-100 rounded-lg"
                      >
                        <Check className="w-4 h-4" />
                      </button>
                      <button
                        onClick={() => handleReject(media.id)}
                        className="p-2 text-red-600 hover:bg-red-100 rounded-lg"
                      >
                        <X className="w-4 h-4" />
                      </button>
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default StreamerMediaDashboard; 