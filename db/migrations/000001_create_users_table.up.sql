CREATE TABLE IF NOT EXISTS users (
    id         uuid         PRIMARY KEY DEFAULT uuidv7(),
    email      varchar(255) NOT NULL,
    name       varchar(100) NOT NULL,
    created_at timestamptz  NOT NULL DEFAULT now(),
    updated_at timestamptz  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_key ON users (lower(email));