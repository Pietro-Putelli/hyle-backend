CREATE TABLE sessions (
    id SERIAL PRIMARY KEY NOT NULL,
    guid UUID DEFAULT uuid_generate_v4 () NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    device_id VARCHAR(255) NOT NULL,
    device_token VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- create indexes
CREATE INDEX sessions_user_id_index ON sessions (user_id, guid);