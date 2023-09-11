package audit

import (
	"context"
	"database/sql"
	"fmt"
)

// Recorder persists an audit trail entry for a state change.
type Recorder interface {
	Record(ctx context.Context, action string, entityID int64) error
}

type SQLRecorder struct {
	db *sql.DB
}

func NewSQLRecorder(db *sql.DB) *SQLRecorder {
	return &SQLRecorder{db: db}
}

func (r *SQLRecorder) Record(ctx context.Context, action string, entityID int64) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO audit_log (action, entity_id, recorded_at) VALUES ($1, $2, now())`, action, entityID)
	if err != nil {
		return fmt.Errorf("record audit %s: %w", action, err)
	}
	return nil
}
