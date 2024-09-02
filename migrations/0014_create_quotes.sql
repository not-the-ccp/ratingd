CREATE TABLE quotes (
    id BIGSERIAL PRIMARY KEY,
    shipment_id BIGINT NOT NULL,
    contract_id BIGINT NOT NULL,
    net_amount BIGINT NOT NULL,
    fuel_amount BIGINT NOT NULL,
    currency TEXT NOT NULL DEFAULT '',
    expires_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL DEFAULT '',
    accepted_at TIMESTAMPTZ,
    superseded_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
