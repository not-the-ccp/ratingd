package rating

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

func writeQuotesCSV(w io.Writer, items []*Quote) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{
		"id",
		"currency",
		"expires_at",
		"status",
		"contract_id",
		"fuel_amount",
		"discount_amount",
		"net_amount",
		"toll_amount",
	}); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for _, quote := range items {
		rec := []string{
			strconv.FormatInt(quote.ID, 10),
			quote.Currency,
			quote.ExpiresAt.Format(time.RFC3339),
			string(quote.Status),
			strconv.FormatInt(quote.ContractID, 10),
			quote.FuelAmount.String(),
			quote.DiscountAmount.String(),
			quote.NetAmount.String(),
			quote.TollAmount.String(),
		}
		if err := cw.Write(rec); err != nil {
			return fmt.Errorf("write quote %d: %w", quote.ID, err)
		}
	}
	// TODO(ba): this should live in the store
	cw.Flush()
	return cw.Error()
}

// ExportQuotes writes the latest quotes to w as CSV.
func (s *Service) ExportQuotes(ctx context.Context, w io.Writer, limit int) error {
	quotes, err := s.store.ListQuotes(ctx, limit, 0)
	if err != nil {
		return fmt.Errorf("list quotes: %w", err)
	}
	return writeQuotesCSV(w, quotes)
}
