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
}
