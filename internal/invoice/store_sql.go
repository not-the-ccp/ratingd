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

func (s *SQLStore) ListInvoices(ctx context.Context, limit int, offset int) ([]*Invoice, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, number, customer_id, paid_at, issued_on, purchase_order, net_amount, issued_at, currency, exchange_rate, voided_at, tax_amount, status, due_at, created_at, updated_at FROM invoices ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not list invoices: %w", err)
	}
	defer rows.Close()

	var invoices []*Invoice
	for rows.Next() {
		var invoice Invoice
		if err := rows.Scan(&invoice.ID, &invoice.Number, &invoice.CustomerID, &invoice.PaidAt, &invoice.IssuedOn, &invoice.PurchaseOrder, &invoice.NetAmount, &invoice.IssuedAt, &invoice.Currency, &invoice.ExchangeRate, &invoice.VoidedAt, &invoice.TaxAmount, &invoice.Status, &invoice.DueAt, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan invoice: %w", err)
		}
		invoices = append(invoices, &invoice)
	}
	return invoices, rows.Err()
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) ListInvoicesByStatus(ctx context.Context, status InvoiceStatus) ([]*Invoice, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, number, customer_id, paid_at, issued_on, purchase_order, net_amount, issued_at, currency, exchange_rate, voided_at, tax_amount, status, due_at, created_at, updated_at FROM invoices WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("could not list invoices: %w", err)
	}
	defer rows.Close()

	var invoices []*Invoice
	for rows.Next() {
		var invoice Invoice
		if err := rows.Scan(&invoice.ID, &invoice.Number, &invoice.CustomerID, &invoice.PaidAt, &invoice.IssuedOn, &invoice.PurchaseOrder, &invoice.NetAmount, &invoice.IssuedAt, &invoice.Currency, &invoice.ExchangeRate, &invoice.VoidedAt, &invoice.TaxAmount, &invoice.Status, &invoice.DueAt, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan invoice: %w", err)
		}
		invoices = append(invoices, &invoice)
	}
	return invoices, rows.Err()
}

func (s *SQLStore) ListCreditNotes(ctx context.Context, limit int, offset int) ([]*CreditNote, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, reason, amount, status, number, applied_at, invoice_id, created_at, updated_at FROM credit_notes ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not list credit notes: %w", err)
	}
	defer rows.Close()
	var notes []*CreditNote
	for rows.Next() {
		var note CreditNote
		if err := rows.Scan(&note.ID, &note.Reason, &note.Amount, &note.Status, &note.Number, &note.AppliedAt, &note.InvoiceID, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan credit note: %w", err)
		}
		notes = append(notes, &note)
	}
	return notes, rows.Err()
}

func (s *SQLStore) UpdateCreditNote(ctx context.Context, note *CreditNote) error {
	res, err := s.db.ExecContext(ctx, `UPDATE credit_notes SET reason = $2, amount = $3, status = $4, number = $5, applied_at = $6, invoice_id = $7, updated_at = $8 WHERE id = $1`, note.ID, note.Reason, note.Amount, note.Status, note.Number, note.AppliedAt, note.InvoiceID, note.UpdatedAt)
	if err != nil {
		return fmt.Errorf("could not update credit note: %w", err)
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

func (s *SQLStore) ListCreditNotesByStatus(ctx context.Context, status CreditNoteStatus) ([]*CreditNote, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, reason, amount, status, number, applied_at, invoice_id, created_at, updated_at FROM credit_notes WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("could not list credit notes: %w", err)
	}

	defer rows.Close()

	var notes []*CreditNote
	for rows.Next() {
		var note CreditNote
		if err := rows.Scan(&note.ID, &note.Reason, &note.Amount, &note.Status, &note.Number, &note.AppliedAt, &note.InvoiceID, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan credit note: %w", err)
		}
		notes = append(notes, &note)
	}
	return notes, rows.Err()
}

func (s *SQLStore) FetchCreditNote(ctx context.Context, id int64) (*CreditNote, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, reason, amount, status, number, applied_at, invoice_id, created_at, updated_at FROM credit_notes WHERE id = $1`, id)
	var creditNote CreditNote
	if err := row.Scan(&creditNote.ID, &creditNote.Reason, &creditNote.Amount, &creditNote.Status, &creditNote.Number, &creditNote.AppliedAt, &creditNote.InvoiceID, &creditNote.CreatedAt, &creditNote.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("could not query credit note %d: %w", id, err)
	}
	return &creditNote, nil
}

func (s *SQLStore) ListInvoiceLines(ctx context.Context, limit int, offset int) ([]*InvoiceLine, error) {
	// list invoice lines first so we fail early
	rows, err := s.db.QueryContext(ctx, `SELECT id, unit_price, description, amount, quantity, discount_amount, shipment_id, invoice_id, created_at, updated_at FROM invoice_lines ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not list invoice lines: %w", err)
	}

	defer rows.Close()

	var lines []*InvoiceLine
	for rows.Next() {
		var invoiceLine InvoiceLine
		if err := rows.Scan(&invoiceLine.ID, &invoiceLine.UnitPrice, &invoiceLine.Label, &invoiceLine.Amount, &invoiceLine.Quantity, &invoiceLine.DiscountAmount, &invoiceLine.ShipmentID, &invoiceLine.InvoiceID, &invoiceLine.CreatedAt, &invoiceLine.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan invoice line: %w", err)
		}
		lines = append(lines, &invoiceLine)
	}
	return lines, rows.Err()
}

func (s *SQLStore) ListInvoiceLinesForInvoice(ctx context.Context, invoiceID int64) ([]*InvoiceLine, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, unit_price, description, amount, quantity, discount_amount, shipment_id, invoice_id, created_at, updated_at FROM invoice_lines WHERE invoice_id = $1 ORDER BY id`, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("could not list invoice lines: %w", err)
	}

	defer rows.Close()

	var lines []*InvoiceLine
	for rows.Next() {
		var line InvoiceLine
		if err := rows.Scan(&line.ID, &line.UnitPrice, &line.Label, &line.Amount, &line.Quantity, &line.DiscountAmount, &line.ShipmentID, &line.InvoiceID, &line.CreatedAt, &line.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan invoice line: %w", err)
		}
		lines = append(lines, &line)
	}
	return lines, rows.Err()
}

func (s *SQLStore) ListCreditNotesForInvoice(ctx context.Context, invoiceID int64) ([]*CreditNote, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, reason, amount, status, number, applied_at, invoice_id, created_at, updated_at FROM credit_notes WHERE invoice_id = $1 ORDER BY id`, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("could not list credit notes: %w", err)
	}

	defer rows.Close()

	var creditNotes []*CreditNote
	for rows.Next() {
		var note CreditNote
		if err := rows.Scan(&note.ID, &note.Reason, &note.Amount, &note.Status, &note.Number, &note.AppliedAt, &note.InvoiceID, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan credit note: %w", err)
		}
		creditNotes = append(creditNotes, &note)
	}
	return creditNotes, rows.Err()
}

func (s *SQLStore) FetchInvoice(ctx context.Context, id int64) (*Invoice, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, number, customer_id, paid_at, issued_on, purchase_order, net_amount, issued_at, currency, exchange_rate, voided_at, tax_amount, status, due_at, created_at, updated_at FROM invoices WHERE id = $1`, id)
	var invoice Invoice
	if err := row.Scan(&invoice.ID, &invoice.Number, &invoice.CustomerID, &invoice.PaidAt, &invoice.IssuedOn, &invoice.PurchaseOrder, &invoice.NetAmount, &invoice.IssuedAt, &invoice.Currency, &invoice.ExchangeRate, &invoice.VoidedAt, &invoice.TaxAmount, &invoice.Status, &invoice.DueAt, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("could not query invoice %d: %w", id, err)
	}
	return &invoice, nil
}

func (s *SQLStore) InsertInvoiceLine(ctx context.Context, line *InvoiceLine) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO invoice_lines (unit_price, description, amount, quantity, discount_amount, shipment_id, invoice_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`, line.UnitPrice, line.Label, line.Amount, line.Quantity, line.DiscountAmount, line.ShipmentID, line.InvoiceID, line.CreatedAt, line.UpdatedAt).Scan(&line.ID); err != nil {
		return fmt.Errorf("could not insert invoice line: %w", err)
	}

	return nil
}

func (s *SQLStore) InsertCreditNote(ctx context.Context, note *CreditNote) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO credit_notes (reason, amount, status, number, applied_at, invoice_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, note.Reason, note.Amount, note.Status, note.Number, note.AppliedAt, note.InvoiceID, note.CreatedAt, note.UpdatedAt).Scan(&note.ID); err != nil {
		return fmt.Errorf("could not insert credit note: %w", err)
	}

	return nil
}

func (s *SQLStore) UpdateInvoiceLine(ctx context.Context, line *InvoiceLine) error {
	res, err := s.db.ExecContext(ctx, `UPDATE invoice_lines SET unit_price = $2, description = $3, amount = $4, quantity = $5, discount_amount = $6, shipment_id = $7, invoice_id = $8, updated_at = $9 WHERE id = $1`, line.ID, line.UnitPrice, line.Label, line.Amount, line.Quantity, line.DiscountAmount, line.ShipmentID, line.InvoiceID, line.UpdatedAt)
	if err != nil {
		return fmt.Errorf("could not update invoice line: %w", err)
	}

	// we need this for the checks below: rows affected
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// GetInvoiceLine implements Store.
func (s *SQLStore) FetchInvoiceLine(ctx context.Context, id int64) (*InvoiceLine, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, unit_price, description, amount, quantity, discount_amount, shipment_id, invoice_id, created_at, updated_at FROM invoice_lines WHERE id = $1`, id)
	var line InvoiceLine
	if err := row.Scan(&line.ID, &line.UnitPrice, &line.Label, &line.Amount, &line.Quantity, &line.DiscountAmount, &line.ShipmentID, &line.InvoiceID, &line.CreatedAt, &line.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("could not query invoice line %d: %w", id, err)
	}
	return &line, nil
}

func (s *SQLStore) InsertInvoice(ctx context.Context, invoice *Invoice) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO invoices (number, customer_id, paid_at, issued_on, purchase_order, net_amount, issued_at, currency, exchange_rate, voided_at, tax_amount, status, due_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id`, invoice.Number, invoice.CustomerID, invoice.PaidAt, invoice.IssuedOn, invoice.PurchaseOrder, invoice.NetAmount, invoice.IssuedAt, invoice.Currency, invoice.ExchangeRate, invoice.VoidedAt, invoice.TaxAmount, invoice.Status, invoice.DueAt, invoice.CreatedAt, invoice.UpdatedAt).Scan(&invoice.ID); err != nil {
		return fmt.Errorf("could not insert invoice: %w", err)
	}

	return nil
}

func (s *SQLStore) UpdateInvoice(ctx context.Context, invoice *Invoice) error {
	res, err := s.db.ExecContext(ctx, `UPDATE invoices SET number = $2, customer_id = $3, paid_at = $4, issued_on = $5, purchase_order = $6, net_amount = $7, issued_at = $8, currency = $9, exchange_rate = $10, voided_at = $11, tax_amount = $12, status = $13, due_at = $14, updated_at = $15 WHERE id = $1`, invoice.ID, invoice.Number, invoice.CustomerID, invoice.PaidAt, invoice.IssuedOn, invoice.PurchaseOrder, invoice.NetAmount, invoice.IssuedAt, invoice.Currency, invoice.ExchangeRate, invoice.VoidedAt, invoice.TaxAmount, invoice.Status, invoice.DueAt, invoice.UpdatedAt)
	if err != nil {
		return fmt.Errorf("could not update invoice: %w", err)
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
