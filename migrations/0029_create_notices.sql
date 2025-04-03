CREATE TABLE notices (
    id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL,
    level INTEGER NOT NULL DEFAULT 0,
    fee_amount BIGINT NOT NULL,
    sent_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL DEFAULT '',
    dispatched_at TIMESTAMPTZ,
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
