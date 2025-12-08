package export

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// SQLStore implements Store on top of PostgreSQL.
type SQLStore struct {
	db *sql.DB
}

func (s *SQLStore) FindBatch(ctx context.Context, id int64) (*Batch, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, row_count, failed_at, finished_at, status, file_name, kind, started_at, created_at, updated_at FROM batches WHERE id = $1`, id)
	var batch Batch
	if err := row.Scan(&batch.ID, &batch.RowCount, &batch.FailedAt, &batch.FinishedAt, &batch.Status, &batch.FileName, &batch.Kind, &batch.StartedAt, &batch.CreatedAt, &batch.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query batch %d: %w", id, err)
	}
	return &batch, nil
}

func (s *SQLStore) ListBatches(ctx context.Context, limit int, offset int) ([]*Batch, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, row_count, failed_at, finished_at, status, file_name, kind, started_at, created_at, updated_at FROM batches ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list batches: %w", err)
	}
	defer rows.Close()

	var batches []*Batch
	for rows.Next() {
		var batch Batch
		if err := rows.Scan(&batch.ID, &batch.RowCount, &batch.FailedAt, &batch.FinishedAt, &batch.Status, &batch.FileName, &batch.Kind, &batch.StartedAt, &batch.CreatedAt, &batch.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan batch: %w", err)
		}
		batches = append(batches, &batch)
	}
	return batches, rows.Err()
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) ListBatchesByStatus(ctx context.Context, status BatchStatus) ([]*Batch, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, row_count, failed_at, finished_at, status, file_name, kind, started_at, created_at, updated_at FROM batches WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("list batches: %w", err)
	}
	defer rows.Close()

	var batches []*Batch
	for rows.Next() {
		var batch Batch
		if err := rows.Scan(&batch.ID, &batch.RowCount, &batch.FailedAt, &batch.FinishedAt, &batch.Status, &batch.FileName, &batch.Kind, &batch.StartedAt, &batch.CreatedAt, &batch.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan batch: %w", err)
		}
		batches = append(batches, &batch)
	}
	return batches, rows.Err()
}

func (s *SQLStore) UpdateBatch(ctx context.Context, batch *Batch) error {
	res, err := s.db.ExecContext(ctx, `UPDATE batches SET row_count = $2, failed_at = $3, finished_at = $4, status = $5, file_name = $6, kind = $7, started_at = $8, updated_at = $9 WHERE id = $1`, batch.ID, batch.RowCount, batch.FailedAt, batch.FinishedAt, batch.Status, batch.FileName, batch.Kind, batch.StartedAt, batch.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update batch: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLStore) InsertBatch(ctx context.Context, batch *Batch) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO batches (row_count, failed_at, finished_at, status, file_name, kind, started_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`, batch.RowCount, batch.FailedAt, batch.FinishedAt, batch.Status, batch.FileName, batch.Kind, batch.StartedAt, batch.CreatedAt, batch.UpdatedAt).Scan(&batch.ID); err != nil {
		return err
	}

	return nil
}
