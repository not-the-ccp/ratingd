CREATE TABLE fuel_surcharges (
    id BIGSERIAL PRIMARY KEY,
    tariff_id BIGINT NOT NULL,
    valid_from TIMESTAMPTZ NOT NULL,
    percent DOUBLE PRECISION NOT NULL DEFAULT 0,
    diesel_price BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
