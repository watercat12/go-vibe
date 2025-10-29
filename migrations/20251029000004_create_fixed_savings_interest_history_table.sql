-- +migrate Up
CREATE TABLE fixed_savings_interest_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    calculation_period VARCHAR(50) NOT NULL,
    total_interest_amount DECIMAL(15,2) NOT NULL,
    is_early_withdrawal BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_fixed_savings_interest_history_account_id ON fixed_savings_interest_history(account_id);

-- +migrate Down
DROP TABLE fixed_savings_interest_history;