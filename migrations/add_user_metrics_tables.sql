-- User Metrics Tables Migration
-- This adds tables required for user metrics collection

-- User activities tracking table
CREATE TABLE IF NOT EXISTS user_activities (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    activity_type VARCHAR(50) NOT NULL,
    metadata JSONB DEFAULT '{}',
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_activities_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- User sessions tracking table
CREATE TABLE IF NOT EXISTS user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    session_token VARCHAR(255) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT true,
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_sessions_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for performance optimization
CREATE INDEX IF NOT EXISTS idx_user_activities_user_id_created_at ON user_activities(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_activities_created_at ON user_activities(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_activities_activity_type ON user_activities(activity_type);

CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_active ON user_sessions(is_active, last_activity DESC);
CREATE INDEX IF NOT EXISTS idx_user_sessions_token ON user_sessions(session_token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at ON user_sessions(expires_at);

-- Add last_login_at and last_activity_at to users table if not exists
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='last_login_at') THEN
        ALTER TABLE users ADD COLUMN last_login_at TIMESTAMP;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='last_activity_at') THEN
        ALTER TABLE users ADD COLUMN last_activity_at TIMESTAMP;
    END IF;
END $$;

-- Create function to update last_activity_at automatically
CREATE OR REPLACE FUNCTION update_user_last_activity()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE users 
    SET last_activity_at = NEW.created_at 
    WHERE id = NEW.user_id;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to auto-update user last activity
DROP TRIGGER IF EXISTS trigger_update_user_last_activity ON user_activities;
CREATE TRIGGER trigger_update_user_last_activity
    AFTER INSERT ON user_activities
    FOR EACH ROW EXECUTE FUNCTION update_user_last_activity();

-- Create view for active users statistics
CREATE OR REPLACE VIEW user_metrics_summary AS
SELECT 
    (SELECT COUNT(*) FROM users WHERE deleted_at IS NULL) as total_users,
    (SELECT COUNT(DISTINCT user_id) FROM user_activities WHERE created_at >= NOW() - INTERVAL '24 hours') as active_users_24h,
    (SELECT COUNT(DISTINCT user_id) FROM user_activities WHERE created_at >= NOW() - INTERVAL '7 days') as active_users_7d,
    (SELECT COUNT(DISTINCT user_id) FROM user_activities WHERE created_at >= NOW() - INTERVAL '30 days') as active_users_30d,
    (SELECT COUNT(DISTINCT user_id) FROM user_sessions WHERE is_active = true AND last_activity >= NOW() - INTERVAL '1 hour') as active_users_current,
    (SELECT COUNT(DISTINCT user_id) FROM user_sessions WHERE is_active = true AND last_activity >= NOW() - INTERVAL '5 minutes') as online_users,
    (SELECT COUNT(*) FROM user_sessions WHERE is_active = true AND expires_at > NOW()) as active_sessions;

-- Insert sample activity types
INSERT INTO user_activities (user_id, activity_type) 
SELECT 1, 'login' WHERE EXISTS (SELECT 1 FROM users WHERE id = 1)
ON CONFLICT DO NOTHING;

-- Comments for documentation
COMMENT ON TABLE user_activities IS 'Tracks all user activities for metrics collection';
COMMENT ON TABLE user_sessions IS 'Tracks user sessions for active/online user metrics';
COMMENT ON VIEW user_metrics_summary IS 'Provides aggregated user metrics for monitoring';

-- Grant permissions (adjust as needed)
-- GRANT SELECT ON user_activities TO monitoring_user;
-- GRANT SELECT ON user_sessions TO monitoring_user;
-- GRANT SELECT ON user_metrics_summary TO monitoring_user; 