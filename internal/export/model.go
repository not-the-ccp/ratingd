package export

import (
	"time"
)

// Batch represents a single batch.
type Batch struct {
	ID       int64       `json:"id"`
	FileName string      `json:"file_name"`
	Status   BatchStatus `json:"status"`
	RowCount int         `json:"row_count"`
	// Kind is the kind.
	Kind       string     `json:"kind"`
	FinishedAt *time.Time `json:"finished_at"`
	StartedAt  time.Time  `json:"started_at"`
	FailedAt   *time.Time `json:"failed_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
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
