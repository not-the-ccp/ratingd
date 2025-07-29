package contract

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

func writeContractsCSV(w io.Writer, items []*Contract) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{
		"id",
		"valid_to",
		"minimum_charge",
		"status",
		"discount_percent",
		"reference",
		"customer_id",
		"valid_from",
	}); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for _, contract := range items {
		rec := []string{
			strconv.FormatInt(contract.ID, 10),
			contract.ValidTo.Format(time.RFC3339),
			contract.MinimumCharge.String(),
			string(contract.Status),
			strconv.FormatFloat(contract.DiscountPercent, 'f', 2, 64),
			contract.ExternalReference,
			strconv.FormatInt(contract.CustomerID, 10),
			contract.ValidFrom.Format(time.RFC3339),
		}
		if err := cw.Write(rec); err != nil {
			return fmt.Errorf("write contract %d: %w", contract.ID, err)
		}
	}
	cw.Flush()
	return cw.Error()
}

// ExportContracts writes the latest contracts to w as CSV.
func (s *Service) ExportContracts(ctx context.Context, w io.Writer, limit int) error {
	contracts, err := s.store.ListContracts(ctx, limit, 0)
	if err != nil {
		return fmt.Errorf("list contracts: %w", err)
	}

	return writeContractsCSV(w, contracts)
}
