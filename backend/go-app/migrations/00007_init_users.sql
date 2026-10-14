-- +goose Up

CREATE TYPE user_role AS ENUM ('superuser', 'businessman');

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    nickname TEXT NOT NULL UNIQUE,
    u_role user_role NOT NULL,
    pass_hash TEXT NOT NULL
);

-- +goose Down

DROP TABLE users;
DROP TYPE user_role;
