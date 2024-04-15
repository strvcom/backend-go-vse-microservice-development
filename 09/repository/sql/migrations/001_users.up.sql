CREATE TABLE users (
    id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at  TIMESTAMPTZ  NOT NULL,
    updated_at  TIMESTAMPTZ  NOT NULL
);

INSERT INTO users (
    id,
    created_at,
    updated_at
) VALUES (
    '00000000-0000-0000-0000-000000000000',
    NOW(),
    NOW()
);
