-- Migration: Add Platform Integration Tables
-- Description: Add tables to support YouTube and TikTok integration
-- Date: 2024-01-XX

-- Tabel untuk menyimpan informasi platform streaming
CREATE TABLE IF NOT EXISTS streaming_platforms (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    platform_type VARCHAR(20) NOT NULL CHECK (platform_type IN ('youtube', 'tiktok', 'twitch')),
    platform_user_id VARCHAR(255) NOT NULL,
    platform_username VARCHAR(255) NOT NULL,
    channel_url TEXT NOT NULL,
    channel_name VARCHAR(255),
    profile_image_url TEXT,
    follower_count INTEGER DEFAULT 0,
    is_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, platform_type, platform_user_id)
);

-- Tabel untuk menyimpan konten streaming aktif
CREATE TABLE IF NOT EXISTS streaming_content (
    id SERIAL PRIMARY KEY,
    platform_id INTEGER REFERENCES streaming_platforms(id) ON DELETE CASCADE,
    content_type VARCHAR(20) NOT NULL CHECK (content_type IN ('live', 'video', 'short')),
    content_id VARCHAR(255) NOT NULL,
    content_url TEXT NOT NULL,
    title VARCHAR(500),
    description TEXT,
    thumbnail_url TEXT,
    duration INTEGER, -- dalam detik, NULL untuk live stream
    view_count INTEGER DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    is_live BOOLEAN DEFAULT FALSE,
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(platform_id, content_id)
);

-- Tabel untuk tracking donasi per konten
CREATE TABLE IF NOT EXISTS content_donations (
    id SERIAL PRIMARY KEY,
    donation_id INTEGER REFERENCES donations(id) ON DELETE CASCADE,
    content_id INTEGER REFERENCES streaming_content(id) ON DELETE SET NULL,
    platform_type VARCHAR(20) NOT NULL,
    content_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk performa
CREATE INDEX IF NOT EXISTS idx_streaming_platforms_user_platform ON streaming_platforms(user_id, platform_type);
CREATE INDEX IF NOT EXISTS idx_streaming_platforms_active ON streaming_platforms(is_active, platform_type);
CREATE INDEX IF NOT EXISTS idx_streaming_content_platform_live ON streaming_content(platform_id, is_live);
CREATE INDEX IF NOT EXISTS idx_streaming_content_type ON streaming_content(content_type, is_live);
CREATE INDEX IF NOT EXISTS idx_content_donations_content ON content_donations(content_id);
CREATE INDEX IF NOT EXISTS idx_content_donations_donation ON content_donations(donation_id);
CREATE INDEX IF NOT EXISTS idx_content_donations_platform ON content_donations(platform_type);

-- Trigger untuk update timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply trigger ke tabel yang membutuhkan
CREATE TRIGGER update_streaming_platforms_updated_at 
    BEFORE UPDATE ON streaming_platforms 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_streaming_content_updated_at 
    BEFORE UPDATE ON streaming_content 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert sample data untuk testing
INSERT INTO streaming_platforms (user_id, platform_type, platform_user_id, platform_username, channel_url, channel_name, is_verified) 
VALUES 
    (1, 'youtube', 'UC_sample_channel_1', 'sample_creator_1', 'https://www.youtube.com/@sample_creator_1', 'Sample Gaming Channel', true),
    (2, 'tiktok', 'sample_creator_2', 'sample_creator_2', 'https://www.tiktok.com/@sample_creator_2', 'Sample TikTok Creator', false)
ON CONFLICT (user_id, platform_type, platform_user_id) DO NOTHING;

-- Insert sample content
INSERT INTO streaming_content (platform_id, content_type, content_id, content_url, title, is_live)
VALUES 
    (1, 'live', 'sample_live_123', 'https://www.youtube.com/watch?v=sample_live_123', 'Live Gaming Stream', true),
    (1, 'video', 'sample_video_456', 'https://www.youtube.com/watch?v=sample_video_456', 'Gaming Tutorial', false),
    (2, 'video', '7234567890123456789', 'https://www.tiktok.com/@sample_creator_2/video/7234567890123456789', 'Funny Gaming Moment', false)
ON CONFLICT (platform_id, content_id) DO NOTHING;

-- Verify tables created successfully
SELECT 'streaming_platforms' as table_name, count(*) as row_count FROM streaming_platforms
UNION ALL
SELECT 'streaming_content' as table_name, count(*) as row_count FROM streaming_content
UNION ALL
SELECT 'content_donations' as table_name, count(*) as row_count FROM content_donations;

-- Show table structure
\d streaming_platforms;
\d streaming_content;
\d content_donations; 