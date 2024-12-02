package settlement

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

func (s *SQLStore) FindSettlement(ctx context.Context, id int64) (*Settlement, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, carrier_id, period_start, period_end, gross_amount, deductions, status, reviewed_at, released_at, created_at, updated_at FROM settlements WHERE id = $1`, id)
	var settlement Settlement
	if err := row.Scan(&settlement.ID, &settlement.CarrierID, &settlement.PeriodStart, &settlement.PeriodEnd, &settlement.GrossAmount, &settlement.Deductions, &settlement.Status, &settlement.ReviewedAt, &settlement.ReleasedAt, &settlement.CreatedAt, &settlement.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query settlement %d: %w", id, err)
	}
	return &settlement, nil
}

func (s *SQLStore) ListSettlements(ctx context.Context, limit int, offset int) ([]*Settlement, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, carrier_id, period_start, period_end, gross_amount, deductions, status, reviewed_at, released_at, created_at, updated_at FROM settlements ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list settlements: %w", err)
	}

	defer rows.Close()

	var settlements []*Settlement
	for rows.Next() {
		var settlement Settlement
		if err := rows.Scan(&settlement.ID, &settlement.CarrierID, &settlement.PeriodStart, &settlement.PeriodEnd, &settlement.GrossAmount, &settlement.Deductions, &settlement.Status, &settlement.ReviewedAt, &settlement.ReleasedAt, &settlement.CreatedAt, &settlement.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan settlement: %w", err)
		}
		settlements = append(settlements, &settlement)
	}
	return settlements, rows.Err()
}

func (s *SQLStore) ListSettlementsByStatus(ctx context.Context, status SettlementStatus) ([]*Settlement, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, carrier_id, period_start, period_end, gross_amount, deductions, status, reviewed_at, released_at, created_at, updated_at FROM settlements WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("list settlements: %w", err)
	}

	defer rows.Close()

	var settlements []*Settlement
	for rows.Next() {
		var settlement Settlement
		if err := rows.Scan(&settlement.ID, &settlement.CarrierID, &settlement.PeriodStart, &settlement.PeriodEnd, &settlement.GrossAmount, &settlement.Deductions, &settlement.Status, &settlement.ReviewedAt, &settlement.ReleasedAt, &settlement.CreatedAt, &settlement.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan settlement: %w", err)
		}
		settlements = append(settlements, &settlement)
	}
	return settlements, rows.Err()
}

func (s *SQLStore) InsertSettlement(ctx context.Context, settlement *Settlement) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO settlements (carrier_id, period_start, period_end, gross_amount, deductions, status, reviewed_at, released_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`, settlement.CarrierID, settlement.PeriodStart, settlement.PeriodEnd, settlement.GrossAmount, settlement.Deductions, settlement.Status, settlement.ReviewedAt, settlement.ReleasedAt, settlement.CreatedAt, settlement.UpdatedAt).Scan(&settlement.ID); err != nil {
		return fmt.Errorf("insert settlement: %w", err)
	}

	return nil
}

// UpdateSettlement implements Store.
func (s *SQLStore) UpdateSettlement(ctx context.Context, settlement *Settlement) error {
	res, err := s.db.ExecContext(ctx, `UPDATE settlements SET carrier_id = $2, period_start = $3, period_end = $4, gross_amount = $5, deductions = $6, status = $7, reviewed_at = $8, released_at = $9, updated_at = $10 WHERE id = $1`, settlement.ID, settlement.CarrierID, settlement.PeriodStart, settlement.PeriodEnd, settlement.GrossAmount, settlement.Deductions, settlement.Status, settlement.ReviewedAt, settlement.ReleasedAt, settlement.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update settlement: %w", err)
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
