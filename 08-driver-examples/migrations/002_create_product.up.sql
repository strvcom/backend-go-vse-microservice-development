CREATE TABLE product (
    id          UUID PRIMARY KEY,
    name        VARCHAR(255)  NOT NULL,
    description VARCHAR(1024) NOT NULL,
    price       NUMERIC       NOT NULL,
    created_at  TIMESTAMP     NOT NULL,
    updated_at  TIMESTAMP     NOT NULL
);
