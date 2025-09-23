-- +migrate Up
CREATE TABLE IF NOT EXISTS user_profiles (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255),
    email VARCHAR(255),
    avatar VARCHAR(255),
    phone VARCHAR(20),
    id_number VARCHAR(20),
    birth_year INT,
    gender VARCHAR(10),
    team VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE user_profiles;