-- +migrate Up
CREATE TABLE profiles (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    display_name VARCHAR(100),
    avatar_url TEXT,
    phone_number VARCHAR(15),
    national_id VARCHAR(20),
    birth_year INT,
    gender VARCHAR(10),
    team VARCHAR(20),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_profiles_user_id ON profiles(user_id);

-- +migrate Down
DROP TABLE profiles;