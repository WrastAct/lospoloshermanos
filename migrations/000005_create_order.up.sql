CREATE TABLE IF NOT EXISTS "order" (
    order_id            bigserial PRIMARY KEY,
    sender_id           bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    receiver_id         bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    mass                DECIMAL(8, 2) NOT NULL,
    insurance_id        bigint NOT NULL REFERENCES insurance ON DELETE CASCADE,
    value               DECIMAL(8, 2) NOT NULL,
    insurance_coverage  DECIMAL(8, 2) NOT NULL,
    sender_address      TEXT NOT NULL,
    receiver_address    TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS properties (
    properties_id INTEGER PRIMARY KEY,
    name          VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS order_properties (
    order_id        bigint NOT NULL REFERENCES "order" ON DELETE CASCADE,
    properties_id   INTEGER NOT NULL REFERENCES properties ON DELETE CASCADE,
    PRIMARY KEY (order_id, properties_id)
);

