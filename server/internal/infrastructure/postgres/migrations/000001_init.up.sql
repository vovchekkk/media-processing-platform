-- Active: 1786720779745@@127.0.0.1@5433@postgres@public
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    filter JSONB,
    image TEXT,
    status VARCHAR(255) NOT NULL DEFAULT 'in_progress',
    result TEXT,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);