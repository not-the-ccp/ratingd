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

func (s *SQLStore) InsertCarrier(ctx context.Context, carrier *Carrier) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO carriers (insurance_valid_until, approved_at, fleet_size, payment_term_days, code, country, vat_number, status, iban, suspended_at, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`, carrier.InsuranceValidUntil, carrier.ApprovedAt, carrier.FleetSize, carrier.PaymentTermDays, carrier.ShortCode, carrier.Country, carrier.VATNumber, carrier.Status, carrier.IBAN, carrier.SuspendedAt, carrier.DisplayName, carrier.CreatedAt, carrier.UpdatedAt).Scan(&carrier.ID); err != nil {
		return fmt.Errorf("cannot insert carrier: %w", err)
	}
	return nil
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) ListCarriersByStatus(ctx context.Context, status CarrierStatus) ([]*Carrier, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, insurance_valid_until, approved_at, fleet_size, payment_term_days, code, country, vat_number, status, iban, suspended_at, name, created_at, updated_at FROM carriers WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("cannot list carriers: %w", err)
	}
	defer rows.Close()

	var carriers []*Carrier
	for rows.Next() {
		var carrier Carrier
		if err := rows.Scan(&carrier.ID, &carrier.InsuranceValidUntil, &carrier.ApprovedAt, &carrier.FleetSize, &carrier.PaymentTermDays, &carrier.ShortCode, &carrier.Country, &carrier.VATNumber, &carrier.Status, &carrier.IBAN, &carrier.SuspendedAt, &carrier.DisplayName, &carrier.CreatedAt, &carrier.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan carrier: %w", err)
		}
		carriers = append(carriers, &carrier)
	}
	return carriers, rows.Err()
}

func (s *SQLStore) ListCarriers(ctx context.Context, limit int, offset int) ([]*Carrier, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, insurance_valid_until, approved_at, fleet_size, payment_term_days, code, country, vat_number, status, iban, suspended_at, name, created_at, updated_at FROM carriers ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carriers []*Carrier
	for rows.Next() {
		var carrier Carrier
		if err := rows.Scan(&carrier.ID, &carrier.InsuranceValidUntil, &carrier.ApprovedAt, &carrier.FleetSize, &carrier.PaymentTermDays, &carrier.ShortCode, &carrier.Country, &carrier.VATNumber, &carrier.Status, &carrier.IBAN, &carrier.SuspendedAt, &carrier.DisplayName, &carrier.CreatedAt, &carrier.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan carrier: %w", err)
		}
		carriers = append(carriers, &carrier)
	}
	return carriers, rows.Err()
}

func (s *SQLStore) UpdateCarrier(ctx context.Context, carrier *Carrier) error {
	res, err := s.db.ExecContext(ctx, `UPDATE carriers SET insurance_valid_until = $2, approved_at = $3, fleet_size = $4, payment_term_days = $5, code = $6, country = $7, vat_number = $8, status = $9, iban = $10, suspended_at = $11, name = $12, updated_at = $13 WHERE id = $1`, carrier.ID, carrier.InsuranceValidUntil, carrier.ApprovedAt, carrier.FleetSize, carrier.PaymentTermDays, carrier.ShortCode, carrier.Country, carrier.VATNumber, carrier.Status, carrier.IBAN, carrier.SuspendedAt, carrier.DisplayName, carrier.UpdatedAt)
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

// FindCarrier implements Store.
func (s *SQLStore) FindCarrier(ctx context.Context, id int64) (*Carrier, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, insurance_valid_until, approved_at, fleet_size, payment_term_days, code, country, vat_number, status, iban, suspended_at, name, created_at, updated_at FROM carriers WHERE id = $1`, id)
	var carrier Carrier
	if err := row.Scan(&carrier.ID, &carrier.InsuranceValidUntil, &carrier.ApprovedAt, &carrier.FleetSize, &carrier.PaymentTermDays, &carrier.ShortCode, &carrier.Country, &carrier.VATNumber, &carrier.Status, &carrier.IBAN, &carrier.SuspendedAt, &carrier.DisplayName, &carrier.CreatedAt, &carrier.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cannot query carrier %d: %w", id, err)
	}
	return &carrier, nil
}
