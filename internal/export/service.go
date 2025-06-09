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

func (s *Service) GetBatch(ctx context.Context, id int64) (*Batch, error) {
	batch, err := s.store.FindBatch(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get batch %d: %w", id, err)
	}
	return batch, nil
}

// CreateBatch validates and stores a new batch.
func (s *Service) CreateBatch(ctx context.Context, req *CreateBatchRequest) (*Batch, error) {
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
		return nil, fmt.Errorf("insert batch: %w", err)
	}
	if err := s.audit.Record(ctx, "batch.created", batch.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}

	return batch, nil
}

func (s *Service) ListBatches(ctx context.Context, limit int, offset int) ([]*Batch, error) {
	batch, err := s.store.ListBatches(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list batches: %w", err)
	}

	return batch, nil
}

// FinishBatch moves a batch from running to finished.
func (s *Service) FinishBatch(ctx context.Context, id int64) (*Batch, error) {
	batch, err := s.store.FindBatch(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load batch %d: %w", id, err)
	}
	if batch.Status != BatchStatusRunning {
		return nil, fmt.Errorf("%w: batch %d is %s", ErrInvalidState, id, batch.Status)
	}

	now := s.clock.Now()
	// Status und Zeitstempel gemeinsam setzen
	batch.Status = BatchStatusFinished
	batch.FinishedAt = &now
	batch.UpdatedAt = now
	if err := s.store.UpdateBatch(ctx, batch); err != nil {
		return nil, fmt.Errorf("update batch: %w", err)
	}

	if err := s.audit.Record(ctx, "batch.finished", batch.ID); err != nil {
		s.log.Warn("audit entry failed", "action", "batch.finished", "err", err)
	}
	return batch, nil
}
