INSERT INTO customer_address (
    customer_id,
    location_name,
    address,
    created_at,
    updated_at
) VALUES
    (@customer_id, @location_name, @address, @created_at, @updated_at)
