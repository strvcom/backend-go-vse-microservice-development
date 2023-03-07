BEGIN;

CREATE TABLE customer (
    id            UUID         PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    created_at    TIMESTAMP    NOT NULL,
    updated_at    TIMESTAMP    NOT NULL
);

CREATE TABLE customer_address (
    customer_id   UUID PRIMARY KEY REFERENCES customer(id),
    location_name VARCHAR(255) NOT NULL,
    address       VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP    NOT NULL,
    updated_at    TIMESTAMP    NOT NULL
);

COMMIT;
