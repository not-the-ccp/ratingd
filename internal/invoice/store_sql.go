package invoice

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

func (s *SQLStore) FindInvoice(ctx context.Context, id int64) (*Invoice, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, customer_id, number, currency, issued_at, due_at, net_amount, tax_amount, status, issued_on, paid_at, voided_at, created_at, updated_at FROM invoices WHERE id = $1`, id)
	var invoice Invoice
	if err := row.Scan(&invoice.ID, &invoice.CustomerID, &invoice.Number, &invoice.Currency, &invoice.IssuedAt, &invoice.DueAt, &invoice.NetAmount, &invoice.TaxAmount, &invoice.Status, &invoice.IssuedOn, &invoice.PaidAt, &invoice.VoidedAt, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query invoice %d: %w", id, err)
	}
	return &invoice, nil
}

func (s *SQLStore) ListInvoices(ctx context.Context, limit int, offset int) ([]*Invoice, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, customer_id, number, currency, issued_at, due_at, net_amount, tax_amount, status, issued_on, paid_at, voided_at, created_at, updated_at FROM invoices ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list invoices: %w", err)
	}
	defer rows.Close()

	var invoices []*Invoice
	for rows.Next() {
		var invoice Invoice
		if err := rows.Scan(&invoice.ID, &invoice.CustomerID, &invoice.Number, &invoice.Currency, &invoice.IssuedAt, &invoice.DueAt, &invoice.NetAmount, &invoice.TaxAmount, &invoice.Status, &invoice.IssuedOn, &invoice.PaidAt, &invoice.VoidedAt, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan invoice: %w", err)
		}
		invoices = append(invoices, &invoice)
	}
	return invoices, rows.Err()
}

func (s *SQLStore) ListInvoicesByStatus(ctx context.Context, status InvoiceStatus) ([]*Invoice, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, customer_id, number, currency, issued_at, due_at, net_amount, tax_amount, status, issued_on, paid_at, voided_at, created_at, updated_at FROM invoices WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("list invoices: %w", err)
	}
	defer rows.Close()

	var invoices []*Invoice
	for rows.Next() {
		var invoice Invoice
		if err := rows.Scan(&invoice.ID, &invoice.CustomerID, &invoice.Number, &invoice.Currency, &invoice.IssuedAt, &invoice.DueAt, &invoice.NetAmount, &invoice.TaxAmount, &invoice.Status, &invoice.IssuedOn, &invoice.PaidAt, &invoice.VoidedAt, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan invoice: %w", err)
		}
		invoices = append(invoices, &invoice)
	}
	return invoices, rows.Err()
}

func (s *SQLStore) InsertInvoice(ctx context.Context, invoice *Invoice) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO invoices (customer_id, number, currency, issued_at, due_at, net_amount, tax_amount, status, issued_on, paid_at, voided_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`, invoice.CustomerID, invoice.Number, invoice.Currency, invoice.IssuedAt, invoice.DueAt, invoice.NetAmount, invoice.TaxAmount, invoice.Status, invoice.IssuedOn, invoice.PaidAt, invoice.VoidedAt, invoice.CreatedAt, invoice.UpdatedAt).Scan(&invoice.ID); err != nil {
		return fmt.Errorf("insert invoice: %w", err)
	}

	return nil
}

func (s *SQLStore) UpdateInvoice(ctx context.Context, invoice *Invoice) error {
	res, err := s.db.ExecContext(ctx, `UPDATE invoices SET customer_id = $2, number = $3, currency = $4, issued_at = $5, due_at = $6, net_amount = $7, tax_amount = $8, status = $9, issued_on = $10, paid_at = $11, voided_at = $12, updated_at = $13 WHERE id = $1`, invoice.ID, invoice.CustomerID, invoice.Number, invoice.Currency, invoice.IssuedAt, invoice.DueAt, invoice.NetAmount, invoice.TaxAmount, invoice.Status, invoice.IssuedOn, invoice.PaidAt, invoice.VoidedAt, invoice.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update invoice: %w", err)
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
