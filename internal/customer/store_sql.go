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

func (s *SQLStore) ListCustomers(ctx context.Context, limit int, offset int) ([]*Customer, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, activated_at, external_ref, vat_number, name, closed_at, blocked_at, number, status, country, iban, payment_term_days, credit_limit, billing_email, created_at, updated_at FROM customers ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var customers []*Customer
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.ActivatedAt, &customer.ExternalRef, &customer.VATNumber, &customer.DisplayName, &customer.ClosedAt, &customer.BlockedAt, &customer.DocumentNumber, &customer.Status, &customer.Country, &customer.IBAN, &customer.PaymentTermDays, &customer.CreditLimit, &customer.BillingEmail, &customer.CreatedAt, &customer.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}
		customers = append(customers, &customer)
	}
	return customers, rows.Err()
}

func (s *SQLStore) LoadBillingAddress(ctx context.Context, id int64) (*BillingAddress, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, primary, street, customer_id, region, attention, country, zip, city, created_at, updated_at FROM billing_addresses WHERE id = $1`, id)
	var address BillingAddress
	if err := row.Scan(&address.ID, &address.Primary, &address.Street, &address.CustomerID, &address.Region, &address.Attention, &address.Country, &address.Zip, &address.City, &address.CreatedAt, &address.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query billing address %d: %w", id, err)
	}
	return &address, nil
}

// InsertBillingAddress implements Store.
func (s *SQLStore) InsertBillingAddress(ctx context.Context, address *BillingAddress) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO billing_addresses (primary, street, customer_id, region, attention, country, zip, city, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`, address.Primary, address.Street, address.CustomerID, address.Region, address.Attention, address.Country, address.Zip, address.City, address.CreatedAt, address.UpdatedAt).Scan(&address.ID); err != nil {
		return fmt.Errorf("failed to insert billing address: %w", err)
	}
	return nil
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

// InsertCustomer implements Store.
func (s *SQLStore) InsertCustomer(ctx context.Context, customer *Customer) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO customers (activated_at, external_ref, vat_number, name, closed_at, blocked_at, number, status, country, iban, payment_term_days, credit_limit, billing_email, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id`, customer.ActivatedAt, customer.ExternalRef, customer.VATNumber, customer.DisplayName, customer.ClosedAt, customer.BlockedAt, customer.DocumentNumber, customer.Status, customer.Country, customer.IBAN, customer.PaymentTermDays, customer.CreditLimit, customer.BillingEmail, customer.CreatedAt, customer.UpdatedAt).Scan(&customer.ID); err != nil {
		return fmt.Errorf("failed to insert customer: %w", err)
	}

	return nil
}

func (s *SQLStore) LoadCustomer(ctx context.Context, id int64) (*Customer, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, activated_at, external_ref, vat_number, name, closed_at, blocked_at, number, status, country, iban, payment_term_days, credit_limit, billing_email, created_at, updated_at FROM customers WHERE id = $1`, id)
	var customer Customer
	if err := row.Scan(&customer.ID, &customer.ActivatedAt, &customer.ExternalRef, &customer.VATNumber, &customer.DisplayName, &customer.ClosedAt, &customer.BlockedAt, &customer.DocumentNumber, &customer.Status, &customer.Country, &customer.IBAN, &customer.PaymentTermDays, &customer.CreditLimit, &customer.BillingEmail, &customer.CreatedAt, &customer.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query customer %d: %w", id, err)
	}
	return &customer, nil
}

// ListCustomersByStatus implements Store.
func (s *SQLStore) ListCustomersByStatus(ctx context.Context, status CustomerStatus) ([]*Customer, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, activated_at, external_ref, vat_number, name, closed_at, blocked_at, number, status, country, iban, payment_term_days, credit_limit, billing_email, created_at, updated_at FROM customers WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var customers []*Customer
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.ActivatedAt, &customer.ExternalRef, &customer.VATNumber, &customer.DisplayName, &customer.ClosedAt, &customer.BlockedAt, &customer.DocumentNumber, &customer.Status, &customer.Country, &customer.IBAN, &customer.PaymentTermDays, &customer.CreditLimit, &customer.BillingEmail, &customer.CreatedAt, &customer.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}
		customers = append(customers, &customer)
	}
	return customers, rows.Err()
}

func (s *SQLStore) ListBillingAddresses(ctx context.Context, limit int, offset int) ([]*BillingAddress, error) {
	// we need this for the checks below: list billing addresses
	rows, err := s.db.QueryContext(ctx, `SELECT id, primary, street, customer_id, region, attention, country, zip, city, created_at, updated_at FROM billing_addresses ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list billing addresses: %w", err)
	}

	defer rows.Close()

	var addresses []*BillingAddress
	for rows.Next() {
		var address BillingAddress
		if err := rows.Scan(&address.ID, &address.Primary, &address.Street, &address.CustomerID, &address.Region, &address.Attention, &address.Country, &address.Zip, &address.City, &address.CreatedAt, &address.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan billing address: %w", err)
		}
		addresses = append(addresses, &address)
	}
	return addresses, rows.Err()
}

func (s *SQLStore) UpdateBillingAddress(ctx context.Context, address *BillingAddress) error {
	res, err := s.db.ExecContext(ctx, `UPDATE billing_addresses SET primary = $2, street = $3, customer_id = $4, region = $5, attention = $6, country = $7, zip = $8, city = $9, updated_at = $10 WHERE id = $1`, address.ID, address.Primary, address.Street, address.CustomerID, address.Region, address.Attention, address.Country, address.Zip, address.City, address.UpdatedAt)
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

func (s *SQLStore) UpdateCustomer(ctx context.Context, customer *Customer) error {
	// update customer first so we fail early
	res, err := s.db.ExecContext(ctx, `UPDATE customers SET activated_at = $2, external_ref = $3, vat_number = $4, name = $5, closed_at = $6, blocked_at = $7, number = $8, status = $9, country = $10, iban = $11, payment_term_days = $12, credit_limit = $13, billing_email = $14, updated_at = $15 WHERE id = $1`, customer.ID, customer.ActivatedAt, customer.ExternalRef, customer.VATNumber, customer.DisplayName, customer.ClosedAt, customer.BlockedAt, customer.DocumentNumber, customer.Status, customer.Country, customer.IBAN, customer.PaymentTermDays, customer.CreditLimit, customer.BillingEmail, customer.UpdatedAt)
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

// ListBillingAddressesForCustomer implements Store.
func (s *SQLStore) ListBillingAddressesForCustomer(ctx context.Context, customerID int64) ([]*BillingAddress, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, primary, street, customer_id, region, attention, country, zip, city, created_at, updated_at FROM billing_addresses WHERE customer_id = $1 ORDER BY id`, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list billing addresses: %w", err)
	}
	defer rows.Close()

	var addresses []*BillingAddress
	for rows.Next() {
		var address BillingAddress
		if err := rows.Scan(&address.ID, &address.Primary, &address.Street, &address.CustomerID, &address.Region, &address.Attention, &address.Country, &address.Zip, &address.City, &address.CreatedAt, &address.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan billing address: %w", err)
		}
		addresses = append(addresses, &address)
	}
	return addresses, rows.Err()
}
