-- +migrate Up
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    transaction_type VARCHAR(50) NOT NULL CHECK (transaction_type IN ('PAYMENT_INITIATION', 'INTEREST_CREDIT', 'WITHDRAWAL', 'WITHDRAWAL_PENALTY')),
    amount DECIMAL(15,2) NOT NULL,
    transaction_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    description VARCHAR(255),
    is_penalty BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_transactions_account_id ON transactions(account_id);
CREATE INDEX idx_transactions_transaction_type ON transactions(transaction_type);
CREATE INDEX idx_transactions_transaction_date ON transactions(transaction_date);

-- +migrate Down
DROP TABLE transactions;