-- +goose Up
CREATE TABLE users (
   id BIGINT AUTO_INCREMENT PRIMARY KEY,
   password VARCHAR(255) NOT NULL,
   name VARCHAR(100),
   email VARCHAR(100) NOT NULL UNIQUE,
   address VARCHAR(255),
   phone_number VARCHAR(20) UNIQUE,
   date_of_birth DATE,
   gender varchar(6),
   occupation VARCHAR(255),
   role_id BIGINT,
   photo VARCHAR(200),
   verification BOOLEAN DEFAULT FALSE,
   version INT DEFAULT 0 NOT NULL,
   status BOOLEAN DEFAULT TRUE,
   created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
   updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   FOREIGN KEY (role_id) REFERENCES roles(id)
);
--
-- CREATE INDEX idx_users_name ON users(name);
-- CREATE INDEX idx_users_email ON users(email);
-- CREATE INDEX idx_users_name_email ON users(name, email);


-- +goose Down
DROP TABLE users;
