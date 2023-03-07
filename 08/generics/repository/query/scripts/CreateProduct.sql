INSERT INTO product (
    id,
    name,
    description,
    price,
    created_at,
    updated_at
) VALUES
    (@id, @name, @description, @price, @created_at, @updated_at)
