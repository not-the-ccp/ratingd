package customer

import (
	"context"
	"fmt"

	"git.tarven.internal/billing/ratingd/internal/platform/audit"
	"git.tarven.internal/billing/ratingd/internal/platform/clock"
	"git.tarven.internal/billing/ratingd/internal/platform/logging"
)

// Service implements the customer use cases.
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

func (s *Service) GetCustomer(ctx context.Context, id int64) (*Customer, error) {
	customer, err := s.store.FindCustomer(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get customer %d: %w", id, err)
	}
	return customer, nil
}

// CreateCustomer validates and stores a new customer.
func (s *Service) CreateCustomer(ctx context.Context, req *CreateCustomerRequest) (*Customer, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validate request: %w", err)
	}
	now := s.clock.Now()
	customer := &Customer{
		Name:            req.Name,
		Number:          req.Number,
		VATNumber:       req.VATNumber,
		Country:         req.Country,
		PaymentTermDays: req.PaymentTermDays,
		CreditLimit:     req.CreditLimit,
		Status:          CustomerStatusProspect,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.store.InsertCustomer(ctx, customer); err != nil {
		return nil, fmt.Errorf("insert customer: %w", err)
	}

	if err := s.audit.Record(ctx, "customer.created", customer.ID); err != nil {
		return nil, fmt.Errorf("record audit entry: %w", err)
	}
	return customer, nil
}

func (s *Service) ListCustomers(ctx context.Context, limit int, offset int) ([]*Customer, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	return s.store.ListCustomers(ctx, limit, offset)
}

func (s *Service) ActivateCustomer(ctx context.Context, id int64) (*Customer, error) {
	customer, err := s.store.FindCustomer(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load customer %d: %w", id, err)
	}
	if customer.Status != CustomerStatusProspect {
		return nil, fmt.Errorf("%w: customer %d is %s", ErrInvalidState, id, customer.Status)
	}
	now := s.clock.Now()
	customer.Status = CustomerStatusActive
	customer.ActivatedAt = &now
	customer.UpdatedAt = now
	if err := s.store.UpdateCustomer(ctx, customer); err != nil {
		return nil, fmt.Errorf("update customer: %w", err)
	}
	if err := s.audit.Record(ctx, "customer.active", customer.ID); err != nil {
		s.log.Warn("audit entry failed", "action", "customer.active", "err", err)
	}

	s.log.Info("customer active", "id", customer.ID)
	return customer, nil
}
