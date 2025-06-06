-- Language Database Initialization Script
-- This script sets up the schema for translations and language preferences

-- Create database and user if they don't exist
-- (This is handled by docker-entrypoint-initdb.d)

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create language_configs table for system translations
CREATE TABLE IF NOT EXISTS language_configs (
    id SERIAL PRIMARY KEY,
    language VARCHAR(10) NOT NULL,
    key VARCHAR(255) NOT NULL,
    translation TEXT NOT NULL,
    category VARCHAR(100) DEFAULT 'general',
    module VARCHAR(100) DEFAULT 'system',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure unique key per language
    UNIQUE(language, key)
);

-- Create language_info table
CREATE TABLE IF NOT EXISTS language_info (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    native_name VARCHAR(100) NOT NULL,
    flag VARCHAR(10),
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_rtl BOOLEAN NOT NULL DEFAULT false,
    date_format VARCHAR(50) DEFAULT 'DD/MM/YYYY',
    time_format VARCHAR(50) DEFAULT 'HH:mm',
    number_format VARCHAR(50) DEFAULT '1,234.56',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create user_language_preferences table
CREATE TABLE IF NOT EXISTS user_language_preferences (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    primary_language VARCHAR(10) NOT NULL DEFAULT 'id',
    fallback_language VARCHAR(10) NOT NULL DEFAULT 'en',
    auto_detect BOOLEAN NOT NULL DEFAULT true,
    timezone VARCHAR(100) DEFAULT 'Asia/Jakarta',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure one preference per user
    UNIQUE(user_id)
);

-- Create translation_cache table for external API results
CREATE TABLE IF NOT EXISTS translation_cache (
    id SERIAL PRIMARY KEY,
    original_text TEXT NOT NULL,
    from_lang VARCHAR(10) NOT NULL,
    to_lang VARCHAR(10) NOT NULL,
    translated_text TEXT NOT NULL,
    source VARCHAR(50) DEFAULT 'libretranslate',
    confidence DECIMAL(5,4),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() + INTERVAL '24 hours'),
    
    -- Ensure unique translations for text/language pairs
    UNIQUE(original_text, from_lang, to_lang)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_language_configs_language ON language_configs(language);
CREATE INDEX IF NOT EXISTS idx_language_configs_key ON language_configs(key);
CREATE INDEX IF NOT EXISTS idx_language_configs_category ON language_configs(category);
CREATE INDEX IF NOT EXISTS idx_language_configs_module ON language_configs(module);
CREATE INDEX IF NOT EXISTS idx_language_info_code ON language_info(code);
CREATE INDEX IF NOT EXISTS idx_language_info_active ON language_info(is_active);
CREATE INDEX IF NOT EXISTS idx_user_language_preferences_user_id ON user_language_preferences(user_id);
CREATE INDEX IF NOT EXISTS idx_translation_cache_langs ON translation_cache(from_lang, to_lang);
CREATE INDEX IF NOT EXISTS idx_translation_cache_expires ON translation_cache(expires_at);

-- Insert supported languages
INSERT INTO language_info (code, name, native_name, flag, is_rtl, date_format, number_format) VALUES
    ('id', 'Indonesian', 'Bahasa Indonesia', '🇮🇩', false, 'DD/MM/YYYY', '1.234.567,89'),
    ('en', 'English', 'English', '🇺🇸', false, 'MM/DD/YYYY', '1,234,567.89'),
    ('zh', 'Chinese (Mandarin)', '中文', '🇨🇳', false, 'YYYY/MM/DD', '1,234,567.89')
ON CONFLICT (code) DO NOTHING;

-- Insert default system translations
INSERT INTO language_configs (language, key, translation, category, module) VALUES
    -- Indonesian translations
    ('id', 'common.save', 'Simpan', 'common', 'ui'),
    ('id', 'common.cancel', 'Batal', 'common', 'ui'),
    ('id', 'common.submit', 'Kirim', 'common', 'ui'),
    ('id', 'common.loading', 'Memuat...', 'common', 'ui'),
    ('id', 'common.error', 'Terjadi kesalahan', 'common', 'ui'),
    ('id', 'common.success', 'Berhasil', 'common', 'ui'),
    ('id', 'common.confirm', 'Konfirmasi', 'common', 'ui'),
    ('id', 'common.delete', 'Hapus', 'common', 'ui'),
    ('id', 'common.edit', 'Edit', 'common', 'ui'),
    ('id', 'common.view', 'Lihat', 'common', 'ui'),
    ('id', 'common.close', 'Tutup', 'common', 'ui'),
    
    ('id', 'donation.title', 'Donasi', 'donation', 'donation'),
    ('id', 'donation.amount', 'Jumlah', 'donation', 'donation'),
    ('id', 'donation.message', 'Pesan', 'donation', 'donation'),
    ('id', 'donation.submit', 'Kirim Donasi', 'donation', 'donation'),
    ('id', 'donation.success', 'Donasi berhasil dikirim', 'donation', 'donation'),
    ('id', 'donation.failed', 'Donasi gagal', 'donation', 'donation'),
    ('id', 'donation.anonymous', 'Anonim', 'donation', 'donation'),
    ('id', 'donation.display_name', 'Nama Pengirim', 'donation', 'donation'),
    
    ('id', 'payment.processing', 'Memproses pembayaran...', 'payment', 'payment'),
    ('id', 'payment.success', 'Pembayaran berhasil', 'payment', 'payment'),
    ('id', 'payment.failed', 'Pembayaran gagal', 'payment', 'payment'),
    ('id', 'payment.cancelled', 'Pembayaran dibatalkan', 'payment', 'payment'),
    
    ('id', 'auth.login', 'Masuk', 'auth', 'auth'),
    ('id', 'auth.register', 'Daftar', 'auth', 'auth'),
    ('id', 'auth.logout', 'Keluar', 'auth', 'auth'),
    ('id', 'auth.username', 'Nama Pengguna', 'auth', 'auth'),
    ('id', 'auth.email', 'Email', 'auth', 'auth'),
    ('id', 'auth.password', 'Kata Sandi', 'auth', 'auth'),
    ('id', 'auth.confirm_password', 'Konfirmasi Kata Sandi', 'auth', 'auth'),

    -- English translations
    ('en', 'common.save', 'Save', 'common', 'ui'),
    ('en', 'common.cancel', 'Cancel', 'common', 'ui'),
    ('en', 'common.submit', 'Submit', 'common', 'ui'),
    ('en', 'common.loading', 'Loading...', 'common', 'ui'),
    ('en', 'common.error', 'An error occurred', 'common', 'ui'),
    ('en', 'common.success', 'Success', 'common', 'ui'),
    ('en', 'common.confirm', 'Confirm', 'common', 'ui'),
    ('en', 'common.delete', 'Delete', 'common', 'ui'),
    ('en', 'common.edit', 'Edit', 'common', 'ui'),
    ('en', 'common.view', 'View', 'common', 'ui'),
    ('en', 'common.close', 'Close', 'common', 'ui'),
    
    ('en', 'donation.title', 'Donation', 'donation', 'donation'),
    ('en', 'donation.amount', 'Amount', 'donation', 'donation'),
    ('en', 'donation.message', 'Message', 'donation', 'donation'),
    ('en', 'donation.submit', 'Send Donation', 'donation', 'donation'),
    ('en', 'donation.success', 'Donation sent successfully', 'donation', 'donation'),
    ('en', 'donation.failed', 'Donation failed', 'donation', 'donation'),
    ('en', 'donation.anonymous', 'Anonymous', 'donation', 'donation'),
    ('en', 'donation.display_name', 'Display Name', 'donation', 'donation'),
    
    ('en', 'payment.processing', 'Processing payment...', 'payment', 'payment'),
    ('en', 'payment.success', 'Payment successful', 'payment', 'payment'),
    ('en', 'payment.failed', 'Payment failed', 'payment', 'payment'),
    ('en', 'payment.cancelled', 'Payment cancelled', 'payment', 'payment'),
    
    ('en', 'auth.login', 'Login', 'auth', 'auth'),
    ('en', 'auth.register', 'Register', 'auth', 'auth'),
    ('en', 'auth.logout', 'Logout', 'auth', 'auth'),
    ('en', 'auth.username', 'Username', 'auth', 'auth'),
    ('en', 'auth.email', 'Email', 'auth', 'auth'),
    ('en', 'auth.password', 'Password', 'auth', 'auth'),
    ('en', 'auth.confirm_password', 'Confirm Password', 'auth', 'auth'),

    -- Chinese translations
    ('zh', 'common.save', '保存', 'common', 'ui'),
    ('zh', 'common.cancel', '取消', 'common', 'ui'),
    ('zh', 'common.submit', '提交', 'common', 'ui'),
    ('zh', 'common.loading', '加载中...', 'common', 'ui'),
    ('zh', 'common.error', '发生错误', 'common', 'ui'),
    ('zh', 'common.success', '成功', 'common', 'ui'),
    ('zh', 'common.confirm', '确认', 'common', 'ui'),
    ('zh', 'common.delete', '删除', 'common', 'ui'),
    ('zh', 'common.edit', '编辑', 'common', 'ui'),
    ('zh', 'common.view', '查看', 'common', 'ui'),
    ('zh', 'common.close', '关闭', 'common', 'ui'),
    
    ('zh', 'donation.title', '捐赠', 'donation', 'donation'),
    ('zh', 'donation.amount', '金额', 'donation', 'donation'),
    ('zh', 'donation.message', '消息', 'donation', 'donation'),
    ('zh', 'donation.submit', '发送捐赠', 'donation', 'donation'),
    ('zh', 'donation.success', '捐赠发送成功', 'donation', 'donation'),
    ('zh', 'donation.failed', '捐赠失败', 'donation', 'donation'),
    ('zh', 'donation.anonymous', '匿名', 'donation', 'donation'),
    ('zh', 'donation.display_name', '显示名称', 'donation', 'donation'),
    
    ('zh', 'payment.processing', '处理付款中...', 'payment', 'payment'),
    ('zh', 'payment.success', '付款成功', 'payment', 'payment'),
    ('zh', 'payment.failed', '付款失败', 'payment', 'payment'),
    ('zh', 'payment.cancelled', '付款已取消', 'payment', 'payment'),
    
    ('zh', 'auth.login', '登录', 'auth', 'auth'),
    ('zh', 'auth.register', '注册', 'auth', 'auth'),
    ('zh', 'auth.logout', '登出', 'auth', 'auth'),
    ('zh', 'auth.username', '用户名', 'auth', 'auth'),
    ('zh', 'auth.email', '电子邮件', 'auth', 'auth'),
    ('zh', 'auth.password', '密码', 'auth', 'auth'),
    ('zh', 'auth.confirm_password', '确认密码', 'auth', 'auth')
ON CONFLICT (language, key) DO NOTHING;

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers to auto-update updated_at timestamps
CREATE TRIGGER update_language_configs_updated_at 
    BEFORE UPDATE ON language_configs 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_language_info_updated_at 
    BEFORE UPDATE ON language_info 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_language_preferences_updated_at 
    BEFORE UPDATE ON user_language_preferences 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create cleanup job for expired translation cache
CREATE OR REPLACE FUNCTION cleanup_expired_translations()
RETURNS void AS $$
BEGIN
    DELETE FROM translation_cache WHERE expires_at < NOW();
END;
$$ language 'plpgsql';

-- Grant permissions (if needed)
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO language_user;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO language_user; 