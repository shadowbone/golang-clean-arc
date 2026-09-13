CREATE TABLE IF NOT EXISTS orders (
    id         uuid        PRIMARY KEY DEFAULT uuidv7(),
    user_id    uuid        NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    total      bigint      NOT NULL DEFAULT 0 CHECK (total >= 0),
    status     varchar(20) NOT NULL DEFAULT 'pending',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS orders_user_id_idx ON orders (user_id, id DESC);

CREATE TABLE IF NOT EXISTS order_items (
    id         uuid   PRIMARY KEY DEFAULT uuidv7(),
    order_id   uuid   NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id uuid   NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    nama       varchar(150) NOT NULL,
    qty        int    NOT NULL CHECK (qty > 0),
    price      bigint NOT NULL CHECK (price >= 0),
    subtotal   bigint NOT NULL CHECK (subtotal >= 0)
);

CREATE INDEX IF NOT EXISTS order_items_order_id_idx ON order_items (order_id);