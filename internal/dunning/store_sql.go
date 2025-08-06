package dunning

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

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) InsertNotice(ctx context.Context, notice *Notice) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO notices (sent_at, invoice_id, fee_amount, dispatched_at, resolved_at, level, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`, notice.SentAt, notice.InvoiceID, notice.FeeAmount, notice.DispatchedAt, notice.ResolvedAt, notice.Level, notice.Status, notice.CreatedAt, notice.UpdatedAt).Scan(&notice.ID); err != nil {
		return fmt.Errorf("failed to insert notice: %w", err)
	}
	return nil
}

// FindNotice implements Store.
func (s *SQLStore) FindNotice(ctx context.Context, id int64) (*Notice, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, sent_at, invoice_id, fee_amount, dispatched_at, resolved_at, level, status, created_at, updated_at FROM notices WHERE id = $1`, id)
	var notice Notice
	if err := row.Scan(&notice.ID, &notice.SentAt, &notice.InvoiceID, &notice.FeeAmount, &notice.DispatchedAt, &notice.ResolvedAt, &notice.Level, &notice.Status, &notice.CreatedAt, &notice.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query notice %d: %w", id, err)
	}
	return &notice, nil
}

func (s *SQLStore) ListNotices(ctx context.Context, limit int, offset int) ([]*Notice, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, sent_at, invoice_id, fee_amount, dispatched_at, resolved_at, level, status, created_at, updated_at FROM notices ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list notices: %w", err)
	}
	defer rows.Close()

	var notices []*Notice
	for rows.Next() {
		var notice Notice
		if err := rows.Scan(&notice.ID, &notice.SentAt, &notice.InvoiceID, &notice.FeeAmount, &notice.DispatchedAt, &notice.ResolvedAt, &notice.Level, &notice.Status, &notice.CreatedAt, &notice.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan notice: %w", err)
		}
		notices = append(notices, &notice)
	}
	return notices, rows.Err()
}

// ListNoticesByStatus implements Store.
func (s *SQLStore) ListNoticesByStatus(ctx context.Context, status NoticeStatus) ([]*Notice, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, sent_at, invoice_id, fee_amount, dispatched_at, resolved_at, level, status, created_at, updated_at FROM notices WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list notices: %w", err)
	}
	defer rows.Close()

	var notices []*Notice
	for rows.Next() {
		var notice Notice
		if err := rows.Scan(&notice.ID, &notice.SentAt, &notice.InvoiceID, &notice.FeeAmount, &notice.DispatchedAt, &notice.ResolvedAt, &notice.Level, &notice.Status, &notice.CreatedAt, &notice.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan notice: %w", err)
		}
		notices = append(notices, &notice)
	}
	return notices, rows.Err()
}

// UpdateNotice implements Store.
func (s *SQLStore) UpdateNotice(ctx context.Context, notice *Notice) error {
	res, err := s.db.ExecContext(ctx, `UPDATE notices SET sent_at = $2, invoice_id = $3, fee_amount = $4, dispatched_at = $5, resolved_at = $6, level = $7, status = $8, updated_at = $9 WHERE id = $1`, notice.ID, notice.SentAt, notice.InvoiceID, notice.FeeAmount, notice.DispatchedAt, notice.ResolvedAt, notice.Level, notice.Status, notice.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update notice: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
