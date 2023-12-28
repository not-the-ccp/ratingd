CREATE TABLE contracts (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    tariff_id BIGINT NOT NULL,
    reference TEXT NOT NULL DEFAULT '',
    valid_from TIMESTAMPTZ NOT NULL,
    valid_to TIMESTAMPTZ NOT NULL,
    discount_percent DOUBLE PRECISION NOT NULL DEFAULT 0,
    minimum_charge BIGINT NOT NULL,
    status TEXT NOT NULL DEFAULT '',
    signed_at TIMESTAMPTZ,
    terminated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
