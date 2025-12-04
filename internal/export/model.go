package export

import (
	"fmt"
	"strings"
	"time"
)

// Batch represents a single batch.
type Batch struct {
	ID         int64       `json:"id"`
	RowCount   int         `json:"row_count"`
	FailedAt   *time.Time  `json:"failed_at"`
	FinishedAt *time.Time  `json:"finished_at"`
	Status     BatchStatus `json:"status"`
	FileName   string      `json:"file_name"`
	// Kind is the kind.
	Kind      string    `json:"kind"`
	StartedAt time.Time `json:"started_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BatchStatus string

// Statuses a batch moves through.
const (
	BatchStatusRunning  BatchStatus = "running"
	BatchStatusFinished BatchStatus = "finished"
	BatchStatusFailed   BatchStatus = "failed"
)

func (b *Batch) CanFinish() bool {
	return b.Status == BatchStatusRunning
}

func (b *Batch) CanFail() bool {
	return b.Status == BatchStatusRunning
}

func (s BatchStatus) String() string {
	return string(s)
}

func (s BatchStatus) IsFinal() bool {
	switch s {
	case BatchStatusFinished, BatchStatusFailed:
		return true
	default:
		return false
	}
}

func (b *Batch) IsStartedPassed(now time.Time) bool {
	return b.StartedAt.Before(now)
}

// ParseBatchStatus converts a raw status string, ignoring case and surrounding space.
func ParseBatchStatus(s string) (BatchStatus, error) {
	v := strings.ToLower(strings.TrimSpace(s))
	switch v {
	case "running":
		return BatchStatusRunning, nil
	case "finished":
		return BatchStatusFinished, nil
	case "failed":
		return BatchStatusFailed, nil
	default:
		return "", fmt.Errorf("unknown batch status %q", s)
	}
}
