INSERT INTO customer_address (
    id,
    customer_id,
    location_name,
    address,
    created_at,
    updated_at
) VALUES
    (@id, @customer_id, @location_name, @address, @created_at, @updated_at)
