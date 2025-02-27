-- +goose Up
INSERT INTO users (id, name, email, password, role_id, phone_number, date_of_birth, gender, occupation)
VALUES
    (1, 'Admin', 'admin@isi.com', '21f78900fc72b1929a2b12cac3d184c95fad3733621ae370dc1869cc24291057', 1, '+23456789012', null, null, null),
    (2, 'Staff', 'staff@isi.com', '21f78900fc72b1929a2b12cac3d184c95fad3733621ae370dc1869cc24291057', 2, '+23456789013', null, null, null),
    (3, 'Coach', 'coach@isi.com', '21f78900fc72b1929a2b12cac3d184c95fad3733621ae370dc1869cc24291057', 3, '+23456789014', null, null, null),
    (4, 'Coachee', 'coachee@isi.com', '21f78900fc72b1929a2b12cac3d184c95fad3733621ae370dc1869cc24291057', 4, '+23456789015', null, null, null),
    (5, 'Marketing', 'marketing@isi.com', '21f78900fc72b1929a2b12cac3d184c95fad3733621ae370dc1869cc24291057', 5, '+23456789045', null, null, null);

INSERT INTO coaches (user_id) VALUES (3);
