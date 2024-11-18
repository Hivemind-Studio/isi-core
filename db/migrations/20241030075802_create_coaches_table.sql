-- +goose Up
CREATE TABLE coaches (
     id BIGINT AUTO_INCREMENT PRIMARY KEY,
     user_id BIGINT,
     certifications TEXT,
     experiences TEXT,
     education TEXT,
     level INT,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE coaches;
