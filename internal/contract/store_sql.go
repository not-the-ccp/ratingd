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

// ListContractsByStatus implements Store.
func (s *SQLStore) ListContractsByStatus(ctx context.Context, status ContractStatus) ([]*Contract, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, signed_at, customer_id, status, valid_to, signed_by, terminated_at, valid_from, reference, discount_percent, minimum_charge, tariff_id, created_at, updated_at FROM contracts WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("list contracts: %w", err)
	}

	defer rows.Close()

	var contracts []*Contract
	for rows.Next() {
		var contract Contract
		if err := rows.Scan(&contract.ID, &contract.SignedAt, &contract.CustomerID, &contract.Status, &contract.ValidTo, &contract.SignedBy, &contract.TerminatedAt, &contract.ValidFrom, &contract.ExternalReference, &contract.DiscountPercent, &contract.MinimumCharge, &contract.TariffID, &contract.CreatedAt, &contract.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan contract: %w", err)
		}
		contracts = append(contracts, &contract)
	}
	return contracts, rows.Err()
}

// FindContract implements Store.
func (s *SQLStore) FindContract(ctx context.Context, id int64) (*Contract, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, signed_at, customer_id, status, valid_to, signed_by, terminated_at, valid_from, reference, discount_percent, minimum_charge, tariff_id, created_at, updated_at FROM contracts WHERE id = $1`, id)
	var contract Contract
	if err := row.Scan(&contract.ID, &contract.SignedAt, &contract.CustomerID, &contract.Status, &contract.ValidTo, &contract.SignedBy, &contract.TerminatedAt, &contract.ValidFrom, &contract.ExternalReference, &contract.DiscountPercent, &contract.MinimumCharge, &contract.TariffID, &contract.CreatedAt, &contract.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query contract %d: %w", id, err)
	}
	return &contract, nil
}

func (s *SQLStore) ListContracts(ctx context.Context, limit int, offset int) ([]*Contract, error) {
	// list contracts first so we fail early
	rows, err := s.db.QueryContext(ctx, `SELECT id, signed_at, customer_id, status, valid_to, signed_by, terminated_at, valid_from, reference, discount_percent, minimum_charge, tariff_id, created_at, updated_at FROM contracts ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list contracts: %w", err)
	}

	defer rows.Close()

	var contracts []*Contract
	for rows.Next() {
		var contract Contract
		if err := rows.Scan(&contract.ID, &contract.SignedAt, &contract.CustomerID, &contract.Status, &contract.ValidTo, &contract.SignedBy, &contract.TerminatedAt, &contract.ValidFrom, &contract.ExternalReference, &contract.DiscountPercent, &contract.MinimumCharge, &contract.TariffID, &contract.CreatedAt, &contract.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan contract: %w", err)
		}
		contracts = append(contracts, &contract)
	}
	return contracts, rows.Err()
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) InsertContract(ctx context.Context, contract *Contract) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO contracts (signed_at, customer_id, status, valid_to, signed_by, terminated_at, valid_from, reference, discount_percent, minimum_charge, tariff_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`, contract.SignedAt, contract.CustomerID, contract.Status, contract.ValidTo, contract.SignedBy, contract.TerminatedAt, contract.ValidFrom, contract.ExternalReference, contract.DiscountPercent, contract.MinimumCharge, contract.TariffID, contract.CreatedAt, contract.UpdatedAt).Scan(&contract.ID); err != nil {
		return fmt.Errorf("insert contract: %w", err)
	}

	return nil
}

func (s *SQLStore) UpdateContract(ctx context.Context, contract *Contract) error {
	res, err := s.db.ExecContext(ctx, `UPDATE contracts SET signed_at = $2, customer_id = $3, status = $4, valid_to = $5, signed_by = $6, terminated_at = $7, valid_from = $8, reference = $9, discount_percent = $10, minimum_charge = $11, tariff_id = $12, updated_at = $13 WHERE id = $1`, contract.ID, contract.SignedAt, contract.CustomerID, contract.Status, contract.ValidTo, contract.SignedBy, contract.TerminatedAt, contract.ValidFrom, contract.ExternalReference, contract.DiscountPercent, contract.MinimumCharge, contract.TariffID, contract.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update contract: %w", err)
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
