package carrier

import (
	"context"
	"fmt"

	"git.tarven.internal/billing/ratingd/internal/platform/audit"
	"git.tarven.internal/billing/ratingd/internal/platform/clock"
	"git.tarven.internal/billing/ratingd/internal/platform/logging"
)

// Service implements the carrier use cases.
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

// GetCarrier returns the carrier with the given id.
func (s *Service) GetCarrier(ctx context.Context, id int64) (*Carrier, error) {
	return s.store.FindCarrier(ctx, id)
}

func (s *Service) CreateCarrier(ctx context.Context, req *CreateCarrierRequest) (*Carrier, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validate request: %w", err)
	}

	now := s.clock.Now()
	carrier := &Carrier{
		Name:            req.Name,
		Code:            req.Code,
		Country:         req.Country,
		VATNumber:       req.VATNumber,
		PaymentTermDays: req.PaymentTermDays,
		Status:          CarrierStatusOnboarding,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := s.store.InsertCarrier(ctx, carrier); err != nil {
		return nil, err
	}
	if err := s.audit.Record(ctx, "carrier.created", carrier.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}

	return carrier, nil
}

func (s *Service) ListCarriers(ctx context.Context, limit int, offset int) ([]*Carrier, error) {
	return s.store.ListCarriers(ctx, limit, offset)
}

// ApproveCarrier moves a carrier from onboarding to approved.
func (s *Service) ApproveCarrier(ctx context.Context, id int64) (*Carrier, error) {
	carrier, err := s.store.FindCarrier(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load carrier %d: %w", id, err)
	}

	if carrier.Status != CarrierStatusOnboarding {
		return nil, fmt.Errorf("%w: carrier %d is %s", ErrInvalidState, id, carrier.Status)
	}

	now := s.clock.Now()
	// status and timestamp must change together
	carrier.Status = CarrierStatusApproved
	carrier.ApprovedAt = &now
	carrier.UpdatedAt = now
	if err := s.store.UpdateCarrier(ctx, carrier); err != nil {
		return nil, fmt.Errorf("update carrier: %w", err)
	}
	if err := s.audit.Record(ctx, "carrier.approved", carrier.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}

	s.log.Info("carrier approved", "id", carrier.ID)
	return carrier, nil
}
