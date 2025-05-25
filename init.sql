-- Initial database setup for donation system
-- This file will be executed when PostgreSQL container starts

-- Enable UUID extension for generating UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create indexes for better performance (will be created after tables are migrated by GORM)
-- These are just examples, GORM will handle the actual table creation

-- Set timezone
SET timezone = 'Asia/Jakarta';

-- Create a function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql'; 