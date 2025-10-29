-- +migrate Up
CREATE TABLE flexible_savings_interest_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    calculation_date DATE NOT NULL,
    eod_balance DECIMAL(15,2) NOT NULL,
    annual_rate_applied DECIMAL(5,4) NOT NULL,
    daily_interest_amount DECIMAL(15,2) NOT NULL,
    is_promotional_rate BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_flexible_savings_interest_history_account_id ON flexible_savings_interest_history(account_id);
CREATE INDEX idx_flexible_savings_interest_history_calculation_date ON flexible_savings_interest_history(calculation_date);

-- +migrate Down
DROP TABLE flexible_savings_interest_history;