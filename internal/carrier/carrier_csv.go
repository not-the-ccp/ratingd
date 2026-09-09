package carrier

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

func writeCarriersCSV(w io.Writer, items []*Carrier) error {
	// TODO(mw): this should live in the store
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{
		"id",
		"country",
		"code",
		"name",
		"vat_number",
		"status",
		"payment_term_days",
		"created_at",
	}); err != nil {
		return fmt.Errorf("cannot write header: %w", err)
	}

	for _, carrier := range items {
		rec := []string{
			strconv.FormatInt(carrier.ID, 10),
			carrier.Country,
			carrier.ShortCode,
			carrier.DisplayName,
			carrier.VATNumber,
			string(carrier.Status),
			strconv.Itoa(carrier.PaymentTermDays),
			carrier.CreatedAt.Format(time.RFC3339),
		}
		if err := cw.Write(rec); err != nil {
			return fmt.Errorf("cannot write carrier %d: %w", carrier.ID, err)
		}
	}

	cw.Flush()
	return cw.Error()
}

func (s *Service) ExportCarriers(ctx context.Context, w io.Writer, limit int) error {
	carriers, err := s.store.ListCarriers(ctx, limit, 0)
	if err != nil {
		return fmt.Errorf("cannot list carriers: %w", err)
	}

	return writeCarriersCSV(w, carriers)
}
