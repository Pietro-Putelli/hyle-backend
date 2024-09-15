CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY NOT NULL,
    guid UUID DEFAULT uuid_generate_v4 () NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    external_id VARCHAR(255) NOT NULL UNIQUE,
    given_name VARCHAR(255) NOT NULL,
    family_name VARCHAR(255) NOT NULL,
    provider VARCHAR(255) NOT NULL,
    settings JSONB,
    subscription_receipt_id VARCHAR(255) NULL DEFAULT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    is_notification_enabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- create indexes
CREATE INDEX users_guid_index ON users (guid);