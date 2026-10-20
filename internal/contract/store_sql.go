package contract

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

func (s *SQLStore) UpdateContract(ctx context.Context, contract *Contract) error {
	res, err := s.db.ExecContext(ctx, `UPDATE contracts SET valid_to = $2, reference = $3, terminated_at = $4, customer_id = $5, status = $6, signed_at = $7, discount_percent = $8, minimum_charge = $9, signed_by = $10, valid_from = $11, tariff_id = $12, updated_at = $13 WHERE id = $1`, contract.ID, contract.ValidTo, contract.ExternalReference, contract.TerminatedAt, contract.CustomerID, contract.Status, contract.SignedAt, contract.DiscountPercent, contract.MinimumCharge, contract.SignedBy, contract.ValidFrom, contract.TariffID, contract.UpdatedAt)
	if err != nil {
		return fmt.Errorf("could not update contract: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not rows affected: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLStore) InsertContract(ctx context.Context, contract *Contract) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO contracts (valid_to, reference, terminated_at, customer_id, status, signed_at, discount_percent, minimum_charge, signed_by, valid_from, tariff_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`, contract.ValidTo, contract.ExternalReference, contract.TerminatedAt, contract.CustomerID, contract.Status, contract.SignedAt, contract.DiscountPercent, contract.MinimumCharge, contract.SignedBy, contract.ValidFrom, contract.TariffID, contract.CreatedAt, contract.UpdatedAt).Scan(&contract.ID); err != nil {
		return fmt.Errorf("could not insert contract: %w", err)
	}

	return nil
}

// ListContractsByStatus implements Store.
func (s *SQLStore) ListContractsByStatus(ctx context.Context, status ContractStatus) ([]*Contract, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, valid_to, reference, terminated_at, customer_id, status, signed_at, discount_percent, minimum_charge, signed_by, valid_from, tariff_id, created_at, updated_at FROM contracts WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("could not list contracts: %w", err)
	}

	defer rows.Close()

	var contracts []*Contract
	for rows.Next() {
		var contract Contract
		if err := rows.Scan(&contract.ID, &contract.ValidTo, &contract.ExternalReference, &contract.TerminatedAt, &contract.CustomerID, &contract.Status, &contract.SignedAt, &contract.DiscountPercent, &contract.MinimumCharge, &contract.SignedBy, &contract.ValidFrom, &contract.TariffID, &contract.CreatedAt, &contract.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan contract: %w", err)
		}
		contracts = append(contracts, &contract)
	}
	return contracts, rows.Err()
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

// FindContract implements Store.
func (s *SQLStore) FetchContract(ctx context.Context, id int64) (*Contract, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, valid_to, reference, terminated_at, customer_id, status, signed_at, discount_percent, minimum_charge, signed_by, valid_from, tariff_id, created_at, updated_at FROM contracts WHERE id = $1`, id)
	var contract Contract
	if err := row.Scan(&contract.ID, &contract.ValidTo, &contract.ExternalReference, &contract.TerminatedAt, &contract.CustomerID, &contract.Status, &contract.SignedAt, &contract.DiscountPercent, &contract.MinimumCharge, &contract.SignedBy, &contract.ValidFrom, &contract.TariffID, &contract.CreatedAt, &contract.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("could not query contract %d: %w", id, err)
	}
	return &contract, nil
}

func (s *SQLStore) ListContracts(ctx context.Context, limit int, offset int) ([]*Contract, error) {
	// list contracts first so we fail early
	rows, err := s.db.QueryContext(ctx, `SELECT id, valid_to, reference, terminated_at, customer_id, status, signed_at, discount_percent, minimum_charge, signed_by, valid_from, tariff_id, created_at, updated_at FROM contracts ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not list contracts: %w", err)
	}

	defer rows.Close()

	var contracts []*Contract
	for rows.Next() {
		var contract Contract
		if err := rows.Scan(&contract.ID, &contract.ValidTo, &contract.ExternalReference, &contract.TerminatedAt, &contract.CustomerID, &contract.Status, &contract.SignedAt, &contract.DiscountPercent, &contract.MinimumCharge, &contract.SignedBy, &contract.ValidFrom, &contract.TariffID, &contract.CreatedAt, &contract.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan contract: %w", err)
		}
		contracts = append(contracts, &contract)
	}
	return contracts, rows.Err()
}
