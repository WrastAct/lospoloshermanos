CREATE TABLE IF NOT EXISTS insurance (
    insurance_id bigserial PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    region VARCHAR(255) NOT NULL
);