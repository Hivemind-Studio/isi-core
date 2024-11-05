-- +goose Up
CREATE TABLE role_accessibilities (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    role_id BIGINT,
    page VARCHAR(100) NOT NULL,
    accessibility VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

-- +goose Down
DROP TABLE role_accessibilities;
