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
	row := s.db.QueryRowContext(ctx, `SELECT id, name, number, vat_number, country, payment_term_days, credit_limit, status, activated_at, blocked_at, closed_at, created_at, updated_at FROM customers WHERE id = $1`, id)
	var customer Customer
	if err := row.Scan(&customer.ID, &customer.Name, &customer.Number, &customer.VATNumber, &customer.Country, &customer.PaymentTermDays, &customer.CreditLimit, &customer.Status, &customer.ActivatedAt, &customer.BlockedAt, &customer.ClosedAt, &customer.CreatedAt, &customer.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query customer %d: %w", id, err)
	}
	return &customer, nil
}

func (s *SQLStore) ListCustomers(ctx context.Context, limit int, offset int) ([]*Customer, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, number, vat_number, country, payment_term_days, credit_limit, status, activated_at, blocked_at, closed_at, created_at, updated_at FROM customers ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var customers []*Customer
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Number, &customer.VATNumber, &customer.Country, &customer.PaymentTermDays, &customer.CreditLimit, &customer.Status, &customer.ActivatedAt, &customer.BlockedAt, &customer.ClosedAt, &customer.CreatedAt, &customer.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan customer: %w", err)
		}
		customers = append(customers, &customer)
	}
	return customers, rows.Err()
}

// ListCustomersByStatus implements Store.
func (s *SQLStore) ListCustomersByStatus(ctx context.Context, status CustomerStatus) ([]*Customer, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, number, vat_number, country, payment_term_days, credit_limit, status, activated_at, blocked_at, closed_at, created_at, updated_at FROM customers WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var customers []*Customer
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Number, &customer.VATNumber, &customer.Country, &customer.PaymentTermDays, &customer.CreditLimit, &customer.Status, &customer.ActivatedAt, &customer.BlockedAt, &customer.ClosedAt, &customer.CreatedAt, &customer.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan customer: %w", err)
		}
		customers = append(customers, &customer)
	}
	return customers, rows.Err()
}

// InsertCustomer implements Store.
func (s *SQLStore) InsertCustomer(ctx context.Context, customer *Customer) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO customers (name, number, vat_number, country, payment_term_days, credit_limit, status, activated_at, blocked_at, closed_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`, customer.Name, customer.Number, customer.VATNumber, customer.Country, customer.PaymentTermDays, customer.CreditLimit, customer.Status, customer.ActivatedAt, customer.BlockedAt, customer.ClosedAt, customer.CreatedAt, customer.UpdatedAt).Scan(&customer.ID); err != nil {
		return fmt.Errorf("insert customer: %w", err)
	}

	return nil
}

func (s *SQLStore) UpdateCustomer(ctx context.Context, customer *Customer) error {
	// update customer first so we fail early
	res, err := s.db.ExecContext(ctx, `UPDATE customers SET name = $2, number = $3, vat_number = $4, country = $5, payment_term_days = $6, credit_limit = $7, status = $8, activated_at = $9, blocked_at = $10, closed_at = $11, updated_at = $12 WHERE id = $1`, customer.ID, customer.Name, customer.Number, customer.VATNumber, customer.Country, customer.PaymentTermDays, customer.CreditLimit, customer.Status, customer.ActivatedAt, customer.BlockedAt, customer.ClosedAt, customer.UpdatedAt)
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
