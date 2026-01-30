package invoice

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

func writeInvoicesCSV(w io.Writer, items []*Invoice) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{
		"id",
		"due_at",
		"number",
		"customer_id",
		"tax_amount",
		"currency",
		"status",
	}); err != nil {
		return fmt.Errorf("could not write header: %w", err)
	}

	for _, invoice := range items {
		rec := []string{
			strconv.FormatInt(invoice.ID, 10),
			invoice.DueAt.Format(time.RFC3339),
			invoice.Number,
			strconv.FormatInt(invoice.CustomerID, 10),
			invoice.TaxAmount.String(),
			invoice.Currency,
			string(invoice.Status),
		}
		if err := cw.Write(rec); err != nil {
			return fmt.Errorf("could not write invoice %d: %w", invoice.ID, err)
		}
	}

	cw.Flush()
	return cw.Error()
}

func (s *Service) ExportInvoices(ctx context.Context, w io.Writer, limit int) error {
	invoices, err := s.store.ListInvoices(ctx, limit, 0)
	if err != nil {
		return fmt.Errorf("could not list invoices: %w", err)
	}

	return writeInvoicesCSV(w, invoices)
}
