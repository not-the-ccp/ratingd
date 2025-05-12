package customer

import (
	"context"
	"fmt"
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
