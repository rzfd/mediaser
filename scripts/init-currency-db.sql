-- Currency Database Initialization Script
-- This script sets up the schema for currency exchange rates and user preferences

-- Create database and user if they don't exist
-- (This is handled by docker-entrypoint-initdb.d)

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create currency_rates table
CREATE TABLE IF NOT EXISTS currency_rates (
    id SERIAL PRIMARY KEY,
    from_currency VARCHAR(10) NOT NULL,
    to_currency VARCHAR(10) NOT NULL,
    rate DECIMAL(20,8) NOT NULL,
    source VARCHAR(50) NOT NULL DEFAULT 'exchangerate-api',
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_updated BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure unique currency pairs per source
    UNIQUE(from_currency, to_currency, source)
);

-- Create currency_info table
CREATE TABLE IF NOT EXISTS currency_info (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    decimal_unit INTEGER NOT NULL DEFAULT 2,
    is_active BOOLEAN NOT NULL DEFAULT true,
    country VARCHAR(100),
    region VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create user_currency_preferences table
CREATE TABLE IF NOT EXISTS user_currency_preferences (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    primary_currency VARCHAR(10) NOT NULL DEFAULT 'IDR',
    secondary_currency VARCHAR(10) NOT NULL DEFAULT 'USD',
    auto_convert BOOLEAN NOT NULL DEFAULT false,
    show_both_currencies BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure one preference per user
    UNIQUE(user_id)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_currency_rates_from_to ON currency_rates(from_currency, to_currency);
CREATE INDEX IF NOT EXISTS idx_currency_rates_active ON currency_rates(is_active);
CREATE INDEX IF NOT EXISTS idx_currency_rates_last_updated ON currency_rates(last_updated);
CREATE INDEX IF NOT EXISTS idx_currency_info_code ON currency_info(code);
CREATE INDEX IF NOT EXISTS idx_currency_info_active ON currency_info(is_active);
CREATE INDEX IF NOT EXISTS idx_user_currency_preferences_user_id ON user_currency_preferences(user_id);

-- Insert supported currencies
INSERT INTO currency_info (code, name, symbol, decimal_unit, country, region) VALUES
    ('IDR', 'Indonesian Rupiah', 'Rp', 0, 'Indonesia', 'Southeast Asia'),
    ('USD', 'US Dollar', '$', 2, 'United States', 'North America'),
    ('CNY', 'Chinese Yuan', '¥', 2, 'China', 'East Asia'),
    ('EUR', 'Euro', '€', 2, 'European Union', 'Europe'),
    ('JPY', 'Japanese Yen', '¥', 0, 'Japan', 'East Asia'),
    ('SGD', 'Singapore Dollar', 'S$', 2, 'Singapore', 'Southeast Asia'),
    ('MYR', 'Malaysian Ringgit', 'RM', 2, 'Malaysia', 'Southeast Asia')
ON CONFLICT (code) DO NOTHING;

-- Insert some initial exchange rates (will be updated by the service)
INSERT INTO currency_rates (from_currency, to_currency, rate, source, last_updated) VALUES
    ('USD', 'IDR', 15000.00, 'initial', EXTRACT(EPOCH FROM NOW())),
    ('USD', 'CNY', 7.20, 'initial', EXTRACT(EPOCH FROM NOW())),
    ('USD', 'EUR', 0.85, 'initial', EXTRACT(EPOCH FROM NOW())),
    ('USD', 'JPY', 110.00, 'initial', EXTRACT(EPOCH FROM NOW())),
    ('USD', 'SGD', 1.35, 'initial', EXTRACT(EPOCH FROM NOW())),
    ('USD', 'MYR', 4.50, 'initial', EXTRACT(EPOCH FROM NOW()))
ON CONFLICT (from_currency, to_currency, source) DO NOTHING;

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers to auto-update updated_at timestamps
CREATE TRIGGER update_currency_rates_updated_at 
    BEFORE UPDATE ON currency_rates 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_currency_info_updated_at 
    BEFORE UPDATE ON currency_info 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_currency_preferences_updated_at 
    BEFORE UPDATE ON user_currency_preferences 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Grant permissions (if needed)
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO currency_user;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO currency_user; 