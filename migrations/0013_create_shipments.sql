CREATE TABLE shipments (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    carrier_id BIGINT NOT NULL,
    reference TEXT NOT NULL DEFAULT '',
    origin_zip TEXT NOT NULL DEFAULT '',
    destination_zip TEXT NOT NULL DEFAULT '',
    weight_kg DOUBLE PRECISION NOT NULL DEFAULT 0,
    loading_meters DOUBLE PRECISION NOT NULL DEFAULT 0,
    pallets INTEGER NOT NULL DEFAULT 0,
    pickup_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL DEFAULT '',
    picked_up_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
