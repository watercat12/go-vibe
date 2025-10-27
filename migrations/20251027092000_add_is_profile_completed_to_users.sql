-- +migrate Up
ALTER TABLE users ADD COLUMN is_profile_completed BOOLEAN DEFAULT FALSE;

-- +migrate Down
ALTER TABLE users DROP COLUMN is_profile_completed;