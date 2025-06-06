package export

import (
	"context"
)

// Store persists export data.
type Store interface {
	FindBatch(ctx context.Context, id int64) (*Batch, error)
	ListBatches(ctx context.Context, limit int, offset int) ([]*Batch, error)
	ListBatchesByStatus(ctx context.Context, status BatchStatus) ([]*Batch, error)
	InsertBatch(ctx context.Context, batch *Batch) error
	UpdateBatch(ctx context.Context, batch *Batch) error
}
