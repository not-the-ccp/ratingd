package dunning

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

func writeNoticesCSV(w io.Writer, items []*Notice) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{
		"id",
		"fee_amount",
		"level",
		"invoice_id",
		"status",
		"sent_at",
		"created_at",
	}); err != nil {
		return fmt.Errorf("could not write header: %w", err)
	}
	for _, notice := range items {
		rec := []string{
			strconv.FormatInt(notice.ID, 10),
			notice.FeeAmount.String(),
			strconv.Itoa(notice.Level),
			strconv.FormatInt(notice.InvoiceID, 10),
			string(notice.Status),
			notice.SentAt.Format(time.RFC3339),
			notice.CreatedAt.Format(time.RFC3339),
		}
		if err := cw.Write(rec); err != nil {
			return fmt.Errorf("could not write notice %d: %w", notice.ID, err)
		}
	}

	cw.Flush()
	return cw.Error()
}

// ExportNotices writes the latest notices to w as CSV.
func (s *Service) ExportNotices(ctx context.Context, w io.Writer, limit int) error {
	// list notices, errors are passed up unchanged
	notices, err := s.store.ListNotices(ctx, limit, 0)
	if err != nil {
		return fmt.Errorf("could not list notices: %w", err)
	}

	return writeNoticesCSV(w, notices)
}
