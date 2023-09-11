CREATE TABLE audit_log (
    id BIGSERIAL PRIMARY KEY,
    action TEXT NOT NULL,
    entity_id BIGINT NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL
);
