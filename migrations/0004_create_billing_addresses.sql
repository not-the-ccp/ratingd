CREATE TABLE billing_addresses (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    street TEXT NOT NULL DEFAULT '',
    zip TEXT NOT NULL DEFAULT '',
    city TEXT NOT NULL DEFAULT '',
    country TEXT NOT NULL DEFAULT '',
    primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
