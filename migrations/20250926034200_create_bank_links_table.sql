-- +migrate Up
CREATE TABLE bank_links (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bank_code VARCHAR(10) NOT NULL,
    account_type VARCHAR(20) NOT NULL,
    access_token VARCHAR(255),
    refresh_token VARCHAR(255),
    expires_in INT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_bank_links_user_id ON bank_links(user_id);
CREATE INDEX idx_bank_links_bank_code ON bank_links(bank_code);

-- +migrate Down
DROP TABLE bank_links;