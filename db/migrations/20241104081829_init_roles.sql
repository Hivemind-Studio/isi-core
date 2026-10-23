-- +goose Up
INSERT INTO roles
( name )
VALUES
    ('Admin'),
    ('Staff'),
    ('Coach'),
    ('Coachee'),
    ('Marketing');

-- +goose Down
DELETE FROM roles where 1=1;
