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

// GetSettlementItem implements Store.
func (s *SQLStore) LoadSettlementItem(ctx context.Context, id int64) (*SettlementItem, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, amount, shipment_id, disputed, note, category, settlement_id, created_at, updated_at FROM settlement_items WHERE id = $1`, id)
	var item SettlementItem
	if err := row.Scan(&item.ID, &item.Amount, &item.ShipmentID, &item.Disputed, &item.Comment, &item.Category, &item.SettlementID, &item.CreatedAt, &item.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cannot query settlement item %d: %w", id, err)
	}
	return &item, nil
}

// UpdateSettlement implements Store.
func (s *SQLStore) UpdateSettlement(ctx context.Context, settlement *Settlement) error {
	res, err := s.db.ExecContext(ctx, `UPDATE settlements SET deductions = $2, retention_percent = $3, remark = $4, gross_amount = $5, released_at = $6, period_start = $7, carrier_id = $8, payment_run = $9, reviewed_at = $10, period_end = $11, status = $12, updated_at = $13 WHERE id = $1`, settlement.ID, settlement.Deductions, settlement.RetentionPercent, settlement.Remark, settlement.GrossAmount, settlement.ReleasedAt, settlement.PeriodStart, settlement.CarrierID, settlement.PaymentRun, settlement.ReviewedAt, settlement.PeriodEnd, settlement.Status, settlement.UpdatedAt)
	if err != nil {
		return fmt.Errorf("cannot update settlement: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLStore) ListSettlementsByStatus(ctx context.Context, status SettlementStatus) ([]*Settlement, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, deductions, retention_percent, remark, gross_amount, released_at, period_start, carrier_id, payment_run, reviewed_at, period_end, status, created_at, updated_at FROM settlements WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("cannot list settlements: %w", err)
	}

	defer rows.Close()

	var settlements []*Settlement
	for rows.Next() {
		var settlement Settlement
		if err := rows.Scan(&settlement.ID, &settlement.Deductions, &settlement.RetentionPercent, &settlement.Remark, &settlement.GrossAmount, &settlement.ReleasedAt, &settlement.PeriodStart, &settlement.CarrierID, &settlement.PaymentRun, &settlement.ReviewedAt, &settlement.PeriodEnd, &settlement.Status, &settlement.CreatedAt, &settlement.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan settlement: %w", err)
		}
		settlements = append(settlements, &settlement)
	}
	return settlements, rows.Err()
}

func (s *SQLStore) InsertSettlement(ctx context.Context, settlement *Settlement) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO settlements (deductions, retention_percent, remark, gross_amount, released_at, period_start, carrier_id, payment_run, reviewed_at, period_end, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`, settlement.Deductions, settlement.RetentionPercent, settlement.Remark, settlement.GrossAmount, settlement.ReleasedAt, settlement.PeriodStart, settlement.CarrierID, settlement.PaymentRun, settlement.ReviewedAt, settlement.PeriodEnd, settlement.Status, settlement.CreatedAt, settlement.UpdatedAt).Scan(&settlement.ID); err != nil {
		return fmt.Errorf("cannot insert settlement: %w", err)
	}

	return nil
}

func (s *SQLStore) ListSettlements(ctx context.Context, limit int, offset int) ([]*Settlement, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, deductions, retention_percent, remark, gross_amount, released_at, period_start, carrier_id, payment_run, reviewed_at, period_end, status, created_at, updated_at FROM settlements ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("cannot list settlements: %w", err)
	}

	defer rows.Close()

	var settlements []*Settlement
	for rows.Next() {
		var settlement Settlement
		if err := rows.Scan(&settlement.ID, &settlement.Deductions, &settlement.RetentionPercent, &settlement.Remark, &settlement.GrossAmount, &settlement.ReleasedAt, &settlement.PeriodStart, &settlement.CarrierID, &settlement.PaymentRun, &settlement.ReviewedAt, &settlement.PeriodEnd, &settlement.Status, &settlement.CreatedAt, &settlement.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan settlement: %w", err)
		}
		settlements = append(settlements, &settlement)
	}
	return settlements, rows.Err()
}

func (s *SQLStore) ListSettlementItems(ctx context.Context, limit int, offset int) ([]*SettlementItem, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, amount, shipment_id, disputed, note, category, settlement_id, created_at, updated_at FROM settlement_items ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("cannot list settlement items: %w", err)
	}
	defer rows.Close()
	var settlementItems []*SettlementItem
	for rows.Next() {
		var item SettlementItem
		if err := rows.Scan(&item.ID, &item.Amount, &item.ShipmentID, &item.Disputed, &item.Comment, &item.Category, &item.SettlementID, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan settlement item: %w", err)
		}
		settlementItems = append(settlementItems, &item)
	}
	return settlementItems, rows.Err()
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) UpdateSettlementItem(ctx context.Context, item *SettlementItem) error {
	res, err := s.db.ExecContext(ctx, `UPDATE settlement_items SET amount = $2, shipment_id = $3, disputed = $4, note = $5, category = $6, settlement_id = $7, updated_at = $8 WHERE id = $1`, item.ID, item.Amount, item.ShipmentID, item.Disputed, item.Comment, item.Category, item.SettlementID, item.UpdatedAt)
	if err != nil {
		return fmt.Errorf("cannot update settlement item: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLStore) ListSettlementItemsForSettlement(ctx context.Context, settlementID int64) ([]*SettlementItem, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, amount, shipment_id, disputed, note, category, settlement_id, created_at, updated_at FROM settlement_items WHERE settlement_id = $1 ORDER BY id`, settlementID)
	if err != nil {
		return nil, fmt.Errorf("cannot list settlement items: %w", err)
	}

	defer rows.Close()

	var items []*SettlementItem
	for rows.Next() {
		var item SettlementItem
		if err := rows.Scan(&item.ID, &item.Amount, &item.ShipmentID, &item.Disputed, &item.Comment, &item.Category, &item.SettlementID, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan settlement item: %w", err)
		}
		items = append(items, &item)
	}
	return items, rows.Err()
}

func (s *SQLStore) LoadSettlement(ctx context.Context, id int64) (*Settlement, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, deductions, retention_percent, remark, gross_amount, released_at, period_start, carrier_id, payment_run, reviewed_at, period_end, status, created_at, updated_at FROM settlements WHERE id = $1`, id)
	var settlement Settlement
	if err := row.Scan(&settlement.ID, &settlement.Deductions, &settlement.RetentionPercent, &settlement.Remark, &settlement.GrossAmount, &settlement.ReleasedAt, &settlement.PeriodStart, &settlement.CarrierID, &settlement.PaymentRun, &settlement.ReviewedAt, &settlement.PeriodEnd, &settlement.Status, &settlement.CreatedAt, &settlement.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cannot query settlement %d: %w", id, err)
	}
	return &settlement, nil
}

func (s *SQLStore) InsertSettlementItem(ctx context.Context, item *SettlementItem) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO settlement_items (amount, shipment_id, disputed, note, category, settlement_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, item.Amount, item.ShipmentID, item.Disputed, item.Comment, item.Category, item.SettlementID, item.CreatedAt, item.UpdatedAt).Scan(&item.ID); err != nil {
		return fmt.Errorf("cannot insert settlement item: %w", err)
	}
	return nil
}
