CREATE TABLE cart (
    customer_id UUID NOT NULL REFERENCES customer(id),
    product_id  UUID NOT NULL REFERENCES product(id),
    inserted_at TIMESTAMP NOT NULL,
    PRIMARY KEY (customer_id, product_id)
);
