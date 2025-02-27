-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN google_id VARCHAR(255) UNIQUE after email;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
