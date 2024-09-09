CREATE TABLE invoices (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    number TEXT NOT NULL DEFAULT '',
    currency TEXT NOT NULL DEFAULT '',
    issued_at TIMESTAMPTZ NOT NULL,
    due_at TIMESTAMPTZ NOT NULL,
    net_amount BIGINT NOT NULL,
    tax_amount BIGINT NOT NULL,
    status TEXT NOT NULL DEFAULT '',
    issued_on TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    voided_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
