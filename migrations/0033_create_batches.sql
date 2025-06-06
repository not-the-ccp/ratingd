CREATE TABLE batches (
    id BIGSERIAL PRIMARY KEY,
    kind TEXT NOT NULL DEFAULT '',
    file_name TEXT NOT NULL DEFAULT '',
    row_count INTEGER NOT NULL DEFAULT 0,
    started_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL DEFAULT '',
    finished_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
