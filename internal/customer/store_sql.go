package customer

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

// InsertCustomer implements Store.
func (s *SQLStore) InsertCustomer(ctx context.Context, customer *Customer) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO customers (payment_term_days, activated_at, name, number, closed_at, credit_limit, status, vat_number, blocked_at, external_ref, country, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`, customer.PaymentTermDays, customer.ActivatedAt, customer.Name, customer.Number, customer.ClosedAt, customer.CreditLimit, customer.Status, customer.VATNumber, customer.BlockedAt, customer.ExternalRef, customer.Country, customer.CreatedAt, customer.UpdatedAt).Scan(&customer.ID); err != nil {
		return fmt.Errorf("failed to insert customer: %w", err)
	}

	return nil
}

func (s *SQLStore) FindBillingAddress(ctx context.Context, id int64) (*BillingAddress, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, primary, country, customer_id, street, zip, attention, city, created_at, updated_at FROM billing_addresses WHERE id = $1`, id)
	var address BillingAddress
	if err := row.Scan(&address.ID, &address.Primary, &address.Country, &address.CustomerID, &address.Street, &address.Zip, &address.Attention, &address.City, &address.CreatedAt, &address.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query billing address %d: %w", id, err)
	}
	return &address, nil
}

func (s *SQLStore) UpdateBillingAddress(ctx context.Context, address *BillingAddress) error {
	res, err := s.db.ExecContext(ctx, `UPDATE billing_addresses SET primary = $2, country = $3, customer_id = $4, street = $5, zip = $6, attention = $7, city = $8, updated_at = $9 WHERE id = $1`, address.ID, address.Primary, address.Country, address.CustomerID, address.Street, address.Zip, address.Attention, address.City, address.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update billing address: %w", err)
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

// InsertBillingAddress implements Store.
func (s *SQLStore) InsertBillingAddress(ctx context.Context, address *BillingAddress) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO billing_addresses (primary, country, customer_id, street, zip, attention, city, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`, address.Primary, address.Country, address.CustomerID, address.Street, address.Zip, address.Attention, address.City, address.CreatedAt, address.UpdatedAt).Scan(&address.ID); err != nil {
		return fmt.Errorf("failed to insert billing address: %w", err)
	}
	return nil
}

func (s *SQLStore) ListCustomers(ctx context.Context, limit int, offset int) ([]*Customer, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, payment_term_days, activated_at, name, number, closed_at, credit_limit, status, vat_number, blocked_at, external_ref, country, created_at, updated_at FROM customers ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var customers []*Customer
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.PaymentTermDays, &customer.ActivatedAt, &customer.Name, &customer.Number, &customer.ClosedAt, &customer.CreditLimit, &customer.Status, &customer.VATNumber, &customer.BlockedAt, &customer.ExternalRef, &customer.Country, &customer.CreatedAt, &customer.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}
		customers = append(customers, &customer)
	}
	return customers, rows.Err()
}

func (s *SQLStore) UpdateCustomer(ctx context.Context, customer *Customer) error {
	// update customer first so we fail early
	res, err := s.db.ExecContext(ctx, `UPDATE customers SET payment_term_days = $2, activated_at = $3, name = $4, number = $5, closed_at = $6, credit_limit = $7, status = $8, vat_number = $9, blocked_at = $10, external_ref = $11, country = $12, updated_at = $13 WHERE id = $1`, customer.ID, customer.PaymentTermDays, customer.ActivatedAt, customer.Name, customer.Number, customer.ClosedAt, customer.CreditLimit, customer.Status, customer.VATNumber, customer.BlockedAt, customer.ExternalRef, customer.Country, customer.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLStore) ListBillingAddresses(ctx context.Context, limit int, offset int) ([]*BillingAddress, error) {
	// we need this for the checks below: list billing addresses
	rows, err := s.db.QueryContext(ctx, `SELECT id, primary, country, customer_id, street, zip, attention, city, created_at, updated_at FROM billing_addresses ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list billing addresses: %w", err)
	}

	defer rows.Close()

	var addresses []*BillingAddress
	for rows.Next() {
		var address BillingAddress
		if err := rows.Scan(&address.ID, &address.Primary, &address.Country, &address.CustomerID, &address.Street, &address.Zip, &address.Attention, &address.City, &address.CreatedAt, &address.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan billing address: %w", err)
		}
		addresses = append(addresses, &address)
	}
	return addresses, rows.Err()
}

// ListCustomersByStatus implements Store.
func (s *SQLStore) ListCustomersByStatus(ctx context.Context, status CustomerStatus) ([]*Customer, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, payment_term_days, activated_at, name, number, closed_at, credit_limit, status, vat_number, blocked_at, external_ref, country, created_at, updated_at FROM customers WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var customers []*Customer
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.PaymentTermDays, &customer.ActivatedAt, &customer.Name, &customer.Number, &customer.ClosedAt, &customer.CreditLimit, &customer.Status, &customer.VATNumber, &customer.BlockedAt, &customer.ExternalRef, &customer.Country, &customer.CreatedAt, &customer.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}
		customers = append(customers, &customer)
	}
	return customers, rows.Err()
}

// ListBillingAddressesForCustomer implements Store.
func (s *SQLStore) ListBillingAddressesForCustomer(ctx context.Context, customerID int64) ([]*BillingAddress, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, primary, country, customer_id, street, zip, attention, city, created_at, updated_at FROM billing_addresses WHERE customer_id = $1 ORDER BY id`, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list billing addresses: %w", err)
	}
	defer rows.Close()

	var addresses []*BillingAddress
	for rows.Next() {
		var address BillingAddress
		if err := rows.Scan(&address.ID, &address.Primary, &address.Country, &address.CustomerID, &address.Street, &address.Zip, &address.Attention, &address.City, &address.CreatedAt, &address.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan billing address: %w", err)
		}
		addresses = append(addresses, &address)
	}
	return addresses, rows.Err()
}

func (s *SQLStore) FindCustomer(ctx context.Context, id int64) (*Customer, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, payment_term_days, activated_at, name, number, closed_at, credit_limit, status, vat_number, blocked_at, external_ref, country, created_at, updated_at FROM customers WHERE id = $1`, id)
	var customer Customer
	if err := row.Scan(&customer.ID, &customer.PaymentTermDays, &customer.ActivatedAt, &customer.Name, &customer.Number, &customer.ClosedAt, &customer.CreditLimit, &customer.Status, &customer.VATNumber, &customer.BlockedAt, &customer.ExternalRef, &customer.Country, &customer.CreatedAt, &customer.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query customer %d: %w", id, err)
	}
	return &customer, nil
}
