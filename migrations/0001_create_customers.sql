CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL DEFAULT '',
    number TEXT NOT NULL DEFAULT '',
    vat_number TEXT NOT NULL DEFAULT '',
    country TEXT NOT NULL DEFAULT '',
    payment_term_days INTEGER NOT NULL DEFAULT 0,
    credit_limit BIGINT NOT NULL,
    status TEXT NOT NULL DEFAULT '',
    activated_at TIMESTAMPTZ,
    blocked_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
