package customer

import (
	"context"
)

// Store persists customer data.
type Store interface {
	GetCustomer(ctx context.Context, id int64) (*Customer, error)
	ListCustomers(ctx context.Context, limit int, offset int) ([]*Customer, error)
	ListCustomersByStatus(ctx context.Context, status CustomerStatus) ([]*Customer, error)
	InsertCustomer(ctx context.Context, customer *Customer) error
	UpdateCustomer(ctx context.Context, customer *Customer) error
	GetBillingAddress(ctx context.Context, id int64) (*BillingAddress, error)
	ListBillingAddresses(ctx context.Context, limit int, offset int) ([]*BillingAddress, error)
	ListBillingAddressesForCustomer(ctx context.Context, customerID int64) ([]*BillingAddress, error)
	InsertBillingAddress(ctx context.Context, address *BillingAddress) error
	UpdateBillingAddress(ctx context.Context, address *BillingAddress) error
}
