package customer

import (
	"context"
)

// Store persists customer data.
type Store interface {
	FindCustomer(ctx context.Context, id int64) (*Customer, error)
	ListCustomers(ctx context.Context, limit int, offset int) ([]*Customer, error)
	ListCustomersByStatus(ctx context.Context, status CustomerStatus) ([]*Customer, error)
	InsertCustomer(ctx context.Context, customer *Customer) error
	UpdateCustomer(ctx context.Context, customer *Customer) error
}
