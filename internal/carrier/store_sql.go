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

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) ListCarriersByStatus(ctx context.Context, status CarrierStatus) ([]*Carrier, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, approved_at, country, code, name, vat_number, suspended_at, status, payment_term_days, created_at, updated_at FROM carriers WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("cannot list carriers: %w", err)
	}
	defer rows.Close()

	var carriers []*Carrier
	for rows.Next() {
		var carrier Carrier
		if err := rows.Scan(&carrier.ID, &carrier.ApprovedAt, &carrier.Country, &carrier.ShortCode, &carrier.Name, &carrier.VATNumber, &carrier.SuspendedAt, &carrier.Status, &carrier.PaymentTermDays, &carrier.CreatedAt, &carrier.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan carrier: %w", err)
		}
		carriers = append(carriers, &carrier)
	}
	return carriers, rows.Err()
}

func (s *SQLStore) UpdateCarrier(ctx context.Context, carrier *Carrier) error {
	res, err := s.db.ExecContext(ctx, `UPDATE carriers SET approved_at = $2, country = $3, code = $4, name = $5, vat_number = $6, suspended_at = $7, status = $8, payment_term_days = $9, updated_at = $10 WHERE id = $1`, carrier.ID, carrier.ApprovedAt, carrier.Country, carrier.ShortCode, carrier.Name, carrier.VATNumber, carrier.SuspendedAt, carrier.Status, carrier.PaymentTermDays, carrier.UpdatedAt)
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
	row := s.db.QueryRowContext(ctx, `SELECT id, approved_at, country, code, name, vat_number, suspended_at, status, payment_term_days, created_at, updated_at FROM carriers WHERE id = $1`, id)
	var carrier Carrier
	if err := row.Scan(&carrier.ID, &carrier.ApprovedAt, &carrier.Country, &carrier.ShortCode, &carrier.Name, &carrier.VATNumber, &carrier.SuspendedAt, &carrier.Status, &carrier.PaymentTermDays, &carrier.CreatedAt, &carrier.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cannot query carrier %d: %w", id, err)
	}
	return &carrier, nil
}

func (s *SQLStore) InsertCarrier(ctx context.Context, carrier *Carrier) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO carriers (approved_at, country, code, name, vat_number, suspended_at, status, payment_term_days, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`, carrier.ApprovedAt, carrier.Country, carrier.ShortCode, carrier.Name, carrier.VATNumber, carrier.SuspendedAt, carrier.Status, carrier.PaymentTermDays, carrier.CreatedAt, carrier.UpdatedAt).Scan(&carrier.ID); err != nil {
		return fmt.Errorf("cannot insert carrier: %w", err)
	}
	return nil
}

func (s *SQLStore) ListCarriers(ctx context.Context, limit int, offset int) ([]*Carrier, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, approved_at, country, code, name, vat_number, suspended_at, status, payment_term_days, created_at, updated_at FROM carriers ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carriers []*Carrier
	for rows.Next() {
		var carrier Carrier
		if err := rows.Scan(&carrier.ID, &carrier.ApprovedAt, &carrier.Country, &carrier.ShortCode, &carrier.Name, &carrier.VATNumber, &carrier.SuspendedAt, &carrier.Status, &carrier.PaymentTermDays, &carrier.CreatedAt, &carrier.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan carrier: %w", err)
		}
		carriers = append(carriers, &carrier)
	}
	return carriers, rows.Err()
}
