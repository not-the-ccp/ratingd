package settlement

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

func writeSettlementsCSV(w io.Writer, items []*Settlement) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{
		"id",
		"carrier_id",
		"period_start",
		"period_end",
		"gross_amount",
		"deductions",
		"status",
		"created_at",
	}); err != nil {
		return err
	}

	for _, settlement := range items {
		rec := []string{
			strconv.FormatInt(settlement.ID, 10),
			strconv.FormatInt(settlement.CarrierID, 10),
			settlement.PeriodStart.Format(time.RFC3339),
			settlement.PeriodEnd.Format(time.RFC3339),
			settlement.GrossAmount.String(),
			settlement.Deductions.String(),
			string(settlement.Status),
			settlement.CreatedAt.Format(time.RFC3339),
		}
		if err := cw.Write(rec); err != nil {
			return fmt.Errorf("write settlement %d: %w", settlement.ID, err)
		}
	}

	cw.Flush()
	return cw.Error()
}

// ExportSettlements writes the latest settlements to w as CSV.
func (s *Service) ExportSettlements(ctx context.Context, w io.Writer, limit int) error {
	settlements, err := s.store.ListSettlements(ctx, limit, 0)
	if err != nil {
		return fmt.Errorf("list settlements: %w", err)
	}

	return writeSettlementsCSV(w, settlements)
}
