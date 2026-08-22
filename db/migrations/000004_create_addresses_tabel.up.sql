CREATE TABLE IF NOT EXISTS addresses (
    id         uuid         PRIMARY KEY DEFAULT uuidv7(),
    user_id    uuid         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    label      varchar(30)  NOT NULL,
    street     text         NOT NULL,
    city       varchar(80)  NOT NULL,
    is_default boolean      NOT NULL DEFAULT false,
    created_at timestamptz  NOT NULL DEFAULT now(),
    updated_at timestamptz  NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS addresses_user_idx ON addresses (user_id);

CREATE UNIQUE INDEX IF NOT EXISTS addresses_one_default_per_user
    ON addresses (user_id) WHERE is_default;