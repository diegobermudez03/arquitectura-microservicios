﻿-- Database schema reference for the cards microservice
-- Note: This file is for reference only. GORM will automatically handle migrations.
-- The actual schema is created automatically when the application starts.

-- Users table (auto-created by GORM)
-- CREATE TABLE users (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     user_token TEXT UNIQUE NOT NULL,
--     name TEXT NOT NULL,
--     lastname TEXT NOT NULL,
--     birth_date DATE NOT NULL,
--     country_code TEXT NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     deleted_at TIMESTAMP NULL
-- );

-- Issued cards table (auto-created by GORM)
-- CREATE TABLE issued_cards (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     user_id UUID NOT NULL REFERENCES users(id),
--     user_token TEXT NOT NULL,
--     pan TEXT NOT NULL,
--     cvv TEXT NOT NULL,
--     expiry_date DATE NOT NULL,
--     card_type TEXT NOT NULL,
--     status TEXT NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     deleted_at TIMESTAMP NULL
-- );

-- Failed attempts table (auto-created by GORM)
-- CREATE TABLE failed_attempts (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     user_id UUID NOT NULL REFERENCES users(id),
--     user_token TEXT NOT NULL,
--     card_type TEXT NOT NULL,
--     decline_reason TEXT NOT NULL,
--     status TEXT NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     deleted_at TIMESTAMP NULL
-- );

-- Indexes (auto-created by GORM)
-- CREATE INDEX idx_users_user_token ON users(user_token);
-- CREATE INDEX idx_users_deleted_at ON users(deleted_at);
-- CREATE INDEX idx_issued_cards_user_id ON issued_cards(user_id);
-- CREATE INDEX idx_issued_cards_deleted_at ON issued_cards(deleted_at);
-- CREATE INDEX idx_failed_attempts_user_id ON failed_attempts(user_id);
-- CREATE INDEX idx_failed_attempts_deleted_at ON failed_attempts(deleted_at);

-- To manually create the database (if needed):
-- CREATE DATABASE cards_db;
-- \c cards_db;
-- The application will automatically create all tables and indexes on startup.
