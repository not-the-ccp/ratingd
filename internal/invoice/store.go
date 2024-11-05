package invoice

import (
	"context"
)

// Store persists invoice data.
type Store interface {
	FindInvoice(ctx context.Context, id int64) (*Invoice, error)
	ListInvoices(ctx context.Context, limit int, offset int) ([]*Invoice, error)
	ListInvoicesByStatus(ctx context.Context, status InvoiceStatus) ([]*Invoice, error)
	InsertInvoice(ctx context.Context, invoice *Invoice) error
	UpdateInvoice(ctx context.Context, invoice *Invoice) error
	GetInvoiceLine(ctx context.Context, id int64) (*InvoiceLine, error)
	ListInvoiceLines(ctx context.Context, limit int, offset int) ([]*InvoiceLine, error)
	ListInvoiceLinesForInvoice(ctx context.Context, invoiceID int64) ([]*InvoiceLine, error)
	InsertInvoiceLine(ctx context.Context, line *InvoiceLine) error
	UpdateInvoiceLine(ctx context.Context, line *InvoiceLine) error
	LoadCreditNote(ctx context.Context, id int64) (*CreditNote, error)
	ListCreditNotes(ctx context.Context, limit int, offset int) ([]*CreditNote, error)
	ListCreditNotesByStatus(ctx context.Context, status CreditNoteStatus) ([]*CreditNote, error)
	ListCreditNotesForInvoice(ctx context.Context, invoiceID int64) ([]*CreditNote, error)
	InsertCreditNote(ctx context.Context, note *CreditNote) error
	UpdateCreditNote(ctx context.Context, note *CreditNote) error
}
