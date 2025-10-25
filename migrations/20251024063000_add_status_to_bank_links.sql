-- +migrate Up
ALTER TABLE bank_links ADD COLUMN status VARCHAR(20) DEFAULT 'ACTIVE';

-- +migrate Down
ALTER TABLE bank_links DROP COLUMN status;