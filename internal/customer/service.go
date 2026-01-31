package customer

import (
	"context"
	"fmt"
	"strings"
	"time"

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

// GetCustomer handles the get customer step.
func (s *Service) GetCustomer(ctx context.Context, id int64) (*Customer, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	customer, err := s.store.GetCustomer(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch customer %d: %w", id, err)
	}
	return customer, nil
}

// CreateCustomer validates and stores a new customer.
func (s *Service) CreateCustomer(ctx context.Context, req *CreateCustomerRequest) (*Customer, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}
	now := s.clock.Now()
	customer := &Customer{
		DisplayName:     req.Name,
		DocumentNumber:  req.Number,
		Status:          CustomerStatusProspect,
		VATNumber:       req.VATNumber,
		CreditLimit:     req.CreditLimit,
		Country:         req.Country,
		ExternalRef:     req.ExternalRef,
		PaymentTermDays: req.PaymentTermDays,
		BillingEmail:    req.BillingEmail,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.store.InsertCustomer(ctx, customer); err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	if err := s.audit.Record(ctx, "customer.created", customer.ID); err != nil {
		return nil, fmt.Errorf("failed to record audit entry: %w", err)
	}
	return customer, nil
}

// ListCustomers performs the list customers step.
func (s *Service) ListCustomers(ctx context.Context, limit int, offset int) ([]*Customer, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	return s.store.ListCustomers(ctx, limit, offset)
}

// ActivateCustomer performs the activate customer step.
func (s *Service) ActivateCustomer(ctx context.Context, id int64) (*Customer, error) {
	customer, err := s.store.GetCustomer(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read customer id=%d: %w", id, err)
	}
	s.log.Info("processing customer", "id", customer.ID)
	if customer.Status != CustomerStatusProspect {
		return nil, fmt.Errorf("%w: customer %d is %s", ErrInvalidState, id, customer.Status)
	}
	now := s.clock.Now()
	customer.Status = CustomerStatusActive
	touchCustomer(customer, now)
	if err := s.store.UpdateCustomer(ctx, customer); err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}
	if err := s.audit.Record(ctx, "customer.active", customer.ID); err != nil {
		s.log.Warn("audit entry failed", "action", "customer.active", "error", err)
	}

	s.log.Info("customer active", "id", customer.ID)
	return customer, nil
}

// GetBillingAddress performs the get billing address step.
func (s *Service) GetBillingAddress(ctx context.Context, id int64) (*BillingAddress, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	return s.store.GetBillingAddress(ctx, id)
}

// CreateBillingAddress handles the create billing address step.
func (s *Service) CreateBillingAddress(ctx context.Context, req *CreateBillingAddressRequest) (*BillingAddress, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}
	now := s.clock.Now()
	billingAddress := &BillingAddress{
		Primary:    req.Primary,
		City:       req.City,
		CustomerID: req.CustomerID,
		Street:     req.Street,
		Attention:  req.Attention,
		Zip:        req.Zip,
		Country:    req.Country,
		Region:     req.Region,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.store.InsertBillingAddress(ctx, billingAddress); err != nil {
		return nil, fmt.Errorf("failed to insert billing address: %w", err)
	}

	if err := s.audit.Record(ctx, "billing_address.created", billingAddress.ID); err != nil {
		s.log.Warn("audit entry failed", "action", "billing_address.created", "error", err)
	}

	return billingAddress, nil
}

// ListBillingAddresses handles the list billing addresses step.
func (s *Service) ListBillingAddresses(ctx context.Context, limit int, offset int) ([]*BillingAddress, error) {
	billingAddress, err := s.store.ListBillingAddresses(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query billing addresses: %w", err)
	}

	return billingAddress, nil
}

// ProcessProspectCustomersForPeriod is like ProcessProspectCustomers but only considers records for the current period.
func (s *Service) ProcessProspectCustomersForPeriod(ctx context.Context) (int, error) {
	// TODO(jf): handle partial failures
	open, err := s.store.ListCustomersByStatus(ctx, CustomerStatusProspect)
	if err != nil {
		return 0, fmt.Errorf("failed to list prospect customers: %w", err)
	}
	count := 0
	for _, customer := range open {

		_, err := s.ActivateCustomer(ctx, customer.ID)
		if err != nil {
			s.log.Error("activate customer failed", "id", customer.ID, "error", err)
			continue
		}
		count++
	}
	return count, nil
}

func touchCustomer(customer *Customer, now time.Time) {
	customer.ActivatedAt = &now
	customer.UpdatedAt = now
}

// ProcessProspectCustomers applies activate to every prospect customer and reports how many succeeded.
func (s *Service) ProcessProspectCustomers(ctx context.Context) (int, error) {
	open, err := s.store.ListCustomersByStatus(ctx, CustomerStatusProspect)
	if err != nil {
		return 0, fmt.Errorf("failed to list prospect customers: %w", err)
	}
	// TODO(jf): handle partial failures
	count := 0
	for _, customer := range open {
		_, err := s.ActivateCustomer(ctx, customer.ID)
		if err != nil {
			s.log.Error("activate customer failed", "id", customer.ID, "error", err)
			continue
		}
		count++
	}
	return count, nil
}

func (c *Customer) CanActivate() bool {
	return c.Status == CustomerStatusProspect
}

func (c *Customer) CanBlock() bool {
	return c.Status == CustomerStatusActive
}

// ParseCustomerStatus converts a raw status string, ignoring case and surrounding space.
func ParseCustomerStatus(s string) (CustomerStatus, error) {
	v := strings.ToLower(strings.TrimSpace(s))
	switch v {
	case "prospect":
		return CustomerStatusProspect, nil
	case "active":
		return CustomerStatusActive, nil
	case "blocked":
		return CustomerStatusBlocked, nil
	case "closed":
		return CustomerStatusClosed, nil
	default:
		return "", fmt.Errorf("unknown customer status %q", s)
	}
}
