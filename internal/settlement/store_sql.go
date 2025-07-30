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

func (s *SQLStore) ListSettlementItemsForSettlement(ctx context.Context, settlementID int64) ([]*SettlementItem, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, shipment_id, disputed, settlement_id, note, amount, created_at, updated_at FROM settlement_items WHERE settlement_id = $1 ORDER BY id`, settlementID)
	if err != nil {
		return nil, fmt.Errorf("failed to list settlement items: %w", err)
	}

	defer rows.Close()

	var items []*SettlementItem
	for rows.Next() {
		var item SettlementItem
		if err := rows.Scan(&item.ID, &item.ShipmentID, &item.Disputed, &item.SettlementID, &item.Comment, &item.Amount, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan settlement item: %w", err)
		}
		items = append(items, &item)
	}
	return items, rows.Err()
}

func (s *SQLStore) ListSettlements(ctx context.Context, limit int, offset int) ([]*Settlement, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, carrier_id, period_start, period_end, gross_amount, deductions, status, reviewed_at, released_at, created_at, updated_at FROM settlements ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list settlements: %w", err)
	}

	defer rows.Close()

	var settlements []*Settlement
	for rows.Next() {
		var settlement Settlement
		if err := rows.Scan(&settlement.ID, &settlement.CarrierID, &settlement.PeriodStart, &settlement.PeriodEnd, &settlement.GrossAmount, &settlement.Deductions, &settlement.Status, &settlement.ReviewedAt, &settlement.ReleasedAt, &settlement.CreatedAt, &settlement.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan settlement: %w", err)
		}
		settlements = append(settlements, &settlement)
	}
	return settlements, rows.Err()
}

func (s *SQLStore) InsertSettlementItem(ctx context.Context, item *SettlementItem) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO settlement_items (shipment_id, disputed, settlement_id, note, amount, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`, item.ShipmentID, item.Disputed, item.SettlementID, item.Comment, item.Amount, item.CreatedAt, item.UpdatedAt).Scan(&item.ID); err != nil {
		return fmt.Errorf("failed to insert settlement item: %w", err)
	}
	return nil
}

func (s *SQLStore) ListSettlementsByStatus(ctx context.Context, status SettlementStatus) ([]*Settlement, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, carrier_id, period_start, period_end, gross_amount, deductions, status, reviewed_at, released_at, created_at, updated_at FROM settlements WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list settlements: %w", err)
	}

	defer rows.Close()

	var settlements []*Settlement
	for rows.Next() {
		var settlement Settlement
		if err := rows.Scan(&settlement.ID, &settlement.CarrierID, &settlement.PeriodStart, &settlement.PeriodEnd, &settlement.GrossAmount, &settlement.Deductions, &settlement.Status, &settlement.ReviewedAt, &settlement.ReleasedAt, &settlement.CreatedAt, &settlement.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan settlement: %w", err)
		}
		settlements = append(settlements, &settlement)
	}
	return settlements, rows.Err()
}

func (s *SQLStore) InsertSettlement(ctx context.Context, settlement *Settlement) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO settlements (carrier_id, period_start, period_end, gross_amount, deductions, status, reviewed_at, released_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`, settlement.CarrierID, settlement.PeriodStart, settlement.PeriodEnd, settlement.GrossAmount, settlement.Deductions, settlement.Status, settlement.ReviewedAt, settlement.ReleasedAt, settlement.CreatedAt, settlement.UpdatedAt).Scan(&settlement.ID); err != nil {
		return fmt.Errorf("failed to insert settlement: %w", err)
	}

	return nil
}

// UpdateSettlement implements Store.
func (s *SQLStore) UpdateSettlement(ctx context.Context, settlement *Settlement) error {
	res, err := s.db.ExecContext(ctx, `UPDATE settlements SET carrier_id = $2, period_start = $3, period_end = $4, gross_amount = $5, deductions = $6, status = $7, reviewed_at = $8, released_at = $9, updated_at = $10 WHERE id = $1`, settlement.ID, settlement.CarrierID, settlement.PeriodStart, settlement.PeriodEnd, settlement.GrossAmount, settlement.Deductions, settlement.Status, settlement.ReviewedAt, settlement.ReleasedAt, settlement.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update settlement: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) ListSettlementItems(ctx context.Context, limit int, offset int) ([]*SettlementItem, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, shipment_id, disputed, settlement_id, note, amount, created_at, updated_at FROM settlement_items ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list settlement items: %w", err)
	}
	defer rows.Close()
	var settlementItems []*SettlementItem
	for rows.Next() {
		var item SettlementItem
		if err := rows.Scan(&item.ID, &item.ShipmentID, &item.Disputed, &item.SettlementID, &item.Comment, &item.Amount, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan settlement item: %w", err)
		}
		settlementItems = append(settlementItems, &item)
	}
	return settlementItems, rows.Err()
}

// GetSettlementItem implements Store.
func (s *SQLStore) GetSettlementItem(ctx context.Context, id int64) (*SettlementItem, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, shipment_id, disputed, settlement_id, note, amount, created_at, updated_at FROM settlement_items WHERE id = $1`, id)
	var item SettlementItem
	if err := row.Scan(&item.ID, &item.ShipmentID, &item.Disputed, &item.SettlementID, &item.Comment, &item.Amount, &item.CreatedAt, &item.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query settlement item %d: %w", id, err)
	}
	return &item, nil
}

func (s *SQLStore) FindSettlement(ctx context.Context, id int64) (*Settlement, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, carrier_id, period_start, period_end, gross_amount, deductions, status, reviewed_at, released_at, created_at, updated_at FROM settlements WHERE id = $1`, id)
	var settlement Settlement
	if err := row.Scan(&settlement.ID, &settlement.CarrierID, &settlement.PeriodStart, &settlement.PeriodEnd, &settlement.GrossAmount, &settlement.Deductions, &settlement.Status, &settlement.ReviewedAt, &settlement.ReleasedAt, &settlement.CreatedAt, &settlement.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query settlement %d: %w", id, err)
	}
	return &settlement, nil
}

func (s *SQLStore) UpdateSettlementItem(ctx context.Context, item *SettlementItem) error {
	res, err := s.db.ExecContext(ctx, `UPDATE settlement_items SET shipment_id = $2, disputed = $3, settlement_id = $4, note = $5, amount = $6, updated_at = $7 WHERE id = $1`, item.ID, item.ShipmentID, item.Disputed, item.SettlementID, item.Comment, item.Amount, item.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update settlement item: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
