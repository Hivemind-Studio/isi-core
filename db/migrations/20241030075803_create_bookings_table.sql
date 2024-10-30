-- +goose Up
CREATE TABLE bookings (
                          id BIGINT AUTO_INCREMENT PRIMARY KEY,
                          coach_id BIGINT,
                          user_id BIGINT,
                          period_start TIMESTAMP NOT NULL,
                          period_end TIMESTAMP NOT NULL,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                          created_by BIGINT,
                          updated_by BIGINT,
                          FOREIGN KEY (coach_id) REFERENCES coaches(id),
                          FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE bookings;
