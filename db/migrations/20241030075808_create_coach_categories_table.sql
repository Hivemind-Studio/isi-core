-- +goose Up
CREATE TABLE coach_categories (
                                  coach_id BIGINT,
                                  category_id BIGINT,
                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                  created_by BIGINT,
                                  updated_by BIGINT,
                                  PRIMARY KEY (coach_id, category_id),
                                  FOREIGN KEY (coach_id) REFERENCES coaches(id),
                                  FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- +goose Down
DROP TABLE coach_categories;
