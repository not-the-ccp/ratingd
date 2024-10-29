CREATE TABLE surcharges (
    id BIGSERIAL PRIMARY KEY,
    quote_id BIGINT NOT NULL,
    code TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
