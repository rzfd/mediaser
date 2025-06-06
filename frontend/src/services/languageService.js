import api from '../utils/api';

export const languageService = {
  // Translate text
  async translateText(text, from, to) {
    try {
      const response = await api.post('/api/language/translate', {
        text,
        from,
        to
      });
      return response.data;
    } catch (error) {
      console.error('Translation error:', error);
      throw new Error(error.response?.data?.message || 'Failed to translate text');
    }
  },

  // Detect language
  async detectLanguage(text) {
    try {
      const response = await api.post('/api/language/detect', {
        text
      });
      return response.data;
    } catch (error) {
      console.error('Language detection error:', error);
      throw new Error(error.response?.data?.message || 'Failed to detect language');
    }
  },

  // Bulk translate multiple texts
  async bulkTranslate(texts, from, to) {
    try {
      const response = await api.post('/api/language/bulk-translate', {
        texts,
        from,
        to
      });
      return response.data;
    } catch (error) {
      console.error('Bulk translation error:', error);
      throw new Error(error.response?.data?.message || 'Failed to translate texts');
    }
  },

  // Get system translations
  async getSystemTranslations(language) {
    try {
      const response = await api.get(`/api/language/system/${language}`);
      return response.data;
    } catch (error) {
      console.error('Get system translations error:', error);
      throw new Error(error.response?.data?.message || 'Failed to get system translations');
    }
  },

  // Get language preferences
  async getLanguagePreferences() {
    try {
      const response = await api.get('/api/language/preferences');
      return response.data;
    } catch (error) {
      console.error('Get language preferences error:', error);
      throw new Error(error.response?.data?.message || 'Failed to get language preferences');
    }
  },

  // Update language preferences
  async updateLanguagePreferences(preferences) {
    try {
      const response = await api.put('/api/language/preferences', preferences);
      return response.data;
    } catch (error) {
      console.error('Update language preferences error:', error);
      throw new Error(error.response?.data?.message || 'Failed to update language preferences');
    }
  },

  // Get supported languages
  async getSupportedLanguages() {
    try {
      const response = await api.get('/api/language/list');
      return response.data;
    } catch (error) {
      console.error('Get supported languages error:', error);
      // Return default languages if API fails
      return {
        languages: [
          { code: 'en', name: 'English', nativeName: 'English' },
          { code: 'id', name: 'Indonesian', nativeName: 'Bahasa Indonesia' },
          { code: 'zh', name: 'Chinese', nativeName: '中文' },
          { code: 'ja', name: 'Japanese', nativeName: '日本語' },
          { code: 'ko', name: 'Korean', nativeName: '한국어' },
          { code: 'th', name: 'Thai', nativeName: 'ไทย' },
          { code: 'vi', name: 'Vietnamese', nativeName: 'Tiếng Việt' },
          { code: 'ms', name: 'Malay', nativeName: 'Bahasa Melayu' },
          { code: 'tl', name: 'Filipino', nativeName: 'Filipino' },
          { code: 'es', name: 'Spanish', nativeName: 'Español' },
          { code: 'fr', name: 'French', nativeName: 'Français' },
          { code: 'de', name: 'German', nativeName: 'Deutsch' },
          { code: 'pt', name: 'Portuguese', nativeName: 'Português' },
          { code: 'ru', name: 'Russian', nativeName: 'Русский' },
          { code: 'ar', name: 'Arabic', nativeName: 'العربية' },
          { code: 'hi', name: 'Hindi', nativeName: 'हिन्दी' }
        ]
      };
    }
  },

  // Get language name by code
  getLanguageName(code, languages = []) {
    const defaultLanguages = {
      'en': 'English',
      'id': 'Indonesian',
      'zh': 'Chinese',
      'ja': 'Japanese',
      'ko': 'Korean',
      'th': 'Thai',
      'vi': 'Vietnamese',
      'ms': 'Malay',
      'tl': 'Filipino',
      'es': 'Spanish',
      'fr': 'French',
      'de': 'German',
      'pt': 'Portuguese',
      'ru': 'Russian',
      'ar': 'Arabic',
      'hi': 'Hindi'
    };

    if (languages.length > 0) {
      const language = languages.find(lang => lang.code === code);
      return language ? language.name : code.toUpperCase();
    }

    return defaultLanguages[code] || code.toUpperCase();
  },

  // Get language flag emoji by code
  getLanguageFlag(code) {
    const flags = {
      'en': '🇺🇸',
      'id': '🇮🇩',
      'zh': '🇨🇳',
      'ja': '🇯🇵',
      'ko': '🇰🇷',
      'th': '🇹🇭',
      'vi': '🇻🇳',
      'ms': '🇲🇾',
      'tl': '🇵🇭',
      'es': '🇪🇸',
      'fr': '🇫🇷',
      'de': '🇩🇪',
      'pt': '🇵🇹',
      'ru': '🇷🇺',
      'ar': '🇸🇦',
      'hi': '🇮🇳'
    };

    return flags[code] || '🌐';
  },

  // Check if text needs translation
  needsTranslation(text, targetLanguage, sourceLanguage = null) {
    if (!text || !targetLanguage) return false;
    if (sourceLanguage && sourceLanguage === targetLanguage) return false;
    
    // Simple heuristic: if text contains non-ASCII characters and target is English, probably needs translation
    // This is a basic implementation - in production you'd want more sophisticated detection
    const hasNonAscii = /[^\x00-\x7F]/.test(text);
    
    if (targetLanguage === 'en' && hasNonAscii) return true;
    if (targetLanguage !== 'en' && !hasNonAscii) return true;
    
    return false;
  },

  // Format translation confidence
  formatConfidence(confidence) {
    if (typeof confidence !== 'number') return 'Unknown';
    return `${Math.round(confidence * 100)}%`;
  }
};

export default languageService; 