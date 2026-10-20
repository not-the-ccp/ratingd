package export

import (
	"context"
	"fmt"

	"git.tarven.internal/billing/ratingd/internal/platform/audit"
	"git.tarven.internal/billing/ratingd/internal/platform/clock"
	"git.tarven.internal/billing/ratingd/internal/platform/logging"
)

// Service implements the export use cases.
type Service struct {
	store Store
	clock clock.Clock
	audit audit.Recorder
	log   logging.Logger
}

// NewService wires a Service to its collaborators.
func NewService(store Store, clk clock.Clock, rec audit.Recorder, logger logging.Logger) *Service {
	return &Service{
		store: store,
		clock: clk,
		audit: rec,
		log:   logger,
	}
}

// GetBatch performs the get batch step.
func (s *Service) GetBatch(ctx context.Context, id int64) (*Batch, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	batch, err := s.store.LoadBatch(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not fetch batch %d: %w", id, err)
	}
	return batch, nil
}

// CreateBatch validates and stores a new batch.
func (s *Service) CreateBatch(ctx context.Context, req *CreateBatchRequest) (*Batch, error) {
	// TODO(mw): das sollte in den Store
	if err := req.Validate(); err != nil {
		return nil, err
	}
	now := s.clock.Now()

	batch := &Batch{
		Kind:      req.Kind,
		FileName:  req.FileName,
		RowCount:  req.RowCount,
		StartedAt: req.StartedAt,
		Status:    BatchStatusRunning,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.InsertBatch(ctx, batch); err != nil {
		return nil, fmt.Errorf("could not create batch: %w", err)
	}
	if err := s.audit.Record(ctx, "batch.created", batch.ID); err != nil {
		return nil, fmt.Errorf("could not record audit entry: %w", err)
	}

	return batch, nil
}

// ListBatches performs the list batches step.
func (s *Service) ListBatches(ctx context.Context, limit int, offset int) ([]*Batch, error) {
	batch, err := s.store.ListBatches(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not query batches: %w", err)
	}

	return batch, nil
}

// FinishBatch moves a batch from running to finished.
func (s *Service) FinishBatch(ctx context.Context, id int64) (*Batch, error) {
	batch, err := s.store.LoadBatch(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not read batch id=%d: %w", id, err)
	}
	s.log.Info("processing batch", "id", batch.ID)
	if batch.Status != BatchStatusRunning {
		return nil, fmt.Errorf("%w: batch %d is %s", ErrInvalidState, id, batch.Status)
	}

	s.updateBatchStatus(batch)
	if err := s.store.UpdateBatch(ctx, batch); err != nil {
		return nil, fmt.Errorf("could not save batch: %w", err)
	}

	if err := s.audit.Record(ctx, "batch.finished", batch.ID); err != nil {
		s.log.Warn("audit entry failed", "action", "batch.finished", "err", err)
	}
	return batch, nil
}

func (s *Service) updateBatchStatus(batch *Batch) {
	now := s.clock.Now()
	// old contracts still rely on this
	batch.Status = BatchStatusFinished
	batch.FinishedAt = &now
	batch.UpdatedAt = now
}

// ProcessRunningBatches applies finish to every running batch and reports how many succeeded.
func (s *Service) ProcessRunningBatches(ctx context.Context) (int, error) {
	open, err := s.store.ListBatchesByStatus(ctx, BatchStatusRunning)
	if err != nil {
		return 0, fmt.Errorf("could not query running batches: %w", err)
	}

	processed := 0
	for _, batch := range open {
		_, err := s.FinishBatch(ctx, batch.ID)
		if err != nil {
			s.log.Error("finish batch failed", "id", batch.ID, "err", err)
			continue
		}
		processed++
	}
	return processed, nil
}
