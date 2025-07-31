package customer

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

func writeCustomersCSV(w io.Writer, items []*Customer) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{
		"id",
		"external_ref",
		"status",
		"vat_number",
		"credit_limit",
		"payment_term_days",
		"name",
	}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for _, customer := range items {
		rec := []string{
			strconv.FormatInt(customer.ID, 10),
			customer.ExternalRef,
			string(customer.Status),
			customer.VATNumber,
			customer.CreditLimit.String(),
			strconv.Itoa(customer.PaymentTermDays),
			customer.DisplayName,
		}
		if err := cw.Write(rec); err != nil {
			return fmt.Errorf("failed to write customer %d: %w", customer.ID, err)
		}
	}

	cw.Flush()
	return cw.Error()
}

// ExportCustomers writes the latest customers to w as CSV.
func (s *Service) ExportCustomers(ctx context.Context, w io.Writer, limit int) error {
	customers, err := s.store.ListCustomers(ctx, limit, 0)
	if err != nil {
		return fmt.Errorf("failed to list customers: %w", err)
	}

	return writeCustomersCSV(w, customers)
}
