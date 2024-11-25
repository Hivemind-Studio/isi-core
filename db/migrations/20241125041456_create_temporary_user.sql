-- +goose up

CREATE TABLE temporary_user (
   id BIGINT AUTO_INCREMENT PRIMARY KEY,
   email VARCHAR(200) NOT NULL UNIQUE,
   created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE temporary_user;