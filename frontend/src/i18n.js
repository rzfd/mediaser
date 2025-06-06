import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

import translationEN from './locales/en/translation.json';
import translationID from './locales/id/translation.json';
import translationZH from './locales/zh/translation.json';

const resources = {
  en: {
    translation: translationEN
  },
  id: {
    translation: translationID
  },
  zh: {
    translation: translationZH
  }
};

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: 'en',
    debug: process.env.NODE_ENV === 'development',

    interpolation: {
      escapeValue: false, // not needed for react as it escapes by default
    },

    detection: {
      order: ['localStorage', 'navigator', 'htmlTag'],
      lookupLocalStorage: 'mediashar-language',
      caches: ['localStorage'],
    },

    react: {
      useSuspense: false
    }
  });

export default i18n; 