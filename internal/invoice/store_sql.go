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
	rows, err := s.db.QueryContext(ctx, `SELECT id, status, issued_at, paid_at, currency, customer_id, net_amount, number, due_at, issued_on, tax_amount, voided_at, created_at, updated_at FROM invoices ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not list invoices: %w", err)
	}
	defer rows.Close()

	var invoices []*Invoice
	for rows.Next() {
		var invoice Invoice
		if err := rows.Scan(&invoice.ID, &invoice.Status, &invoice.IssuedAt, &invoice.PaidAt, &invoice.Currency, &invoice.CustomerID, &invoice.NetAmount, &invoice.Number, &invoice.DueAt, &invoice.IssuedOn, &invoice.TaxAmount, &invoice.VoidedAt, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
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
	rows, err := s.db.QueryContext(ctx, `SELECT id, status, issued_at, paid_at, currency, customer_id, net_amount, number, due_at, issued_on, tax_amount, voided_at, created_at, updated_at FROM invoices WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("could not list invoices: %w", err)
	}
	defer rows.Close()

	var invoices []*Invoice
	for rows.Next() {
		var invoice Invoice
		if err := rows.Scan(&invoice.ID, &invoice.Status, &invoice.IssuedAt, &invoice.PaidAt, &invoice.Currency, &invoice.CustomerID, &invoice.NetAmount, &invoice.Number, &invoice.DueAt, &invoice.IssuedOn, &invoice.TaxAmount, &invoice.VoidedAt, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan invoice: %w", err)
		}
		invoices = append(invoices, &invoice)
	}
	return invoices, rows.Err()
}

func (s *SQLStore) ListCreditNotes(ctx context.Context, limit int, offset int) ([]*CreditNote, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, applied_at, invoice_id, amount, number, reason, status, created_at, updated_at FROM credit_notes ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not list credit notes: %w", err)
	}
	defer rows.Close()
	var notes []*CreditNote
	for rows.Next() {
		var note CreditNote
		if err := rows.Scan(&note.ID, &note.AppliedAt, &note.InvoiceID, &note.Amount, &note.Number, &note.Reason, &note.Status, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan credit note: %w", err)
		}
		notes = append(notes, &note)
	}
	return notes, rows.Err()
}

func (s *SQLStore) InsertInvoice(ctx context.Context, invoice *Invoice) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO invoices (status, issued_at, paid_at, currency, customer_id, net_amount, number, due_at, issued_on, tax_amount, voided_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`, invoice.Status, invoice.IssuedAt, invoice.PaidAt, invoice.Currency, invoice.CustomerID, invoice.NetAmount, invoice.Number, invoice.DueAt, invoice.IssuedOn, invoice.TaxAmount, invoice.VoidedAt, invoice.CreatedAt, invoice.UpdatedAt).Scan(&invoice.ID); err != nil {
		return fmt.Errorf("could not insert invoice: %w", err)
	}

	return nil
}

func (s *SQLStore) ListCreditNotesByStatus(ctx context.Context, status CreditNoteStatus) ([]*CreditNote, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, applied_at, invoice_id, amount, number, reason, status, created_at, updated_at FROM credit_notes WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("could not list credit notes: %w", err)
	}

	defer rows.Close()

	var notes []*CreditNote
	for rows.Next() {
		var note CreditNote
		if err := rows.Scan(&note.ID, &note.AppliedAt, &note.InvoiceID, &note.Amount, &note.Number, &note.Reason, &note.Status, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan credit note: %w", err)
		}
		notes = append(notes, &note)
	}
	return notes, rows.Err()
}

func (s *SQLStore) LoadCreditNote(ctx context.Context, id int64) (*CreditNote, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, applied_at, invoice_id, amount, number, reason, status, created_at, updated_at FROM credit_notes WHERE id = $1`, id)
	var creditNote CreditNote
	if err := row.Scan(&creditNote.ID, &creditNote.AppliedAt, &creditNote.InvoiceID, &creditNote.Amount, &creditNote.Number, &creditNote.Reason, &creditNote.Status, &creditNote.CreatedAt, &creditNote.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("could not query credit note %d: %w", id, err)
	}
	return &creditNote, nil
}

func (s *SQLStore) ListInvoiceLines(ctx context.Context, limit int, offset int) ([]*InvoiceLine, error) {
	// list invoice lines first so we fail early
	rows, err := s.db.QueryContext(ctx, `SELECT id, unit_price, amount, invoice_id, discount_amount, quantity, shipment_id, description, created_at, updated_at FROM invoice_lines ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not list invoice lines: %w", err)
	}

	defer rows.Close()

	var lines []*InvoiceLine
	for rows.Next() {
		var invoiceLine InvoiceLine
		if err := rows.Scan(&invoiceLine.ID, &invoiceLine.UnitPrice, &invoiceLine.Amount, &invoiceLine.InvoiceID, &invoiceLine.DiscountAmount, &invoiceLine.Quantity, &invoiceLine.ShipmentID, &invoiceLine.Label, &invoiceLine.CreatedAt, &invoiceLine.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan invoice line: %w", err)
		}
		lines = append(lines, &invoiceLine)
	}
	return lines, rows.Err()
}

func (s *SQLStore) ListInvoiceLinesForInvoice(ctx context.Context, invoiceID int64) ([]*InvoiceLine, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, unit_price, amount, invoice_id, discount_amount, quantity, shipment_id, description, created_at, updated_at FROM invoice_lines WHERE invoice_id = $1 ORDER BY id`, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("could not list invoice lines: %w", err)
	}

	defer rows.Close()

	var lines []*InvoiceLine
	for rows.Next() {
		var line InvoiceLine
		if err := rows.Scan(&line.ID, &line.UnitPrice, &line.Amount, &line.InvoiceID, &line.DiscountAmount, &line.Quantity, &line.ShipmentID, &line.Label, &line.CreatedAt, &line.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan invoice line: %w", err)
		}
		lines = append(lines, &line)
	}
	return lines, rows.Err()
}

func (s *SQLStore) ListCreditNotesForInvoice(ctx context.Context, invoiceID int64) ([]*CreditNote, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, applied_at, invoice_id, amount, number, reason, status, created_at, updated_at FROM credit_notes WHERE invoice_id = $1 ORDER BY id`, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("could not list credit notes: %w", err)
	}

	defer rows.Close()

	var creditNotes []*CreditNote
	for rows.Next() {
		var note CreditNote
		if err := rows.Scan(&note.ID, &note.AppliedAt, &note.InvoiceID, &note.Amount, &note.Number, &note.Reason, &note.Status, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan credit note: %w", err)
		}
		creditNotes = append(creditNotes, &note)
	}
	return creditNotes, rows.Err()
}

func (s *SQLStore) LoadInvoice(ctx context.Context, id int64) (*Invoice, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, status, issued_at, paid_at, currency, customer_id, net_amount, number, due_at, issued_on, tax_amount, voided_at, created_at, updated_at FROM invoices WHERE id = $1`, id)
	var invoice Invoice
	if err := row.Scan(&invoice.ID, &invoice.Status, &invoice.IssuedAt, &invoice.PaidAt, &invoice.Currency, &invoice.CustomerID, &invoice.NetAmount, &invoice.Number, &invoice.DueAt, &invoice.IssuedOn, &invoice.TaxAmount, &invoice.VoidedAt, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("could not query invoice %d: %w", id, err)
	}
	return &invoice, nil
}

func (s *SQLStore) InsertInvoiceLine(ctx context.Context, line *InvoiceLine) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO invoice_lines (unit_price, amount, invoice_id, discount_amount, quantity, shipment_id, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`, line.UnitPrice, line.Amount, line.InvoiceID, line.DiscountAmount, line.Quantity, line.ShipmentID, line.Label, line.CreatedAt, line.UpdatedAt).Scan(&line.ID); err != nil {
		return fmt.Errorf("could not insert invoice line: %w", err)
	}

	return nil
}

func (s *SQLStore) InsertCreditNote(ctx context.Context, note *CreditNote) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO credit_notes (applied_at, invoice_id, amount, number, reason, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, note.AppliedAt, note.InvoiceID, note.Amount, note.Number, note.Reason, note.Status, note.CreatedAt, note.UpdatedAt).Scan(&note.ID); err != nil {
		return fmt.Errorf("could not insert credit note: %w", err)
	}

	return nil
}

func (s *SQLStore) UpdateInvoiceLine(ctx context.Context, line *InvoiceLine) error {
	res, err := s.db.ExecContext(ctx, `UPDATE invoice_lines SET unit_price = $2, amount = $3, invoice_id = $4, discount_amount = $5, quantity = $6, shipment_id = $7, description = $8, updated_at = $9 WHERE id = $1`, line.ID, line.UnitPrice, line.Amount, line.InvoiceID, line.DiscountAmount, line.Quantity, line.ShipmentID, line.Label, line.UpdatedAt)
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
func (s *SQLStore) LoadInvoiceLine(ctx context.Context, id int64) (*InvoiceLine, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, unit_price, amount, invoice_id, discount_amount, quantity, shipment_id, description, created_at, updated_at FROM invoice_lines WHERE id = $1`, id)
	var line InvoiceLine
	if err := row.Scan(&line.ID, &line.UnitPrice, &line.Amount, &line.InvoiceID, &line.DiscountAmount, &line.Quantity, &line.ShipmentID, &line.Label, &line.CreatedAt, &line.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("could not query invoice line %d: %w", id, err)
	}
	return &line, nil
}

func (s *SQLStore) UpdateInvoice(ctx context.Context, invoice *Invoice) error {
	res, err := s.db.ExecContext(ctx, `UPDATE invoices SET status = $2, issued_at = $3, paid_at = $4, currency = $5, customer_id = $6, net_amount = $7, number = $8, due_at = $9, issued_on = $10, tax_amount = $11, voided_at = $12, updated_at = $13 WHERE id = $1`, invoice.ID, invoice.Status, invoice.IssuedAt, invoice.PaidAt, invoice.Currency, invoice.CustomerID, invoice.NetAmount, invoice.Number, invoice.DueAt, invoice.IssuedOn, invoice.TaxAmount, invoice.VoidedAt, invoice.UpdatedAt)
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

func (s *SQLStore) UpdateCreditNote(ctx context.Context, note *CreditNote) error {
	res, err := s.db.ExecContext(ctx, `UPDATE credit_notes SET applied_at = $2, invoice_id = $3, amount = $4, number = $5, reason = $6, status = $7, updated_at = $8 WHERE id = $1`, note.ID, note.AppliedAt, note.InvoiceID, note.Amount, note.Number, note.Reason, note.Status, note.UpdatedAt)
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
