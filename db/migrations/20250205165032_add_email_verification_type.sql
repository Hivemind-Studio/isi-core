-- +goose Up
-- +goose StatementBegin
ALTER TABLE isi.email_verifications ADD `type` varchar(50) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
