CREATE TABLE credit_notes (
    id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL,
    number TEXT NOT NULL DEFAULT '',
    reason TEXT NOT NULL DEFAULT '',
    amount BIGINT NOT NULL,
    status TEXT NOT NULL DEFAULT '',
    applied_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
