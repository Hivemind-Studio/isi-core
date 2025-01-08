-- +goose Up
INSERT INTO roles
( name )
VALUES
    ('Admin'),
    ('Staff'),
    ('Coach'),
    ('Coachee');

-- +goose Down
DELETE FROM roles where 1=1;
