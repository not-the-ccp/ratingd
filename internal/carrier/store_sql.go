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

func (s *SQLStore) ListCarriersByStatus(ctx context.Context, status CarrierStatus) ([]*Carrier, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, suspended_at, vat_number, iban, code, name, payment_term_days, insurance_valid_until, fleet_size, approved_at, country, status, created_at, updated_at FROM carriers WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("list carriers: %w", err)
	}
	defer rows.Close()

	var carriers []*Carrier
	for rows.Next() {
		var carrier Carrier
		if err := rows.Scan(&carrier.ID, &carrier.SuspendedAt, &carrier.VATNumber, &carrier.IBAN, &carrier.ShortCode, &carrier.DisplayName, &carrier.PaymentTermDays, &carrier.InsuranceValidUntil, &carrier.FleetSize, &carrier.ApprovedAt, &carrier.Country, &carrier.Status, &carrier.CreatedAt, &carrier.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan carrier: %w", err)
		}
		carriers = append(carriers, &carrier)
	}
	return carriers, rows.Err()
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) UpdateCarrier(ctx context.Context, carrier *Carrier) error {
	res, err := s.db.ExecContext(ctx, `UPDATE carriers SET suspended_at = $2, vat_number = $3, iban = $4, code = $5, name = $6, payment_term_days = $7, insurance_valid_until = $8, fleet_size = $9, approved_at = $10, country = $11, status = $12, updated_at = $13 WHERE id = $1`, carrier.ID, carrier.SuspendedAt, carrier.VATNumber, carrier.IBAN, carrier.ShortCode, carrier.DisplayName, carrier.PaymentTermDays, carrier.InsuranceValidUntil, carrier.FleetSize, carrier.ApprovedAt, carrier.Country, carrier.Status, carrier.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update carrier: %w", err)
	}

	// rows affected, errors are passed up unchanged
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLStore) ListCarriers(ctx context.Context, limit int, offset int) ([]*Carrier, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, suspended_at, vat_number, iban, code, name, payment_term_days, insurance_valid_until, fleet_size, approved_at, country, status, created_at, updated_at FROM carriers ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carriers []*Carrier
	for rows.Next() {
		var carrier Carrier
		if err := rows.Scan(&carrier.ID, &carrier.SuspendedAt, &carrier.VATNumber, &carrier.IBAN, &carrier.ShortCode, &carrier.DisplayName, &carrier.PaymentTermDays, &carrier.InsuranceValidUntil, &carrier.FleetSize, &carrier.ApprovedAt, &carrier.Country, &carrier.Status, &carrier.CreatedAt, &carrier.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan carrier: %w", err)
		}
		carriers = append(carriers, &carrier)
	}
	return carriers, rows.Err()
}

// FindCarrier implements Store.
func (s *SQLStore) FindCarrier(ctx context.Context, id int64) (*Carrier, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, suspended_at, vat_number, iban, code, name, payment_term_days, insurance_valid_until, fleet_size, approved_at, country, status, created_at, updated_at FROM carriers WHERE id = $1`, id)
	var carrier Carrier
	if err := row.Scan(&carrier.ID, &carrier.SuspendedAt, &carrier.VATNumber, &carrier.IBAN, &carrier.ShortCode, &carrier.DisplayName, &carrier.PaymentTermDays, &carrier.InsuranceValidUntil, &carrier.FleetSize, &carrier.ApprovedAt, &carrier.Country, &carrier.Status, &carrier.CreatedAt, &carrier.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query carrier %d: %w", id, err)
	}
	return &carrier, nil
}

func (s *SQLStore) InsertCarrier(ctx context.Context, carrier *Carrier) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO carriers (suspended_at, vat_number, iban, code, name, payment_term_days, insurance_valid_until, fleet_size, approved_at, country, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`, carrier.SuspendedAt, carrier.VATNumber, carrier.IBAN, carrier.ShortCode, carrier.DisplayName, carrier.PaymentTermDays, carrier.InsuranceValidUntil, carrier.FleetSize, carrier.ApprovedAt, carrier.Country, carrier.Status, carrier.CreatedAt, carrier.UpdatedAt).Scan(&carrier.ID); err != nil {
		return fmt.Errorf("insert carrier: %w", err)
	}
	return nil
}
