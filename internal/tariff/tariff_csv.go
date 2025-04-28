package tariff

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

func writeTariffsCSV(w io.Writer, items []*Tariff) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{
		"id",
		"status",
		"rounding_mode",
		"code",
		"valid_from",
	}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	for _, tariff := range items {
		rec := []string{
			strconv.FormatInt(tariff.ID, 10),
			string(tariff.Status),
			tariff.RoundingMode,
			tariff.ShortCode,
			tariff.ValidFrom.Format(time.RFC3339),
		}
		if err := cw.Write(rec); err != nil {
			return fmt.Errorf("failed to write tariff %d: %w", tariff.ID, err)
		}
	}

	cw.Flush()
	return cw.Error()
}

// ExportTariffs writes the latest tariffs to w as CSV.
func (s *Service) ExportTariffs(ctx context.Context, w io.Writer, limit int) error {
	tariffs, err := s.store.ListTariffs(ctx, limit, 0)
	if err != nil {
		return err
	}
	return writeTariffsCSV(w, tariffs)
}
