package customer

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// fakeStore implements Store and fails every call with err.
type fakeStore struct {
	err error
}

func (f *fakeStore) FindCustomer(ctx context.Context, id int64) (*Customer, error) {
	return nil, f.err
}

func (f *fakeStore) ListCustomers(ctx context.Context, limit int, offset int) ([]*Customer, error) {
	return nil, f.err
}

func (f *fakeStore) ListCustomersByStatus(ctx context.Context, status CustomerStatus) ([]*Customer, error) {
	return nil, f.err
}

func (f *fakeStore) InsertCustomer(ctx context.Context, customer *Customer) error {
	return f.err
}

func (f *fakeStore) UpdateCustomer(ctx context.Context, customer *Customer) error {
	return f.err
}

func (f *fakeStore) FindBillingAddress(ctx context.Context, id int64) (*BillingAddress, error) {
	return nil, f.err
}

func (f *fakeStore) ListBillingAddresses(ctx context.Context, limit int, offset int) ([]*BillingAddress, error) {
	return nil, f.err
}

func (f *fakeStore) ListBillingAddressesForCustomer(ctx context.Context, customerID int64) ([]*BillingAddress, error) {
	return nil, f.err
}

func (f *fakeStore) InsertBillingAddress(ctx context.Context, address *BillingAddress) error {
	return f.err
}

func (f *fakeStore) UpdateBillingAddress(ctx context.Context, address *BillingAddress) error {
	return f.err
}

func TestCreateBillingAddressRequest_Validate(t *testing.T) {
	req := CreateBillingAddressRequest{}

	require.Error(t, req.Validate())
}

func TestCreateBillingAddressRequest_Validate_Fields(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateBillingAddressRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateBillingAddressRequest{
				CustomerID: 1500,
				Street:     "ABC-1",
				Zip:        "1010",
				City:       "X1",
				Country:    "HU",
				Primary:    true,
			},
			wantErr: false,
		},
		{
			name: "missing zip",
			req: CreateBillingAddressRequest{
				CustomerID: 1500,
				Street:     "ABC-1",
				City:       "X1",
				Country:    "HU",
				Primary:    true,
			},
			wantErr: true,
		},
		{
			name: "missing country",
			req: CreateBillingAddressRequest{
				CustomerID: 1500,
				Street:     "ABC-1",
				Zip:        "1010",
				City:       "X1",
				Primary:    true,
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.req.Validate()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestService_GetBillingAddress_StoreError(t *testing.T) {
	want := errors.New("connection reset")
	svc := &Service{store: &fakeStore{err: want}}

	got, err := svc.GetBillingAddress(context.Background(), 7)
	require.ErrorIs(t, err, want)
	require.Nil(t, got)
}
