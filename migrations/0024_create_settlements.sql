CREATE TABLE settlements (
    id BIGSERIAL PRIMARY KEY,
    carrier_id BIGINT NOT NULL,
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    gross_amount BIGINT NOT NULL,
    deductions BIGINT NOT NULL,
    status TEXT NOT NULL DEFAULT '',
    reviewed_at TIMESTAMPTZ,
    released_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
