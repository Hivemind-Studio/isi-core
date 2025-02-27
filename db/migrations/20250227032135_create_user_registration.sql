-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_registration (
   id BIGINT AUTO_INCREMENT PRIMARY KEY,
   user_id BIGINT NOT NULL,
   campaign_id BIGINT NOT NULL,
   registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   ip_address VARCHAR(45),
   user_agent VARCHAR(255),
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   FOREIGN KEY (campaign_id) REFERENCES campaign(id),
   INDEX idx_campaign_id (campaign_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
