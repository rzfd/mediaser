import React from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Heart, Globe, DollarSign, Users, Zap, Shield, Youtube, Music } from 'lucide-react';

const HomePage = () => {
  const { t } = useTranslation();

  const features = [
    {
      icon: <Globe className="w-8 h-8 text-primary-600" />,
      title: t('language.title'),
      description: 'Support for multiple languages with real-time translation'
    },
    {
      icon: <DollarSign className="w-8 h-8 text-green-600" />,
      title: t('currency.title'),
      description: 'Multi-currency support with live exchange rates'
    },
    {
      icon: <Users className="w-8 h-8 text-blue-600" />,
      title: 'Global Community',
      description: 'Connect with streamers and creators worldwide'
    },
    {
      icon: <Zap className="w-8 h-8 text-yellow-600" />,
      title: 'Instant Donations',
      description: 'Fast and secure payment processing'
    },
    {
      icon: <Shield className="w-8 h-8 text-purple-600" />,
      title: 'Secure Platform',
      description: 'Bank-level security for all transactions'
    },
    {
      icon: <Heart className="w-8 h-8 text-red-600" />,
      title: 'Support Creators',
      description: 'Help your favorite creators grow and succeed'
    }
  ];

  return (
    <div className="space-y-16">
      {/* Hero Section */}
      <section className="text-center py-16">
        <div className="max-w-4xl mx-auto">
          <h1 className="text-5xl font-bold text-gray-900 mb-6">
            <span className="gradient-text">{t('app.title')}</span>
          </h1>
          <p className="text-xl text-gray-600 mb-8 leading-relaxed">
            {t('app.description')}
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              to="/donate"
              className="btn-primary text-lg px-8 py-3"
            >
              {t('navigation.donate')}
            </Link>
            <Link
              to="/streamers"
              className="btn-outline text-lg px-8 py-3"
            >
              {t('navigation.streamers')}
            </Link>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-16">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">
              Why Choose MediaShar?
            </h2>
            <p className="text-gray-600 max-w-2xl mx-auto">
              Experience the future of creator support with our innovative multi-language and multi-currency platform
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {features.map((feature, index) => (
              <div key={index} className="card card-hover p-6 text-center">
                <div className="flex justify-center mb-4">
                  {feature.icon}
                </div>
                <h3 className="text-xl font-semibold text-gray-900 mb-2">
                  {feature.title}
                </h3>
                <p className="text-gray-600">
                  {feature.description}
                </p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Tools Section */}
      <section className="py-16 bg-gray-50 -mx-4">
        <div className="max-w-6xl mx-auto px-4">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">
              Powerful Tools
            </h2>
            <p className="text-gray-600 max-w-2xl mx-auto">
              Access our integrated language and currency tools to break down barriers and connect globally
            </p>
          </div>
          
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            {/* Currency Converter */}
            <div className="card p-8">
              <div className="flex items-center mb-4">
                <DollarSign className="w-8 h-8 text-green-600 mr-3" />
                <h3 className="text-2xl font-semibold text-gray-900">
                  {t('currency.converter')}
                </h3>
              </div>
              <p className="text-gray-600 mb-6">
                Convert between multiple currencies with real-time exchange rates
              </p>
              <Link
                to="/currency"
                className="btn-primary"
              >
                Try Currency Converter
              </Link>
            </div>

            {/* Language Translator */}
            <div className="card p-8">
              <div className="flex items-center mb-4">
                <Globe className="w-8 h-8 text-primary-600 mr-3" />
                <h3 className="text-2xl font-semibold text-gray-900">
                  {t('language.translator')}
                </h3>
              </div>
              <p className="text-gray-600 mb-6">
                Translate text between multiple languages instantly
              </p>
              <Link
                to="/language"
                className="btn-primary"
              >
                Try Language Translator
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Media Section */}
      <section className="py-16">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">
              Share Your Content
            </h2>
            <p className="text-gray-600 max-w-2xl mx-auto">
              Upload and showcase your YouTube and TikTok videos to connect with your audience
            </p>
          </div>
          
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            {/* YouTube Integration */}
            <div className="card p-8">
              <div className="flex items-center mb-4">
                <Youtube className="w-8 h-8 text-red-600 mr-3" />
                <h3 className="text-2xl font-semibold text-gray-900">
                  YouTube Integration
                </h3>
              </div>
              <p className="text-gray-600 mb-6">
                Share your YouTube videos and let supporters discover your content
              </p>
              <Link
                to="/media"
                className="btn-primary"
              >
                Manage YouTube Videos
              </Link>
            </div>

            {/* TikTok Integration */}
            <div className="card p-8">
              <div className="flex items-center mb-4">
                <Music className="w-8 h-8 text-pink-600 mr-3" />
                <h3 className="text-2xl font-semibold text-gray-900">
                  TikTok Integration
                </h3>
              </div>
              <p className="text-gray-600 mb-6">
                Showcase your TikTok content and engage with your community
              </p>
              <Link
                to="/media"
                className="btn-primary"
              >
                Manage TikTok Videos
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-16 text-center">
        <div className="max-w-4xl mx-auto">
          <h2 className="text-3xl font-bold text-gray-900 mb-4">
            Ready to Support Your Favorite Creators?
          </h2>
          <p className="text-xl text-gray-600 mb-8">
            Join thousands of supporters making a difference in the creator economy
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              to="/register"
              className="btn-primary text-lg px-8 py-3"
            >
              {t('navigation.register')}
            </Link>
            <Link
              to="/streamers"
              className="btn-secondary text-lg px-8 py-3"
            >
              Browse Creators
            </Link>
          </div>
        </div>
      </section>
    </div>
  );
};

export default HomePage; 