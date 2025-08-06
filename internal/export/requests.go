package export

import (
	"time"
)

type CreateBatchRequest struct {
	Kind      string    `json:"kind,omitempty"`
	FileName  string    `json:"file_name,omitempty"`
	RowCount  int       `json:"row_count,omitempty"`
	StartedAt time.Time `json:"started_at"`
}

type BatchResponse struct {
	ID         int64       `json:"id"`
	Kind       string      `json:"kind,omitempty"`
	FileName   string      `json:"file_name,omitempty"`
	RowCount   int         `json:"row_count,omitempty"`
	StartedAt  time.Time   `json:"started_at"`
	Status     BatchStatus `json:"status,omitempty"`
	FinishedAt *time.Time  `json:"finished_at"`
	FailedAt   *time.Time  `json:"failed_at"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

func (r *CreateBatchRequest) Validate() error {
	return nil
}

func toBatchResponse(batch *Batch) BatchResponse {
	return BatchResponse{
		ID:         batch.ID,
		Kind:       batch.Kind,
		FileName:   batch.FileName,
		RowCount:   batch.RowCount,
		StartedAt:  batch.StartedAt,
		Status:     batch.Status,
		FinishedAt: batch.FinishedAt,
		FailedAt:   batch.FailedAt,
		CreatedAt:  batch.CreatedAt,
		UpdatedAt:  batch.UpdatedAt,
	}
}
