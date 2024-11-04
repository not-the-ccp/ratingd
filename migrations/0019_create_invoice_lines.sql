CREATE TABLE invoice_lines (
    id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL,
    shipment_id BIGINT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    quantity INTEGER NOT NULL DEFAULT 0,
    unit_price BIGINT NOT NULL,
    amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
