-- +goose Up
CREATE TABLE users (
                       id BIGINT AUTO_INCREMENT PRIMARY KEY,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       name VARCHAR(100),
                       email VARCHAR(100) UNIQUE NOT NULL,
                       address TEXT,
                       phone_number VARCHAR(15),
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                       created_by BIGINT,
                       updated_by BIGINT
);

-- +goose Down
DROP TABLE users;
