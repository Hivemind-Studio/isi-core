-- +goose Up
CREATE TABLE bookings (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    coach_id BIGINT,
    user_id BIGINT,
    period_start DATETIME NOT NULL,
    period_end DATETIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (coach_id) REFERENCES coaches(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE bookings;
