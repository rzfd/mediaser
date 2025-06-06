import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { ChevronDown, Check } from 'lucide-react';
import { languageService } from '../../services/languageService';

const LanguageSelector = () => {
  const { i18n, t } = useTranslation();
  const [isOpen, setIsOpen] = useState(false);
  const [languages, setLanguages] = useState([]);

  const defaultLanguages = [
    { code: 'en', name: 'English', nativeName: 'English', flag: '🇺🇸' },
    { code: 'id', name: 'Indonesian', nativeName: 'Bahasa Indonesia', flag: '🇮🇩' },
    { code: 'zh', name: 'Chinese', nativeName: '中文', flag: '🇨🇳' }
  ];

  useEffect(() => {
    const loadLanguages = async () => {
      try {
        const response = await languageService.getSupportedLanguages();
        const supportedLanguages = response.languages.map(lang => ({
          ...lang,
          flag: languageService.getLanguageFlag(lang.code)
        }));
        setLanguages(supportedLanguages);
      } catch (error) {
        console.warn('Failed to load languages from API, using defaults:', error);
        setLanguages(defaultLanguages);
      }
    };

    loadLanguages();
  }, []);

  const currentLanguage = languages.find(lang => lang.code === i18n.language) || 
                         defaultLanguages.find(lang => lang.code === i18n.language) ||
                         defaultLanguages[0];

  const handleLanguageChange = async (languageCode) => {
    try {
      await i18n.changeLanguage(languageCode);
      setIsOpen(false);
      
      // Save preference to localStorage
      localStorage.setItem('mediashar-language', languageCode);
      
      // Update document language for accessibility
      document.documentElement.lang = languageCode;
      
    } catch (error) {
      console.error('Failed to change language:', error);
    }
  };

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (isOpen && !event.target.closest('.language-selector')) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [isOpen]);

  return (
    <div className="relative language-selector">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center space-x-2 px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 transition-colors duration-200"
        aria-label={t('language.select')}
        aria-expanded={isOpen}
        aria-haspopup="listbox"
      >
        <span className="text-lg">{currentLanguage?.flag || '🌐'}</span>
        <span className="hidden sm:inline-block">
          {currentLanguage?.nativeName || currentLanguage?.name || 'English'}
        </span>
        <ChevronDown 
          className={`w-4 h-4 transition-transform duration-200 ${
            isOpen ? 'rotate-180' : ''
          }`} 
        />
      </button>

      {isOpen && (
        <div className="absolute right-0 mt-2 w-56 bg-white border border-gray-200 rounded-lg shadow-lg z-50 max-h-60 overflow-y-auto">
          <div className="py-2">
            <div className="px-3 py-2 text-xs font-medium text-gray-500 uppercase tracking-wider border-b border-gray-100">
              {t('language.select')}
            </div>
            
            {languages.map((language) => (
              <button
                key={language.code}
                onClick={() => handleLanguageChange(language.code)}
                className={`w-full flex items-center justify-between px-3 py-2 text-sm hover:bg-gray-50 transition-colors duration-200 ${
                  i18n.language === language.code
                    ? 'bg-primary-50 text-primary-700'
                    : 'text-gray-900'
                }`}
                role="option"
                aria-selected={i18n.language === language.code}
              >
                <div className="flex items-center space-x-3">
                  <span className="text-lg">{language.flag || '🌐'}</span>
                  <div className="text-left">
                    <div className="font-medium">{language.nativeName}</div>
                    <div className="text-xs text-gray-500">{language.name}</div>
                  </div>
                </div>
                
                {i18n.language === language.code && (
                  <Check className="w-4 h-4 text-primary-600" />
                )}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default LanguageSelector; 