-- +goose Up
CREATE TABLE users (
   id BIGINT AUTO_INCREMENT PRIMARY KEY,
   password VARCHAR(255) NOT NULL,
   name VARCHAR(100),
   email VARCHAR(100) NOT NULL,
   address VARCHAR(255),
   phone_number VARCHAR(20),
   date_of_birth DATE,
   gender INT,
   occupation VARCHAR(255),
   verification BOOLEAN DEFAULT FALSE,  -- Keep this line
   status BOOLEAN DEFAULT TRUE,          -- Changed field name to `status` (was `verification` again)
   created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
   updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
--
-- CREATE INDEX idx_users_name ON users(name);
-- CREATE INDEX idx_users_email ON users(email);
-- CREATE INDEX idx_users_name_email ON users(name, email);


-- +goose Down
DROP TABLE users;
