-- name: ListCartItems :many
SELECT
    p.id,
    p.name,
    p.description,
    p.price,
    c.inserted_at,
    p.created_at,
    p.updated_at
FROM
    product AS p
JOIN
    cart AS c on c.product_id = p.id
WHERE
    c.customer_id = $1;
