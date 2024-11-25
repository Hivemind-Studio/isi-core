-- +goose up

CREATE TABLE user_verified_account (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT,
    verification_token varchar(255),
    expires_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE user_verified_account;