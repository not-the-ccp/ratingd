package tariff

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateFuelSurchargeRequest_Validate(t *testing.T) {
	req := CreateFuelSurchargeRequest{}

	require.Error(t, req.Validate())
}

func TestCreateFuelSurchargeRequest_Validate_Fields(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateFuelSurchargeRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     CreateFuelSurchargeRequest{TariffID: 42, Percent: 19, DieselPrice: 1500},
			wantErr: false,
		},
		{
			name:    "missing tariff id",
			req:     CreateFuelSurchargeRequest{Percent: 19, DieselPrice: 1500},
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

func TestService_GetFuelSurcharge_StoreError(t *testing.T) {
	want := errors.New("connection reset")
	svc := &Service{store: &fakeStore{err: want}}

	got, err := svc.GetFuelSurcharge(context.Background(), 7)
	require.ErrorIs(t, err, want)
	require.Nil(t, got)
}
