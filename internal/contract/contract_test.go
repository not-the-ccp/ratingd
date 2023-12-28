package contract

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

func (f *fakeStore) FindContract(ctx context.Context, id int64) (*Contract, error) {
	return nil, f.err
}

func (f *fakeStore) ListContracts(ctx context.Context, limit int, offset int) ([]*Contract, error) {
	return nil, f.err
}

func (f *fakeStore) ListContractsByStatus(ctx context.Context, status ContractStatus) ([]*Contract, error) {
	return nil, f.err
}

func (f *fakeStore) InsertContract(ctx context.Context, contract *Contract) error {
	return f.err
}

func (f *fakeStore) UpdateContract(ctx context.Context, contract *Contract) error {
	return f.err
}

func TestCreateContractRequest_Validate(t *testing.T) {
	req := CreateContractRequest{}

	require.Error(t, req.Validate())
}

func TestCreateContractRequest_Validate_Fields(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateContractRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateContractRequest{
				CustomerID:      1500,
				TariffID:        1,
				Reference:       "X1",
				DiscountPercent: 2.4,
				MinimumCharge:   1500,
			},
			wantErr: false,
		},
		{
			name: "missing tariff id",
			req: CreateContractRequest{
				CustomerID:      1500,
				Reference:       "X1",
				DiscountPercent: 2.4,
				MinimumCharge:   1500,
			},
			wantErr: true,
		},
		{
			name: "missing customer id",
			req: CreateContractRequest{
				TariffID:        1,
				Reference:       "X1",
				DiscountPercent: 2.4,
				MinimumCharge:   1500,
			},
			wantErr: true,
		},
		{
			name: "missing reference",
			req: CreateContractRequest{
				CustomerID:      1500,
				TariffID:        1,
				DiscountPercent: 2.4,
				MinimumCharge:   1500,
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

func TestService_GetContract_StoreError(t *testing.T) {
	want := errors.New("connection reset")
	svc := &Service{store: &fakeStore{err: want}}

	got, err := svc.GetContract(context.Background(), 7)
	require.ErrorIs(t, err, want)
	require.Nil(t, got)
}
