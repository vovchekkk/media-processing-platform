DROP TABLE IF EXISTS tasks, sessions, users CASCADE;

CREATE TABLE users (
    id uuid PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
)

CREATE TABLE tasks (
    id uuid PRIMARY KEY,
    status VARCHAR(255) NOT NULL,
    result TEXT
)

CREATE TABLE sessions (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL,

    ADD CCONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
)