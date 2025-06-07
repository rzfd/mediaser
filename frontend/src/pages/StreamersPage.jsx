import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import { Search, Users, Heart, DollarSign, User } from 'lucide-react';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import { apiRequest } from '../utils/tokenUtils';

const StreamersPage = () => {
  const { t } = useTranslation();
  const [streamers, setStreamers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [error, setError] = useState(null);

  // API base URL
  const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

  // Fetch streamers from API
  useEffect(() => {
    fetchStreamers();
  }, []);

  const fetchStreamers = async () => {
    try {
      setLoading(true);
      
      const response = await apiRequest(`${API_BASE_URL}/streamers`, {
        method: 'GET'
      });
      
      if (!response.ok) {
        throw new Error('Failed to fetch streamers');
      }
      
      const data = await response.json();
      setStreamers(data.data || data || []);
    } catch (err) {
      console.error('Error fetching streamers:', err);
      setError(err.message);
      // Set some mock data for demo purposes
      setStreamers([
        {
          id: 1,
          username: 'gaming_master',
          full_name: 'Gaming Master',
          description: 'Professional gamer and content creator',
          total_donations: 150000,
          donation_count: 45,
          avatar_url: null
        },
        {
          id: 2,
          username: 'music_lover',
          full_name: 'Music Lover',
          description: 'Singer and music producer',
          total_donations: 85000,
          donation_count: 28,
          avatar_url: null
        }
      ]);
    } finally {
      setLoading(false);
    }
  };

  // Filter streamers based on search term
  const filteredStreamers = streamers.filter(streamer =>
    streamer.username.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (streamer.full_name && streamer.full_name.toLowerCase().includes(searchTerm.toLowerCase()))
  );

  // Format currency
  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0
    }).format(amount || 0);
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <LoadingSpinner size="xl" />
          <p className="mt-4 text-gray-600">Loading streamers...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto">
      {/* Header */}
      <div className="text-center py-8">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          {t('streamers.title')}
        </h1>
        <p className="text-xl text-gray-600 mb-8">
          Discover and support amazing creators from around the world
        </p>
      </div>

      {/* Search Bar */}
      <div className="mb-8">
        <div className="relative max-w-md mx-auto">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
          <input
            type="text"
            placeholder={t('streamers.search')}
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>
      </div>

      {/* Error Message */}
      {error && (
        <div className="mb-6 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
          <p className="text-yellow-800">
            ⚠️ Could not load data from server. Showing demo data.
          </p>
        </div>
      )}

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <Users className="w-8 h-8 text-blue-600 mr-3" />
            <div>
              <h3 className="text-lg font-semibold text-gray-900">Total Streamers</h3>
              <p className="text-2xl font-bold text-blue-600">{filteredStreamers.length}</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <Heart className="w-8 h-8 text-red-600 mr-3" />
            <div>
              <h3 className="text-lg font-semibold text-gray-900">Total Donations</h3>
              <p className="text-2xl font-bold text-red-600">
                {filteredStreamers.reduce((sum, s) => sum + (s.donation_count || 0), 0)}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <DollarSign className="w-8 h-8 text-green-600 mr-3" />
            <div>
              <h3 className="text-lg font-semibold text-gray-900">Total Amount</h3>
              <p className="text-2xl font-bold text-green-600">
                {formatCurrency(filteredStreamers.reduce((sum, s) => sum + (s.total_donations || 0), 0))}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Streamers Grid */}
      {filteredStreamers.length === 0 ? (
        <div className="text-center py-12">
          <Users className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            {t('streamers.noResults')}
          </h3>
          <p className="text-gray-500">
            {searchTerm ? 'Try adjusting your search terms' : 'No streamers available at the moment'}
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredStreamers.map((streamer) => (
            <div key={streamer.id} className="bg-white rounded-lg shadow hover:shadow-lg transition-shadow duration-200">
              {/* Streamer Avatar */}
              <div className="p-6">
                <div className="flex items-center mb-4">
                  <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center">
                    {streamer.avatar_url ? (
                      <img
                        src={streamer.avatar_url}
                        alt={streamer.username}
                        className="w-16 h-16 rounded-full object-cover"
                      />
                    ) : (
                      <User className="w-8 h-8 text-white" />
                    )}
                  </div>
                  <div className="ml-4 flex-1">
                    <h3 className="text-lg font-semibold text-gray-900">
                      {streamer.full_name || streamer.username}
                    </h3>
                    <p className="text-gray-500">@{streamer.username}</p>
                  </div>
                </div>

                {/* Description */}
                {streamer.description && (
                  <p className="text-gray-600 mb-4 text-sm">
                    {streamer.description}
                  </p>
                )}

                {/* Stats */}
                <div className="grid grid-cols-2 gap-4 mb-4">
                  <div className="text-center">
                    <p className="text-sm text-gray-500">Total Donations</p>
                    <p className="text-lg font-semibold text-gray-900">
                      {formatCurrency(streamer.total_donations || 0)}
                    </p>
                  </div>
                  <div className="text-center">
                    <p className="text-sm text-gray-500">Supporters</p>
                    <p className="text-lg font-semibold text-gray-900">
                      {streamer.donation_count || 0}
                    </p>
                  </div>
                </div>

                {/* Support Button */}
                <Link
                  to={`/donate/${streamer.id}`}
                  className="w-full bg-gradient-to-r from-blue-500 to-purple-600 text-white py-2 px-4 rounded-lg font-medium hover:from-blue-600 hover:to-purple-700 transition-colors duration-200 flex items-center justify-center"
                >
                  <Heart className="w-4 h-4 mr-2" />
                  {t('streamers.supportButton')}
                </Link>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Load More Button (for future pagination) */}
      {filteredStreamers.length > 0 && (
        <div className="text-center mt-8">
          <button
            className="bg-gray-200 text-gray-700 px-6 py-2 rounded-lg hover:bg-gray-300 transition-colors duration-200"
            disabled
          >
            Load More Streamers
          </button>
        </div>
      )}
    </div>
  );
};

export default StreamersPage; 