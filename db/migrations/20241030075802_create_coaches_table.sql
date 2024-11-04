-- +goose Up
CREATE TABLE coaches (
     id BIGINT AUTO_INCREMENT PRIMARY KEY,
     user_id BIGINT,
     certifications TEXT,
     experiences TEXT,
     education TEXT,
     rate DECIMAL(10, 2),
     phone_number VARCHAR(15),
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE coaches;
