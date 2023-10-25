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

func (s *SQLStore) FindCustomer(ctx context.Context, id int64) (*Customer, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, external_ref, status, vat_number, credit_limit, activated_at, closed_at, payment_term_days, name, blocked_at, number, country, created_at, updated_at FROM customers WHERE id = $1`, id)
	var customer Customer
	if err := row.Scan(&customer.ID, &customer.ExternalRef, &customer.Status, &customer.VATNumber, &customer.CreditLimit, &customer.ActivatedAt, &customer.ClosedAt, &customer.PaymentTermDays, &customer.Name, &customer.BlockedAt, &customer.Number, &customer.Country, &customer.CreatedAt, &customer.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query customer %d: %w", id, err)
	}
	return &customer, nil
}

func (s *SQLStore) ListCustomers(ctx context.Context, limit int, offset int) ([]*Customer, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, external_ref, status, vat_number, credit_limit, activated_at, closed_at, payment_term_days, name, blocked_at, number, country, created_at, updated_at FROM customers ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var customers []*Customer
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.ExternalRef, &customer.Status, &customer.VATNumber, &customer.CreditLimit, &customer.ActivatedAt, &customer.ClosedAt, &customer.PaymentTermDays, &customer.Name, &customer.BlockedAt, &customer.Number, &customer.Country, &customer.CreatedAt, &customer.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan customer: %w", err)
		}
		customers = append(customers, &customer)
	}
	return customers, rows.Err()
}

// ListCustomersByStatus implements Store.
func (s *SQLStore) ListCustomersByStatus(ctx context.Context, status CustomerStatus) ([]*Customer, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, external_ref, status, vat_number, credit_limit, activated_at, closed_at, payment_term_days, name, blocked_at, number, country, created_at, updated_at FROM customers WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var customers []*Customer
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.ExternalRef, &customer.Status, &customer.VATNumber, &customer.CreditLimit, &customer.ActivatedAt, &customer.ClosedAt, &customer.PaymentTermDays, &customer.Name, &customer.BlockedAt, &customer.Number, &customer.Country, &customer.CreatedAt, &customer.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan customer: %w", err)
		}
		customers = append(customers, &customer)
	}
	return customers, rows.Err()
}

func (s *SQLStore) UpdateBillingAddress(ctx context.Context, address *BillingAddress) error {
	res, err := s.db.ExecContext(ctx, `UPDATE billing_addresses SET customer_id = $2, street = $3, zip = $4, city = $5, country = $6, primary = $7, attention = $8, updated_at = $9 WHERE id = $1`, address.ID, address.CustomerID, address.Street, address.Zip, address.City, address.Country, address.Primary, address.Attention, address.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update billing address: %w", err)
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

func (s *SQLStore) FindBillingAddress(ctx context.Context, id int64) (*BillingAddress, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, customer_id, street, zip, city, country, primary, attention, created_at, updated_at FROM billing_addresses WHERE id = $1`, id)
	var address BillingAddress
	if err := row.Scan(&address.ID, &address.CustomerID, &address.Street, &address.Zip, &address.City, &address.Country, &address.Primary, &address.Attention, &address.CreatedAt, &address.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query billing address %d: %w", id, err)
	}
	return &address, nil
}

func (s *SQLStore) UpdateCustomer(ctx context.Context, customer *Customer) error {
	// update customer first so we fail early
	res, err := s.db.ExecContext(ctx, `UPDATE customers SET external_ref = $2, status = $3, vat_number = $4, credit_limit = $5, activated_at = $6, closed_at = $7, payment_term_days = $8, name = $9, blocked_at = $10, number = $11, country = $12, updated_at = $13 WHERE id = $1`, customer.ID, customer.ExternalRef, customer.Status, customer.VATNumber, customer.CreditLimit, customer.ActivatedAt, customer.ClosedAt, customer.PaymentTermDays, customer.Name, customer.BlockedAt, customer.Number, customer.Country, customer.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update customer: %w", err)
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
	rows, err := s.db.QueryContext(ctx, `SELECT id, customer_id, street, zip, city, country, primary, attention, created_at, updated_at FROM billing_addresses ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list billing addresses: %w", err)
	}

	defer rows.Close()

	var addresses []*BillingAddress
	for rows.Next() {
		var address BillingAddress
		if err := rows.Scan(&address.ID, &address.CustomerID, &address.Street, &address.Zip, &address.City, &address.Country, &address.Primary, &address.Attention, &address.CreatedAt, &address.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan billing address: %w", err)
		}
		addresses = append(addresses, &address)
	}
	return addresses, rows.Err()
}

// ListBillingAddressesForCustomer implements Store.
func (s *SQLStore) ListBillingAddressesForCustomer(ctx context.Context, customerID int64) ([]*BillingAddress, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, customer_id, street, zip, city, country, primary, attention, created_at, updated_at FROM billing_addresses WHERE customer_id = $1 ORDER BY id`, customerID)
	if err != nil {
		return nil, fmt.Errorf("list billing addresses: %w", err)
	}
	defer rows.Close()

	var addresses []*BillingAddress
	for rows.Next() {
		var address BillingAddress
		if err := rows.Scan(&address.ID, &address.CustomerID, &address.Street, &address.Zip, &address.City, &address.Country, &address.Primary, &address.Attention, &address.CreatedAt, &address.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan billing address: %w", err)
		}
		addresses = append(addresses, &address)
	}
	return addresses, rows.Err()
}

// InsertBillingAddress implements Store.
func (s *SQLStore) InsertBillingAddress(ctx context.Context, address *BillingAddress) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO billing_addresses (customer_id, street, zip, city, country, primary, attention, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`, address.CustomerID, address.Street, address.Zip, address.City, address.Country, address.Primary, address.Attention, address.CreatedAt, address.UpdatedAt).Scan(&address.ID); err != nil {
		return fmt.Errorf("insert billing address: %w", err)
	}
	return nil
}

// InsertCustomer implements Store.
func (s *SQLStore) InsertCustomer(ctx context.Context, customer *Customer) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO customers (external_ref, status, vat_number, credit_limit, activated_at, closed_at, payment_term_days, name, blocked_at, number, country, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`, customer.ExternalRef, customer.Status, customer.VATNumber, customer.CreditLimit, customer.ActivatedAt, customer.ClosedAt, customer.PaymentTermDays, customer.Name, customer.BlockedAt, customer.Number, customer.Country, customer.CreatedAt, customer.UpdatedAt).Scan(&customer.ID); err != nil {
		return fmt.Errorf("insert customer: %w", err)
	}

	return nil
}
