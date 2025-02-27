-- +goose Up
-- +goose StatementBegin
CREATE TABLE campaign (
     id BIGINT AUTO_INCREMENT PRIMARY KEY,
     name varchar(200) NOT NULL,
     channel varchar(50) NOT NULL,
     start_date DATETIME NOT NULL,
     end_date DATETIME NOT NULL,
     link varchar(255) NOT NULL,
     campaign_id VARCHAR(250) NOT NULL,
     status BOOLEAN DEFAULT TRUE,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
