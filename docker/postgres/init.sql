-- Initialize PostgreSQL database for Trading Alchemist
-- This script runs when the container is first created

-- Create database if it doesn't exist
SELECT 'CREATE DATABASE trading_alchemist_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'trading_alchemist_db')\gexec

-- Connect to the database
\c trading_alchemist_db;

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set timezone
SET timezone = 'UTC'; 