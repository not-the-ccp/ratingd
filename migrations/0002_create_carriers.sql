CREATE TABLE carriers (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL DEFAULT '',
    code TEXT NOT NULL DEFAULT '',
    country TEXT NOT NULL DEFAULT '',
    vat_number TEXT NOT NULL DEFAULT '',
    payment_term_days INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT '',
    approved_at TIMESTAMPTZ,
    suspended_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
