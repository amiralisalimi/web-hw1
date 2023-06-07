BEGIN;

CREATE TABLE USERS (
    name VARCHAR(100),
    family VARCHAR(100),
    id INT PRIMARY KEY,
    age INT,
    sex VARCHAR(6),
    createdAt TIMESTAMP WITH TIME ZONE
);

INSERT INTO USERS (name, family, id, age, sex, createdAt) VALUES
    ('amir', 'nejad', 1, 18, 'male', '2023-06-07 23:53:54.325440209 +0330'),
    ('ali', 'salim', 2, 20, 'male', '2022-06-07 23:53:54.325440209 +0330'),
    ('fati', 'paydar', 3, 22, 'female', '2023-06-06 23:53:52.325440209 +0330');

COMMIT;
