-- +goose up
CREATE TABLE email_verifications (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    verification_token VARCHAR(255) UNIQUE,
    trial TINYINT,
    version INT DEFAULT 0 NOT NULL,
    expired_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_email_verification_token_email (verification_token, email),
    INDEX idx_email_date (email, created_at),
    INDEX idx_expired_at (expired_at)
);

-- +goose Down
DROP TABLE email_verifications;