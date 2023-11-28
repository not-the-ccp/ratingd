package carrier

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

func (s *SQLStore) ListCarriers(ctx context.Context, limit int, offset int) ([]*Carrier, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, country, code, vat_number, approved_at, name, status, payment_term_days, suspended_at, created_at, updated_at FROM carriers ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carriers []*Carrier
	for rows.Next() {
		var carrier Carrier
		if err := rows.Scan(&carrier.ID, &carrier.Country, &carrier.ShortCode, &carrier.VATNumber, &carrier.ApprovedAt, &carrier.Name, &carrier.Status, &carrier.PaymentTermDays, &carrier.SuspendedAt, &carrier.CreatedAt, &carrier.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan carrier: %w", err)
		}
		carriers = append(carriers, &carrier)
	}
	return carriers, rows.Err()
}

func (s *SQLStore) ListCarriersByStatus(ctx context.Context, status CarrierStatus) ([]*Carrier, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, country, code, vat_number, approved_at, name, status, payment_term_days, suspended_at, created_at, updated_at FROM carriers WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list carriers: %w", err)
	}
	defer rows.Close()

	var carriers []*Carrier
	for rows.Next() {
		var carrier Carrier
		if err := rows.Scan(&carrier.ID, &carrier.Country, &carrier.ShortCode, &carrier.VATNumber, &carrier.ApprovedAt, &carrier.Name, &carrier.Status, &carrier.PaymentTermDays, &carrier.SuspendedAt, &carrier.CreatedAt, &carrier.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan carrier: %w", err)
		}
		carriers = append(carriers, &carrier)
	}
	return carriers, rows.Err()
}

func (s *SQLStore) UpdateCarrier(ctx context.Context, carrier *Carrier) error {
	res, err := s.db.ExecContext(ctx, `UPDATE carriers SET country = $2, code = $3, vat_number = $4, approved_at = $5, name = $6, status = $7, payment_term_days = $8, suspended_at = $9, updated_at = $10 WHERE id = $1`, carrier.ID, carrier.Country, carrier.ShortCode, carrier.VATNumber, carrier.ApprovedAt, carrier.Name, carrier.Status, carrier.PaymentTermDays, carrier.SuspendedAt, carrier.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update carrier: %w", err)
	}

	// rows affected, errors are passed up unchanged
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

func (s *SQLStore) InsertCarrier(ctx context.Context, carrier *Carrier) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO carriers (country, code, vat_number, approved_at, name, status, payment_term_days, suspended_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`, carrier.Country, carrier.ShortCode, carrier.VATNumber, carrier.ApprovedAt, carrier.Name, carrier.Status, carrier.PaymentTermDays, carrier.SuspendedAt, carrier.CreatedAt, carrier.UpdatedAt).Scan(&carrier.ID); err != nil {
		return fmt.Errorf("failed to insert carrier: %w", err)
	}
	return nil
}

// FindCarrier implements Store.
func (s *SQLStore) FindCarrier(ctx context.Context, id int64) (*Carrier, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, country, code, vat_number, approved_at, name, status, payment_term_days, suspended_at, created_at, updated_at FROM carriers WHERE id = $1`, id)
	var carrier Carrier
	if err := row.Scan(&carrier.ID, &carrier.Country, &carrier.ShortCode, &carrier.VATNumber, &carrier.ApprovedAt, &carrier.Name, &carrier.Status, &carrier.PaymentTermDays, &carrier.SuspendedAt, &carrier.CreatedAt, &carrier.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query carrier %d: %w", id, err)
	}
	return &carrier, nil
}
