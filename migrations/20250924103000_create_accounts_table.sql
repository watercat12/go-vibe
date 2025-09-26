-- +migrate Up
-- 20240610120000_accounts.sql
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_type VARCHAR(20) NOT NULL,
    account_number VARCHAR(30) UNIQUE NOT NULL,
    account_name VARCHAR(100),
    balance NUMERIC(18,2) DEFAULT 0.00,
    interest_rate NUMERIC(5,2),
    fixed_term_months INT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_accounts_type ON accounts(account_type);


-- +migrate Down
DROP TABLE accounts;