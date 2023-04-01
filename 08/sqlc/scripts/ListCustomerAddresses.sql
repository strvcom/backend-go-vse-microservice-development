-- name: ListCustomerAddresses :many
SELECT
    ca.id,
    ca.customer_id,
    ca.location_name,
    ca.address,
    ca.created_at,
    ca.updated_at
FROM
    customer_address AS ca
JOIN
    customer AS c on c.id = ca.customer_id
WHERE
    c.id = $1;
