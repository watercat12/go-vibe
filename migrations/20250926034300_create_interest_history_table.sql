-- +migrate Up
CREATE TABLE interest_history (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    interest_amount NUMERIC(18,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_interest_history_account_id ON interest_history(account_id);
CREATE INDEX idx_interest_history_date ON interest_history(date);


-- +migrate Down
DROP TABLE interest_history;