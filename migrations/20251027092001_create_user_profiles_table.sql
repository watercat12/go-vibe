-- +migrate Up
CREATE TABLE user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    display_name VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(500),
    phone_number VARCHAR(20) NOT NULL UNIQUE,
    national_id VARCHAR(20) NOT NULL UNIQUE,
    birth_year INTEGER NOT NULL,
    gender VARCHAR(10) NOT NULL,
    team VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_user_profiles_phone_number ON user_profiles(phone_number);
CREATE INDEX idx_user_profiles_national_id ON user_profiles(national_id);

-- +migrate Down
DROP TABLE user_profiles;