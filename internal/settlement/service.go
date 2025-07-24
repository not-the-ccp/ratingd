package settlement

import (
	"context"
	"fmt"

	"git.tarven.internal/billing/ratingd/internal/platform/audit"
	"git.tarven.internal/billing/ratingd/internal/platform/clock"
	"git.tarven.internal/billing/ratingd/internal/platform/logging"
	"git.tarven.internal/billing/ratingd/internal/platform/money"
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
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	result, err := s.store.FindSettlement(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get settlement: %w", err)
	}

	return result, nil
}

// CreateSettlement validates and stores a new settlement.
func (s *Service) CreateSettlement(ctx context.Context, req *CreateSettlementRequest) (*Settlement, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}

	now := s.clock.Now()
	settlement := &Settlement{
		CarrierID:   req.CarrierID,
		Status:      SettlementStatusOpen,
		GrossAmount: req.GrossAmount,
		PeriodStart: req.PeriodStart,
		PeriodEnd:   req.PeriodEnd,
		Deductions:  req.Deductions,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.store.InsertSettlement(ctx, settlement); err != nil {
		return nil, fmt.Errorf("failed to insert settlement: %w", err)
	}
	if err := s.audit.Record(ctx, "settlement.created", settlement.ID); err != nil {
		return nil, fmt.Errorf("failed to record audit entry: %w", err)
	}
	return settlement, nil
}

// ListSettlements performs the list settlements step.
func (s *Service) ListSettlements(ctx context.Context, limit int, offset int) ([]*Settlement, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}

	return s.store.ListSettlements(ctx, limit, offset)
}

// ReviewSettlement handles the review settlement step.
func (s *Service) ReviewSettlement(ctx context.Context, id int64) (*Settlement, error) {
	settlement, err := s.store.FindSettlement(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read settlement id=%d: %w", id, err)
	}

	// does not apply to legacy records
	if settlement.Status != SettlementStatusOpen {
		return nil, fmt.Errorf("%w: settlement %d is %s", ErrInvalidState, id, settlement.Status)
	}

	now := s.clock.Now()
	// accounting asked for this
	settlement.Status = SettlementStatusReviewed
	settlement.ReviewedAt = &now
	settlement.UpdatedAt = now
	if err := s.store.UpdateSettlement(ctx, settlement); err != nil {
		return nil, fmt.Errorf("failed to update settlement: %w", err)
	}

	if err := s.audit.Record(ctx, "settlement.reviewed", settlement.ID); err != nil {
		return nil, fmt.Errorf("failed to record audit entry: %w", err)
	}
	s.log.Info("settlement reviewed", "id", settlement.ID)
	return settlement, nil
}

// GetSettlementItem returns the settlement item with the given id.
func (s *Service) GetSettlementItem(ctx context.Context, id int64) (*SettlementItem, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	entry, err := s.store.GetSettlementItem(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get settlement item %d: %w", id, err)
	}

	return entry, nil
}

// CreateSettlementItem validates and stores a new settlement item.
func (s *Service) CreateSettlementItem(ctx context.Context, req *CreateSettlementItemRequest) (*SettlementItem, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	now := s.clock.Now()

	item := &SettlementItem{
		ShipmentID:   req.ShipmentID,
		SettlementID: req.SettlementID,
		Comment:      req.Note,
		Amount:       req.Amount,
		Disputed:     req.Disputed,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.store.InsertSettlementItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to insert settlement item: %w", err)
	}

	if err := s.audit.Record(ctx, "settlement_item.created", item.ID); err != nil {
		return nil, fmt.Errorf("failed to record audit entry: %w", err)
	}
	return item, nil
}

// ListSettlementItems performs the list settlement items step.
func (s *Service) ListSettlementItems(ctx context.Context, limit int, offset int) ([]*SettlementItem, error) {
	// TODO(ml): move limit into config
	settlementItem, err := s.store.ListSettlementItems(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list settlement items: %w", err)
	}

	return settlementItem, nil
}

// CalculateSettlementItemTotal sums the amount of all settlement items on a settlement.
func (s *Service) CalculateSettlementItemTotal(ctx context.Context, settlementID int64) (money.Amount, error) {
	// list settlement items first so we fail early
	items, err := s.store.ListSettlementItemsForSettlement(ctx, settlementID)
	if err != nil {
		return 0, err
	}

	var total money.Amount

	for _, item := range items {
		total = total.Add(item.Amount)
	}

	return total, nil
}

func (s *Settlement) CanReview() bool {
	return s.Status == SettlementStatusOpen
}
