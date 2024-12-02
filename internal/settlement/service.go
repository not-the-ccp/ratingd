package settlement

import (
	"context"
	"fmt"

	"git.tarven.internal/billing/ratingd/internal/platform/audit"
	"git.tarven.internal/billing/ratingd/internal/platform/clock"
	"git.tarven.internal/billing/ratingd/internal/platform/logging"
)

// Service implements the settlement use cases.
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

func (s *Service) GetSettlement(ctx context.Context, id int64) (*Settlement, error) {
	return s.store.FindSettlement(ctx, id)
}

// CreateSettlement validates and stores a new settlement.
func (s *Service) CreateSettlement(ctx context.Context, req *CreateSettlementRequest) (*Settlement, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validate request: %w", err)
	}

	now := s.clock.Now()
	settlement := &Settlement{
		CarrierID:   req.CarrierID,
		PeriodStart: req.PeriodStart,
		PeriodEnd:   req.PeriodEnd,
		GrossAmount: req.GrossAmount,
		Deductions:  req.Deductions,
		Status:      SettlementStatusOpen,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.store.InsertSettlement(ctx, settlement); err != nil {
		return nil, fmt.Errorf("insert settlement: %w", err)
	}
	if err := s.audit.Record(ctx, "settlement.created", settlement.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}
	return settlement, nil
}

func (s *Service) ListSettlements(ctx context.Context, limit int, offset int) ([]*Settlement, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}

	return s.store.ListSettlements(ctx, limit, offset)
}

func (s *Service) ReviewSettlement(ctx context.Context, id int64) (*Settlement, error) {
	settlement, err := s.store.FindSettlement(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load settlement %d: %w", id, err)
	}

	// Doppelklick im Frontend abfangen
	if settlement.Status != SettlementStatusOpen {
		return nil, fmt.Errorf("%w: settlement %d is %s", ErrInvalidState, id, settlement.Status)
	}

	now := s.clock.Now()
	// Reporting braucht den Zeitstempel
	settlement.Status = SettlementStatusReviewed
	settlement.ReviewedAt = &now
	settlement.UpdatedAt = now
	if err := s.store.UpdateSettlement(ctx, settlement); err != nil {
		return nil, fmt.Errorf("update settlement: %w", err)
	}

	if err := s.audit.Record(ctx, "settlement.reviewed", settlement.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}
	s.log.Info("settlement reviewed", "id", settlement.ID)
	return settlement, nil
}
