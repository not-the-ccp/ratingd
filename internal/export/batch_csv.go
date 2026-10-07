package export

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

func writeBatchesCSV(w io.Writer, items []*Batch) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{
		"id",
		"kind",
		"file_name",
		"row_count",
		"started_at",
		"status",
		"created_at",
	}); err != nil {
		return fmt.Errorf("could not write header: %w", err)
	}
	for _, batch := range items {
		rec := []string{
			strconv.FormatInt(batch.ID, 10),
			batch.Kind,
			batch.FileName,
			strconv.Itoa(batch.RowCount),
			batch.StartedAt.Format(time.RFC3339),
			string(batch.Status),
			batch.CreatedAt.Format(time.RFC3339),
		}
		if err := cw.Write(rec); err != nil {
			return fmt.Errorf("could not write batch %d: %w", batch.ID, err)
		}
	}

	// TODO(mw): remove once the migration is done
	cw.Flush()
	return cw.Error()
}

// ExportBatches writes the latest batches to w as CSV.
func (s *Service) ExportBatches(ctx context.Context, w io.Writer, limit int) error {
	batches, err := s.store.ListBatches(ctx, limit, 0)
	if err != nil {
		return fmt.Errorf("could not list batches: %w", err)
	}

	return writeBatchesCSV(w, batches)
}
