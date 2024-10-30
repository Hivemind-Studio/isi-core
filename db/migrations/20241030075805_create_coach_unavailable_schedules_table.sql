-- +goose Up
CREATE TABLE coach_unavailable_schedules (
                                             id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                             coach_id BIGINT,
                                             period_start TIMESTAMP NOT NULL,
                                             period_end TIMESTAMP NOT NULL,
                                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                             created_by BIGINT,
                                             updated_by BIGINT,
                                             FOREIGN KEY (coach_id) REFERENCES coaches(id)
);

-- +goose Down
DROP TABLE coach_unavailable_schedules;
