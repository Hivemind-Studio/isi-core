-- +goose Up
-- +goose StatementBegin
CREATE TABLE campaign (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(200) NOT NULL,
  channel VARCHAR(50) NOT NULL,
  start_date DATETIME NOT NULL,
  end_date DATETIME NOT NULL,
  link VARCHAR(255) NOT NULL,
  campaign_id VARCHAR(255) NOT NULL,
  status BOOLEAN DEFAULT TRUE,
  version INT DEFAULT 0 NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY idx_campaign_id (campaign_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE campaign;
-- +goose StatementEnd