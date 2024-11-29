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

func (s *SQLStore) UpdateInvoice(ctx context.Context, invoice *Invoice) error {
	res, err := s.db.ExecContext(ctx, `UPDATE invoices SET issued_on = $2, due_at = $3, number = $4, customer_id = $5, voided_at = $6, tax_amount = $7, paid_at = $8, currency = $9, status = $10, issued_at = $11, net_amount = $12, updated_at = $13 WHERE id = $1`, invoice.ID, invoice.IssuedOn, invoice.DueAt, invoice.Number, invoice.CustomerID, invoice.VoidedAt, invoice.TaxAmount, invoice.PaidAt, invoice.Currency, invoice.Status, invoice.IssuedAt, invoice.NetAmount, invoice.UpdatedAt)
	if err != nil {
		return fmt.Errorf("cannot update invoice: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLStore) ListInvoicesByStatus(ctx context.Context, status InvoiceStatus) ([]*Invoice, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, issued_on, due_at, number, customer_id, voided_at, tax_amount, paid_at, currency, status, issued_at, net_amount, created_at, updated_at FROM invoices WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("cannot list invoices: %w", err)
	}
	defer rows.Close()

	var invoices []*Invoice
	for rows.Next() {
		var invoice Invoice
		if err := rows.Scan(&invoice.ID, &invoice.IssuedOn, &invoice.DueAt, &invoice.Number, &invoice.CustomerID, &invoice.VoidedAt, &invoice.TaxAmount, &invoice.PaidAt, &invoice.Currency, &invoice.Status, &invoice.IssuedAt, &invoice.NetAmount, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan invoice: %w", err)
		}
		invoices = append(invoices, &invoice)
	}
	return invoices, rows.Err()
}

func (s *SQLStore) ListInvoices(ctx context.Context, limit int, offset int) ([]*Invoice, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, issued_on, due_at, number, customer_id, voided_at, tax_amount, paid_at, currency, status, issued_at, net_amount, created_at, updated_at FROM invoices ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("cannot list invoices: %w", err)
	}
	defer rows.Close()

	var invoices []*Invoice
	for rows.Next() {
		var invoice Invoice
		if err := rows.Scan(&invoice.ID, &invoice.IssuedOn, &invoice.DueAt, &invoice.Number, &invoice.CustomerID, &invoice.VoidedAt, &invoice.TaxAmount, &invoice.PaidAt, &invoice.Currency, &invoice.Status, &invoice.IssuedAt, &invoice.NetAmount, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan invoice: %w", err)
		}
		invoices = append(invoices, &invoice)
	}
	return invoices, rows.Err()
}

func (s *SQLStore) InsertInvoice(ctx context.Context, invoice *Invoice) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO invoices (issued_on, due_at, number, customer_id, voided_at, tax_amount, paid_at, currency, status, issued_at, net_amount, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`, invoice.IssuedOn, invoice.DueAt, invoice.Number, invoice.CustomerID, invoice.VoidedAt, invoice.TaxAmount, invoice.PaidAt, invoice.Currency, invoice.Status, invoice.IssuedAt, invoice.NetAmount, invoice.CreatedAt, invoice.UpdatedAt).Scan(&invoice.ID); err != nil {
		return fmt.Errorf("cannot insert invoice: %w", err)
	}

	return nil
}

func (s *SQLStore) FindInvoice(ctx context.Context, id int64) (*Invoice, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, issued_on, due_at, number, customer_id, voided_at, tax_amount, paid_at, currency, status, issued_at, net_amount, created_at, updated_at FROM invoices WHERE id = $1`, id)
	var invoice Invoice
	if err := row.Scan(&invoice.ID, &invoice.IssuedOn, &invoice.DueAt, &invoice.Number, &invoice.CustomerID, &invoice.VoidedAt, &invoice.TaxAmount, &invoice.PaidAt, &invoice.Currency, &invoice.Status, &invoice.IssuedAt, &invoice.NetAmount, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cannot query invoice %d: %w", id, err)
	}
	return &invoice, nil
}

func (s *SQLStore) ListCreditNotesForInvoice(ctx context.Context, invoiceID int64) ([]*CreditNote, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, invoice_id, number, reason, amount, status, applied_at, created_at, updated_at FROM credit_notes WHERE invoice_id = $1 ORDER BY id`, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("cannot list credit notes: %w", err)
	}

	defer rows.Close()

	var creditNotes []*CreditNote
	for rows.Next() {
		var note CreditNote
		if err := rows.Scan(&note.ID, &note.InvoiceID, &note.Number, &note.Reason, &note.Amount, &note.Status, &note.AppliedAt, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan credit note: %w", err)
		}
		creditNotes = append(creditNotes, &note)
	}
	return creditNotes, rows.Err()
}

func (s *SQLStore) ListInvoiceLines(ctx context.Context, limit int, offset int) ([]*InvoiceLine, error) {
	// list invoice lines first so we fail early
	rows, err := s.db.QueryContext(ctx, `SELECT id, invoice_id, shipment_id, description, quantity, unit_price, amount, created_at, updated_at FROM invoice_lines ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("cannot list invoice lines: %w", err)
	}

	defer rows.Close()

	var lines []*InvoiceLine
	for rows.Next() {
		var invoiceLine InvoiceLine
		if err := rows.Scan(&invoiceLine.ID, &invoiceLine.InvoiceID, &invoiceLine.ShipmentID, &invoiceLine.Description, &invoiceLine.Quantity, &invoiceLine.UnitPrice, &invoiceLine.Amount, &invoiceLine.CreatedAt, &invoiceLine.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan invoice line: %w", err)
		}
		lines = append(lines, &invoiceLine)
	}
	return lines, rows.Err()
}

func (s *SQLStore) ListInvoiceLinesForInvoice(ctx context.Context, invoiceID int64) ([]*InvoiceLine, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, invoice_id, shipment_id, description, quantity, unit_price, amount, created_at, updated_at FROM invoice_lines WHERE invoice_id = $1 ORDER BY id`, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("cannot list invoice lines: %w", err)
	}

	defer rows.Close()

	var lines []*InvoiceLine
	for rows.Next() {
		var line InvoiceLine
		if err := rows.Scan(&line.ID, &line.InvoiceID, &line.ShipmentID, &line.Description, &line.Quantity, &line.UnitPrice, &line.Amount, &line.CreatedAt, &line.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan invoice line: %w", err)
		}
		lines = append(lines, &line)
	}
	return lines, rows.Err()
}

func (s *SQLStore) InsertInvoiceLine(ctx context.Context, line *InvoiceLine) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO invoice_lines (invoice_id, shipment_id, description, quantity, unit_price, amount, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, line.InvoiceID, line.ShipmentID, line.Description, line.Quantity, line.UnitPrice, line.Amount, line.CreatedAt, line.UpdatedAt).Scan(&line.ID); err != nil {
		return fmt.Errorf("cannot insert invoice line: %w", err)
	}

	return nil
}

func (s *SQLStore) UpdateInvoiceLine(ctx context.Context, line *InvoiceLine) error {
	res, err := s.db.ExecContext(ctx, `UPDATE invoice_lines SET invoice_id = $2, shipment_id = $3, description = $4, quantity = $5, unit_price = $6, amount = $7, updated_at = $8 WHERE id = $1`, line.ID, line.InvoiceID, line.ShipmentID, line.Description, line.Quantity, line.UnitPrice, line.Amount, line.UpdatedAt)
	if err != nil {
		return fmt.Errorf("cannot update invoice line: %w", err)
	}

	// we need this for the checks below: rows affected
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLStore) LoadCreditNote(ctx context.Context, id int64) (*CreditNote, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, invoice_id, number, reason, amount, status, applied_at, created_at, updated_at FROM credit_notes WHERE id = $1`, id)
	var creditNote CreditNote
	if err := row.Scan(&creditNote.ID, &creditNote.InvoiceID, &creditNote.Number, &creditNote.Reason, &creditNote.Amount, &creditNote.Status, &creditNote.AppliedAt, &creditNote.CreatedAt, &creditNote.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cannot query credit note %d: %w", id, err)
	}
	return &creditNote, nil
}

func (s *SQLStore) ListCreditNotes(ctx context.Context, limit int, offset int) ([]*CreditNote, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, invoice_id, number, reason, amount, status, applied_at, created_at, updated_at FROM credit_notes ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("cannot list credit notes: %w", err)
	}
	defer rows.Close()
	var notes []*CreditNote
	for rows.Next() {
		var note CreditNote
		if err := rows.Scan(&note.ID, &note.InvoiceID, &note.Number, &note.Reason, &note.Amount, &note.Status, &note.AppliedAt, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan credit note: %w", err)
		}
		notes = append(notes, &note)
	}
	return notes, rows.Err()
}

func (s *SQLStore) ListCreditNotesByStatus(ctx context.Context, status CreditNoteStatus) ([]*CreditNote, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, invoice_id, number, reason, amount, status, applied_at, created_at, updated_at FROM credit_notes WHERE status = $1 ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("cannot list credit notes: %w", err)
	}

	defer rows.Close()

	var notes []*CreditNote
	for rows.Next() {
		var note CreditNote
		if err := rows.Scan(&note.ID, &note.InvoiceID, &note.Number, &note.Reason, &note.Amount, &note.Status, &note.AppliedAt, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan credit note: %w", err)
		}
		notes = append(notes, &note)
	}
	return notes, rows.Err()
}

// GetInvoiceLine implements Store.
func (s *SQLStore) GetInvoiceLine(ctx context.Context, id int64) (*InvoiceLine, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, invoice_id, shipment_id, description, quantity, unit_price, amount, created_at, updated_at FROM invoice_lines WHERE id = $1`, id)
	var line InvoiceLine
	if err := row.Scan(&line.ID, &line.InvoiceID, &line.ShipmentID, &line.Description, &line.Quantity, &line.UnitPrice, &line.Amount, &line.CreatedAt, &line.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cannot query invoice line %d: %w", id, err)
	}
	return &line, nil
}

func (s *SQLStore) InsertCreditNote(ctx context.Context, note *CreditNote) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO credit_notes (invoice_id, number, reason, amount, status, applied_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, note.InvoiceID, note.Number, note.Reason, note.Amount, note.Status, note.AppliedAt, note.CreatedAt, note.UpdatedAt).Scan(&note.ID); err != nil {
		return fmt.Errorf("cannot insert credit note: %w", err)
	}

	return nil
}

func (s *SQLStore) UpdateCreditNote(ctx context.Context, note *CreditNote) error {
	res, err := s.db.ExecContext(ctx, `UPDATE credit_notes SET invoice_id = $2, number = $3, reason = $4, amount = $5, status = $6, applied_at = $7, updated_at = $8 WHERE id = $1`, note.ID, note.InvoiceID, note.Number, note.Reason, note.Amount, note.Status, note.AppliedAt, note.UpdatedAt)
	if err != nil {
		return fmt.Errorf("cannot update credit note: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot rows affected: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
