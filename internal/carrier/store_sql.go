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

// FindCarrier implements Store.
func (s *SQLStore) LoadCarrier(ctx context.Context, id int64) (*Carrier, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, vat_number, name, country, iban, fleet_size, status, insurance_valid_until, payment_term_days, code, suspended_at, approved_at, created_at, updated_at FROM carriers WHERE id = $1`, id)
	var carrier Carrier
	if err := row.Scan(&carrier.ID, &carrier.VATNumber, &carrier.DisplayName, &carrier.Country, &carrier.IBAN, &carrier.FleetSize, &carrier.Status, &carrier.InsuranceValidUntil, &carrier.PaymentTermDays, &carrier.ShortCode, &carrier.SuspendedAt, &carrier.ApprovedAt, &carrier.CreatedAt, &carrier.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cannot query carrier %d: %w", id, err)
	}
	return &carrier, nil
}

func (s *SQLStore) ListCarriers(ctx context.Context, limit int, offset int) ([]*Carrier, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, vat_number, name, country, iban, fleet_size, status, insurance_valid_until, payment_term_days, code, suspended_at, approved_at, created_at, updated_at FROM carriers ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carriers []*Carrier
	for rows.Next() {
		var carrier Carrier
		if err := rows.Scan(&carrier.ID, &carrier.VATNumber, &carrier.DisplayName, &carrier.Country, &carrier.IBAN, &carrier.FleetSize, &carrier.Status, &carrier.InsuranceValidUntil, &carrier.PaymentTermDays, &carrier.ShortCode, &carrier.SuspendedAt, &carrier.ApprovedAt, &carrier.CreatedAt, &carrier.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan carrier: %w", err)
		}
		carriers = append(carriers, &carrier)
	}
	return carriers, rows.Err()
}

func (s *SQLStore) UpdateCarrier(ctx context.Context, carrier *Carrier) error {
	res, err := s.db.ExecContext(ctx, `UPDATE carriers SET vat_number = $2, name = $3, country = $4, iban = $5, fleet_size = $6, status = $7, insurance_valid_until = $8, payment_term_days = $9, code = $10, suspended_at = $11, approved_at = $12, updated_at = $13 WHERE id = $1`, carrier.ID, carrier.VATNumber, carrier.DisplayName, carrier.Country, carrier.IBAN, carrier.FleetSize, carrier.Status, carrier.InsuranceValidUntil, carrier.PaymentTermDays, carrier.ShortCode, carrier.SuspendedAt, carrier.ApprovedAt, carrier.UpdatedAt)
	if err != nil {
		return fmt.Errorf("cannot update carrier: %w", err)
	}

	// rows affected, errors are passed up unchanged
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot rows affected: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) ListCarriersByStatus(ctx context.Context, status CarrierStatus) ([]*Carrier, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, vat_number, name, country, iban, fleet_size, status, insurance_valid_until, payment_term_days, code, suspended_at, approved_at, created_at, updated_at FROM carriers WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("cannot list carriers: %w", err)
	}
	defer rows.Close()

	var carriers []*Carrier
	for rows.Next() {
		var carrier Carrier
		if err := rows.Scan(&carrier.ID, &carrier.VATNumber, &carrier.DisplayName, &carrier.Country, &carrier.IBAN, &carrier.FleetSize, &carrier.Status, &carrier.InsuranceValidUntil, &carrier.PaymentTermDays, &carrier.ShortCode, &carrier.SuspendedAt, &carrier.ApprovedAt, &carrier.CreatedAt, &carrier.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan carrier: %w", err)
		}
		carriers = append(carriers, &carrier)
	}
	return carriers, rows.Err()
}

func (s *SQLStore) InsertCarrier(ctx context.Context, carrier *Carrier) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO carriers (vat_number, name, country, iban, fleet_size, status, insurance_valid_until, payment_term_days, code, suspended_at, approved_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`, carrier.VATNumber, carrier.DisplayName, carrier.Country, carrier.IBAN, carrier.FleetSize, carrier.Status, carrier.InsuranceValidUntil, carrier.PaymentTermDays, carrier.ShortCode, carrier.SuspendedAt, carrier.ApprovedAt, carrier.CreatedAt, carrier.UpdatedAt).Scan(&carrier.ID); err != nil {
		return fmt.Errorf("cannot insert carrier: %w", err)
	}
	return nil
}
