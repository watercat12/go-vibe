-- +migrate Up
CREATE TABLE savings_account_details (
    account_id UUID PRIMARY KEY REFERENCES accounts(id) ON DELETE CASCADE,
    is_fixed_term BOOLEAN NOT NULL,
    term_months INTEGER,
    annual_interest_rate DECIMAL(5,4) NOT NULL,
    start_date DATE NOT NULL,
    maturity_date DATE,
    last_interest_calc_date DATE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_savings_account_details_is_fixed_term ON savings_account_details(is_fixed_term);
CREATE INDEX idx_savings_account_details_maturity_date ON savings_account_details(maturity_date);

-- +migrate Down
DROP TABLE savings_account_details;