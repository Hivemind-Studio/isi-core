-- +goose Up
CREATE TABLE role_accesibilities (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    role_id BIGINT,
    page VARCHAR(100) NOT NULL,
    path VARCHAR(255) NOT NULL,
    `group` VARCHAR(100),
    accessibility VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

-- +goose Down
DROP TABLE role_accessibilities;
