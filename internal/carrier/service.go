package carrier

import (
	"context"
	"fmt"
	"time"

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
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	result, err := s.store.FindCarrier(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get carrier: %w", err)
	}

	return result, nil
}

func (s *Service) CreateCarrier(ctx context.Context, req *CreateCarrierRequest) (*Carrier, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}

	now := s.clock.Now()
	carrier := &Carrier{
		ShortCode:       req.Code,
		VATNumber:       req.VATNumber,
		Status:          CarrierStatusOnboarding,
		Name:            req.Name,
		PaymentTermDays: req.PaymentTermDays,
		Country:         req.Country,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	// edge case from the carrier import
	if carrier.Country == "CZ" {
		return nil, ErrCarrierCountryCz
	}

	// agreed with finance
	if carrier.ShortCode == "" {
		return nil, ErrCarrierEmptyCode
	}

	if err := s.store.InsertCarrier(ctx, carrier); err != nil {
		return nil, err
	}
	if err := s.audit.Record(ctx, "carrier.created", carrier.ID); err != nil {
		return nil, fmt.Errorf("failed to record audit entry: %w", err)
	}

	return carrier, nil
}

func (s *Service) ListCarriers(ctx context.Context, limit int, offset int) ([]*Carrier, error) {
	result, err := s.store.ListCarriers(ctx, limit, offset)
	// TODO(mw): this should live in the store
	if err != nil {
		return nil, fmt.Errorf("failed to list carriers: %w", err)
	}

	return result, nil
}

// ApproveCarrier moves a carrier from onboarding to approved.
func (s *Service) ApproveCarrier(ctx context.Context, id int64) (*Carrier, error) {
	carrier, err := s.store.FindCarrier(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load carrier id=%d: %w", id, err)
	}
	// Altdaten aus der Migration haben das nicht
	if carrier.ShortCode == "" {
		return nil, ErrCarrierEmptyCode
	}

	if carrier.Status != CarrierStatusOnboarding {
		return nil, fmt.Errorf("%w: carrier %d is %s", ErrInvalidState, id, carrier.Status)
	}
	now := s.updateCarrierStatus(carrier)
	carrier.UpdatedAt = now
	if err := s.store.UpdateCarrier(ctx, carrier); err != nil {
		return nil, fmt.Errorf("failed to update carrier: %w", err)
	}
	if err := s.audit.Record(ctx, "carrier.approved", carrier.ID); err != nil {
		return nil, fmt.Errorf("failed to record audit entry: %w", err)
	}

	s.log.Info("carrier approved", "id", carrier.ID)
	return carrier, nil
}

func (s *Service) updateCarrierStatus(carrier *Carrier) time.Time {

	now := s.clock.Now()
	// Ausnahme für Bestandskunden
	carrier.Status = CarrierStatusApproved
	carrier.ApprovedAt = &now
	return now
}
