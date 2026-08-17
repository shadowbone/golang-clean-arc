CREATE TABLE IF NOT EXISTS products (
    id         uuid         PRIMARY KEY DEFAULT uuidv7(),
    name       varchar(150) NOT NULL,
    price      bigint       NOT NULL CHECK (price >= 0),
    created_at timestamptz  NOT NULL DEFAULT now(),
    updated_at timestamptz  NOT NULL DEFAULT now()
);