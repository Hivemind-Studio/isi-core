-- +goose Up
-- +goose StatementBegin
-- Add version column to users table
ALTER TABLE users
ADD COLUMN version INT DEFAULT 0 NOT NULL;

-- Add version column to email_verifications table
ALTER TABLE email_verifications
ADD COLUMN version INT DEFAULT 0 NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
