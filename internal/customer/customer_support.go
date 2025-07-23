package customer

import (
	"context"
	"fmt"
	"strings"
	"time"
)

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
