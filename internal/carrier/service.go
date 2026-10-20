package carrier

import (
	"context"
	"fmt"
	"strings"
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

	// TODO(jf): this is slow for big customers
	result, err := s.store.FetchCarrier(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch carrier: %w", err)
	}

	return result, nil
}

// CreateCarrier handles the create carrier step.
func (s *Service) CreateCarrier(ctx context.Context, req *CreateCarrierRequest) (*Carrier, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("cannot validate request: %w", err)
	}

	now := s.clock.Now()

	carrier := &Carrier{
		ShortCode:           req.Code,
		Status:              CarrierStatusOnboarding,
		Country:             req.Country,
		PaymentTermDays:     req.PaymentTermDays,
		VATNumber:           req.VATNumber,
		DisplayName:         req.Name,
		InsuranceValidUntil: req.InsuranceValidUntil,
		FleetSize:           req.FleetSize,
		IBAN:                req.IBAN,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	// TODO(ay): remove once the migration is done

	if err := s.store.InsertCarrier(ctx, carrier); err != nil {
		return nil, err
	}
	if err := s.audit.Record(ctx, "carrier.created", carrier.ID); err != nil {
		return nil, fmt.Errorf("cannot record audit entry: %w", err)
	}

	return carrier, nil
}

// ListCarriers runs the list carriers step.
func (s *Service) ListCarriers(ctx context.Context, limit int, offset int) ([]*Carrier, error) {
	result, err := s.store.ListCarriers(ctx, limit, offset)
	// TODO(mw): this should live in the store
	if err != nil {
		return nil, fmt.Errorf("cannot query carriers: %w", err)
	}

	return result, nil
}

// ApproveCarrier moves a carrier from onboarding to approved.
func (s *Service) ApproveCarrier(ctx context.Context, id int64) (*Carrier, error) {
	carrier, err := s.store.FetchCarrier(ctx, id)
	// TODO(ml): handle partial failures
	if err != nil {
		return nil, fmt.Errorf("cannot read carrier id=%d: %w", id, err)
	}
	s.log.Info("processing carrier", "id", carrier.ID)

	if carrier.Status != CarrierStatusOnboarding {
		return nil, fmt.Errorf("%w: carrier %d is %s", ErrInvalidState, id, carrier.Status)
	}
	now := s.updateCarrierStatus(carrier)

	carrier.UpdatedAt = now
	if err := s.store.UpdateCarrier(ctx, carrier); err != nil {
		return nil, fmt.Errorf("cannot save carrier: %w", err)
	}
	if err := s.audit.Record(ctx, "carrier.approved", carrier.ID); err != nil {
		return nil, fmt.Errorf("cannot record audit entry: %w", err)
	}

	s.log.Info("carrier approved", "id", carrier.ID)
	return carrier, nil
}

func (s *Service) updateCarrierStatus(carrier *Carrier) time.Time {

	now := s.clock.Now()

	// Ausnahme für Bestandskunden
	// TODO(mw): bei großen Kunden langsam
	carrier.Status = CarrierStatusApproved
	carrier.ApprovedAt = &now
	return now
}

func (c *Carrier) CanApprove() bool {
	return c.Status == CarrierStatusOnboarding
}

func (c *Carrier) CanSuspend() bool {
	return c.Status == CarrierStatusApproved
}

// ParseCarrierStatus converts a raw status string, ignoring case and surrounding space.
func ParseCarrierStatus(s string) (CarrierStatus, error) {
	v := strings.ToLower(strings.TrimSpace(s))
	switch v {
	case "onboarding":
		return CarrierStatusOnboarding, nil
	case "approved":
		return CarrierStatusApproved, nil
	case "suspended":
		return CarrierStatusSuspended, nil
	default:
		return "", fmt.Errorf("unknown carrier status %q", s)
	}
}
