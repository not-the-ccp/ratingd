CREATE TABLE tracking_events (
    id BIGSERIAL PRIMARY KEY,
    shipment_id BIGINT NOT NULL,
    code TEXT NOT NULL DEFAULT '',
    location TEXT NOT NULL DEFAULT '',
    occurred_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
