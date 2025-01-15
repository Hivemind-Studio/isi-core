-- +goose Up
INSERT INTO users (id, name, email, password, role_id, phone_number, date_of_birth, gender, occupation)
VALUES
    (1, 'Admin', 'admin@isi.com', 'd210de044bf170b9d87eadcb4de927b66aa57302e95e585740c67dbf61a47ae6', 1, '+23456789012', null, null, null),
    (2, 'Staff', 'staff@isi.com', 'd210de044bf170b9d87eadcb4de927b66aa57302e95e585740c67dbf61a47ae6', 2, '+23456789013', null, null, null),
    (3, 'Coach', 'coach@isi.com', 'd210de044bf170b9d87eadcb4de927b66aa57302e95e585740c67dbf61a47ae6', 3, '+23456789014', null, null, null),
    (4, 'Coachee', 'coachee@isi.com', 'd210de044bf170b9d87eadcb4de927b66aa57302e95e585740c67dbf61a47ae6', 4, '+23456789015', null, null, null);

INSERT INTO coaches (user_id) VALUES (4);
