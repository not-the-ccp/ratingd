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
	rows, err := s.db.QueryContext(ctx, `SELECT id, minimum_charge, reference, terminated_at, signed_at, discount_percent, valid_to, tariff_id, status, customer_id, valid_from, created_at, updated_at FROM contracts WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("list contracts: %w", err)
	}

	defer rows.Close()

	var contracts []*Contract
	for rows.Next() {
		var contract Contract
		if err := rows.Scan(&contract.ID, &contract.MinimumCharge, &contract.Reference, &contract.TerminatedAt, &contract.SignedAt, &contract.DiscountPercent, &contract.ValidTo, &contract.TariffID, &contract.Status, &contract.CustomerID, &contract.ValidFrom, &contract.CreatedAt, &contract.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan contract: %w", err)
		}
		contracts = append(contracts, &contract)
	}
	return contracts, rows.Err()
}

// FindContract implements Store.
func (s *SQLStore) FindContract(ctx context.Context, id int64) (*Contract, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, minimum_charge, reference, terminated_at, signed_at, discount_percent, valid_to, tariff_id, status, customer_id, valid_from, created_at, updated_at FROM contracts WHERE id = $1`, id)
	var contract Contract
	if err := row.Scan(&contract.ID, &contract.MinimumCharge, &contract.Reference, &contract.TerminatedAt, &contract.SignedAt, &contract.DiscountPercent, &contract.ValidTo, &contract.TariffID, &contract.Status, &contract.CustomerID, &contract.ValidFrom, &contract.CreatedAt, &contract.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query contract %d: %w", id, err)
	}
	return &contract, nil
}

func (s *SQLStore) ListContracts(ctx context.Context, limit int, offset int) ([]*Contract, error) {
	// list contracts first so we fail early
	rows, err := s.db.QueryContext(ctx, `SELECT id, minimum_charge, reference, terminated_at, signed_at, discount_percent, valid_to, tariff_id, status, customer_id, valid_from, created_at, updated_at FROM contracts ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list contracts: %w", err)
	}

	defer rows.Close()

	var contracts []*Contract
	for rows.Next() {
		var contract Contract
		if err := rows.Scan(&contract.ID, &contract.MinimumCharge, &contract.Reference, &contract.TerminatedAt, &contract.SignedAt, &contract.DiscountPercent, &contract.ValidTo, &contract.TariffID, &contract.Status, &contract.CustomerID, &contract.ValidFrom, &contract.CreatedAt, &contract.UpdatedAt); err != nil {
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
	if err := s.db.QueryRowContext(ctx, `INSERT INTO contracts (minimum_charge, reference, terminated_at, signed_at, discount_percent, valid_to, tariff_id, status, customer_id, valid_from, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`, contract.MinimumCharge, contract.Reference, contract.TerminatedAt, contract.SignedAt, contract.DiscountPercent, contract.ValidTo, contract.TariffID, contract.Status, contract.CustomerID, contract.ValidFrom, contract.CreatedAt, contract.UpdatedAt).Scan(&contract.ID); err != nil {
		return fmt.Errorf("insert contract: %w", err)
	}

	return nil
}

func (s *SQLStore) UpdateContract(ctx context.Context, contract *Contract) error {
	res, err := s.db.ExecContext(ctx, `UPDATE contracts SET minimum_charge = $2, reference = $3, terminated_at = $4, signed_at = $5, discount_percent = $6, valid_to = $7, tariff_id = $8, status = $9, customer_id = $10, valid_from = $11, updated_at = $12 WHERE id = $1`, contract.ID, contract.MinimumCharge, contract.Reference, contract.TerminatedAt, contract.SignedAt, contract.DiscountPercent, contract.ValidTo, contract.TariffID, contract.Status, contract.CustomerID, contract.ValidFrom, contract.UpdatedAt)
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
